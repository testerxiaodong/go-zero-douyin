package video

import (
	"context"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/ctxdata"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"golang.org/x/sync/errgroup"
	"io"
	"mime/multipart"
	"path"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	logx.Logger
	ctx              context.Context
	svcCtx           *svc.ServiceContext
	Video            multipart.File
	VideoHeader      *multipart.FileHeader
	VideoCover       multipart.File
	VideoCoverHeader *multipart.FileHeader
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishVideoReq) (resp *types.PublishVideoResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}

	// 获取用户id
	uid := ctxdata.GetUidFromCtx(l.ctx)

	// 参数校验
	validateErrorGroup, _ := errgroup.WithContext(l.ctx)

	tagIdList := strings.Split(req.Tags, ",")

	// 判断标签是否存在
	for _, tag := range tagIdList {
		id, err := strconv.ParseInt(tag, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("标签id不正确"), "req: %v", req)
		}
		validateErrorGroup.Go(func() error {
			_, err := l.svcCtx.VideoRpc.GetTagById(l.ctx, &pb.GetTagByIdReq{Id: id})
			return err
		})
	}

	// 判断分区是否存在
	validateErrorGroup.Go(func() error {
		_, err := l.svcCtx.VideoRpc.GetSectionById(l.ctx, &pb.GetSectionByIdReq{Id: req.SectionId})
		return err
	})

	err = validateErrorGroup.Wait()
	if err != nil {
		return nil, err
	}

	// 文件上传
	fileErrorGroup, _ := errgroup.WithContext(l.ctx)
	// 资源文件本体
	var videoOssSubPath string
	fileErrorGroup.Go(func() error {
		fileContent, err := io.ReadAll(l.Video)
		if err != nil {
			return xerr.NewFileErrMsg("文件内容读取失败")
		}
		defer func(File multipart.File) {
			err := File.Close()
			if err != nil {

			}
		}(l.Video)
		// 检查支持文件类型
		if !utils.JudgeIsSupportedVideo(mimetype.Detect(fileContent)) {
			return xerr.NewFileErrMsg("不支持的视频类型")
		}
		// 上传文件到oss
		filename := req.Title + "-" + strconv.FormatInt(uid, 10) + "-" + utils.NewRandomString(5) + path.Ext(l.VideoHeader.Filename)
		filePath, err := l.svcCtx.OssClient.UploadFile(filename, l.svcCtx.Config.AliCloud.CommonPath, fileContent)
		if err != nil {
			l.Logger.WithFields(logx.Field("err:", err)).Error(fmt.Sprintf("上传文件到cos失败，文件名称:%v,文件路径:%v",
				filename, filePath))
			return xerr.NewErrMsg("系统错误：文件上传失败")
		}
		videoOssSubPath = filePath
		return nil
	})
	// 文件头图
	var videoCoverOssSubPath string
	fileErrorGroup.Go(func() error {
		fileContent, err := io.ReadAll(l.VideoCover)
		if err != nil {
			return xerr.NewFileErrMsg("文件内容读取失败")
		}
		defer func(FilePoster multipart.File) {
			err := FilePoster.Close()
			if err != nil {

			}
		}(l.VideoCover)
		// 不支持图片类型
		if !utils.JudgeIsSupportedImage(mimetype.Detect(fileContent)) {
			return xerr.NewFileErrMsg("不支持的图片类型")
		}
		// 上传文件到oss
		filename := req.Title + "-" + strconv.FormatInt(uid, 10) + "-" + utils.NewRandomString(5) + path.Ext(l.VideoCoverHeader.Filename)
		filePath, err := l.svcCtx.OssClient.UploadFile(filename, l.svcCtx.Config.AliCloud.CommonPath, fileContent)
		if err != nil {
			l.Logger.WithFields(logx.Field("err:", err)).Error(fmt.Sprintf("上传文件到cos失败，文件名称:%v,文件路径:%v",
				filename, filePath))
			return xerr.NewErrMsg("系统错误：文件上传失败")
		}
		videoCoverOssSubPath = filePath
		return nil
	})

	err = fileErrorGroup.Wait()
	if err != nil {
		l.Logger.WithFields(logx.Field("err:", err)).Error("头图和文件资源上传失败！")
		return nil, err
	}
	video, err := l.svcCtx.VideoRpc.PublishVideo(l.ctx, &pb.PublishVideoReq{
		Title:     req.Title,
		SectionId: req.SectionId,
		Tags:      tagIdList,
		OwnerId:   uid,
		PlayUrl:   l.svcCtx.OssClient.GetOssFileFullAccessPath(videoOssSubPath),
		CoverUrl:  l.svcCtx.OssClient.GetOssFileFullAccessPath(videoCoverOssSubPath),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}
	resp = new(types.PublishVideoResp)
	resp.Id = video.Video.Id
	resp.Title = video.Video.Title
	resp.OwnerId = video.Video.OwnerId
	resp.PlayUrl = video.Video.PlayUrl
	resp.CoverUrl = video.Video.CoverUrl
	return resp, nil
}
