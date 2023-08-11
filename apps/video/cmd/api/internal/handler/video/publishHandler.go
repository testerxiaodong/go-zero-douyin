package video

import (
	"go-zero-douyin/common/httpResponse"
	"go-zero-douyin/common/xerr"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-douyin/apps/video/cmd/api/internal/logic/video"
	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"
)

func PublishHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublishVideoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpResponse.ParamErrorResult(r, w, err)
			return
		}
		// 文件大小限制
		err := r.ParseMultipartForm(64 << 24)
		if err != nil {
			// 处理错误信息
			httpResponse.ApiResult(r, w, nil, xerr.NewFileErrMsg("文件超过限制，请上传64M以下文件，并且格式符合xls、pdf、world、xlsx等常见格式"))
		}
		l := video.NewPublishLogic(r.Context(), svcCtx)

		// 获取视频文件内容，传递给logic
		file, header, err := r.FormFile("video")
		l.Video = file
		l.VideoHeader = header
		// 获取视频封面内容，传递给logic
		file, header, err = r.FormFile("cover")
		l.VideoCover = file
		l.VideoCoverHeader = header
		// 调用logic方法
		resp, err := l.Publish(&req)
		httpResponse.ApiResult(r, w, resp, err)
	}
}
