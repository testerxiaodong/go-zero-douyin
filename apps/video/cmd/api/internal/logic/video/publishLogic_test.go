package video_test

import (
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-zero-douyin/apps/video/cmd/api/internal/logic/video"
	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"
	"go-zero-douyin/apps/video/cmd/rpc/mock"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	gloablMock "go-zero-douyin/mock"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

func TestPublishLogic_Publish(t *testing.T) {
	ctl := gomock.NewController(t)

	// 构造需要mock的接口
	mockVideoRpc := mock.NewMockVideo(ctl)
	mockValidator := gloablMock.NewMockValidator(ctl)
	mockOssClient := gloablMock.NewMockOssClient(ctl)

	// 创建publishLogic对象
	serviceContext := &svc.ServiceContext{Validator: mockValidator, VideoRpc: mockVideoRpc, OssClient: mockOssClient}
	publishLogic := video.NewPublishLogic(context.Background(), serviceContext)

	// mock具体的接口方法，实现测试逻辑

	// 参数校验失败mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// 参数校验成功，但OssClient.UploadFile调用失败mock

	// 为publishLogic对象赋值文件指针
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	ossUploadFileError := xerr.NewErrMsg("系统错误：文件上传失败")
	mockOssClient.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("", ossUploadFileError)
	mockOssClient.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("test", nil)

	// 参数校验成功，且OssClient.UploadFile调用成功,但VideoRpc.PublishVideo调用失败mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockOssClient.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("test", nil)
	mockOssClient.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("test", nil)
	mockOssClient.EXPECT().GetOssFileFullAccessPath(gomock.Any()).Return("test")
	mockOssClient.EXPECT().GetOssFileFullAccessPath(gomock.Any()).Return("test")

	videoRpcError := errors.New("videoRpc.PublishVideo error")
	mockVideoRpc.EXPECT().PublishVideo(gomock.Any(), gomock.Any()).Return(nil, videoRpcError)

	// 参数校验成功，且OssClient.UploadFile调用成功，且VideoRpc.PublishVideo调用成功mock

	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockOssClient.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("test", nil)
	mockOssClient.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("test", nil)
	mockOssClient.EXPECT().GetOssFileFullAccessPath(gomock.Any()).Return("test")
	mockOssClient.EXPECT().GetOssFileFullAccessPath(gomock.Any()).Return("test")

	expectedPublishVideoResp := &pb.PublishVideoResp{Video: &pb.VideoInfo{
		Id:       utils.NewRandomInt64(1, 10),
		Title:    utils.NewRandomString(10),
		OwnerId:  utils.NewRandomInt64(1, 10),
		PlayUrl:  utils.NewRandomString(10),
		CoverUrl: utils.NewRandomString(10),
	}}
	mockVideoRpc.EXPECT().PublishVideo(gomock.Any(), gomock.Any()).Return(expectedPublishVideoResp, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *types.PublishVideoReq
		err  error
	}{
		{
			name: "publish_video_with_validate_error",
			req:  &types.PublishVideoReq{Title: "test"},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "publish_video_with_upload_file_error",
			req:  &types.PublishVideoReq{Title: "test"},
			err:  ossUploadFileError,
		},
		{
			name: "publish_video_with_video_rpc_error",
			req:  &types.PublishVideoReq{Title: "test"},
			err:  errors.Wrapf(videoRpcError, "req: %v", &types.PublishVideoReq{Title: "test"}),
		},
		{
			name: "publish_video_success",
			req:  &types.PublishVideoReq{Title: "test"},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			publishLogic.VideoCoverHeader, publishLogic.VideoCover, publishLogic.VideoHeader, publishLogic.Video = setLogicFilePoint(t)
			resp, err := publishLogic.Publish(tc.req)
			if err != nil {
				assert.Equal(t, err.Error(), tc.err.Error())
			} else {
				assert.Equal(t, err, tc.err)
				assert.Equal(t, resp.Id, expectedPublishVideoResp.Video.Id)
				assert.Equal(t, resp.Title, expectedPublishVideoResp.Video.Title)
				assert.Equal(t, resp.OwnerId, expectedPublishVideoResp.Video.OwnerId)
				assert.Equal(t, resp.PlayUrl, expectedPublishVideoResp.Video.PlayUrl)
				assert.Equal(t, resp.CoverUrl, expectedPublishVideoResp.Video.CoverUrl)
			}
		})
	}
}

func createMultipartFileHeader(filePath string, t *testing.T) (*multipart.FileHeader, multipart.File) {
	// open the file
	file, err := os.Open(filePath)
	require.NoError(t, err)
	defer func(file *os.File) {
		err := file.Close()
		require.NoError(t, err)
	}(file)

	// create a buffer to hold the file in memory
	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)

	// create a new form and create a new file field
	formWriter := multipart.NewWriter(buffWriter)
	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
	require.NoError(t, err)

	// copy the content of the file to the form's file field
	_, err = io.Copy(formPart, file)
	require.NoError(t, err)

	// close the form writer after the copying process is finished
	// I don't use defer in here to avoid unexpected EOF error
	err = formWriter.Close()
	assert.NoError(t, err)

	// transform the bytes buffer into a form reader
	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

	// read the form components with max stored memory of 1MB
	multipartForm, err := formReader.ReadForm(1 << 20)
	require.NoError(t, err)

	// return the multipart file header
	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		log.Fatal("multipart file not exists")
		return nil, nil
	}
	f, err := files[0].Open()
	require.NoError(t, err)

	return files[0], f
}

func setLogicFilePoint(t *testing.T) (*multipart.FileHeader, multipart.File, *multipart.FileHeader, multipart.File) {
	coverFileHeader, coverFile := createMultipartFileHeader("../../tmp/go-zero.png", t)
	fileHeader, file := createMultipartFileHeader("../../tmp/SampleVideo_1280x720_1mb.mp4", t)
	return coverFileHeader, coverFile, fileHeader, file
}
