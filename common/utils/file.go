package utils

import (
	"go-zero-douyin/common/xconst"
	"path"
	"runtime"

	"github.com/gabriel-vasile/mimetype"
)

var SupportedMimeFile = map[string]string{
	"aac.aac": "audio/aac",

	"csv.csv":        "text/csv",
	"doc.doc":        "application/msword",
	"docx.1.docx":    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"docx.docx":      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"epub.epub":      "application/epub+zip",
	"json.json":      "application/json",
	"mp3.mp3":        "audio/mpeg",
	"pdf.pdf":        "application/pdf",
	"ppt.ppt":        "application/vnd.ms-powerpoint",
	"pptx.pptx":      "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"utf16bebom.txt": "text/plain; charset=utf-16be",
	"utf16lebom.txt": "text/plain; charset=utf-16le",
	"utf32bebom.txt": "text/plain; charset=utf-32be",
	"utf32lebom.txt": "text/plain; charset=utf-32le",
	"utf8.txt":       "text/plain; charset=utf-8",
	"xls.xls":        "application/vnd.ms-excel",
	"xlsx.1.xlsx":    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"xlsx.2.xlsx":    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"xlsx.xlsx":      "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"xml.xml":        "text/xml; charset=utf-8",
	"xml.withbr.xml": "text/xml; charset=utf-8",
	"zip.zip":        "application/zip",
}

var SupportedMimeImage = map[string]string{
	"apng.png":  "image/vnd.mozilla.apng",
	"jpf.jpf":   "image/jpx",
	"jpg.jpg":   "image/jpeg",
	"png.png":   "image/png",
	"mpeg.mpeg": "video/mpeg",
}

var SupportedMimeVideo = map[string]string{
	"avi.avi":   "video/x-msvideo",
	"m4a.m4a":   "audio/x-m4a",
	"audio.mp4": "audio/mp4",
	"mkv.mkv":   "video/x-matroska",
	"mov.mov":   "video/quicktime",
	"mp4.mp4":   "video/mp4",
}

// JudgeIsSupportedFileType
//
//	@Description: 判断MIME类型是否支持，支持返回简写类型，不支持返回空
//	@param mime
//	@return string
func JudgeIsSupportedFileType(mime *mimetype.MIME) (string, int64) {
	for subName, MIME := range SupportedMimeFile {
		if mime.String() == MIME {
			return subName, xconst.FILE
		}
	}
	for subName, MIME := range SupportedMimeImage {
		if mime.String() == MIME {
			return subName, xconst.PICTURE
		}
	}
	for subName, MIME := range SupportedMimeVideo {
		if mime.String() == MIME {
			return subName, xconst.VIDEO
		}
	}
	return "", xconst.UNKONWN
}

// JudgeIsSupportedImage
//
//	@Description: 判断是不是支持的图片文件
//	@param mime
//	@return bool
func JudgeIsSupportedImage(mime *mimetype.MIME) bool {
	for _, MIME := range SupportedMimeImage {
		if mime.String() == MIME {
			return true
		}
	}
	return false
}

func JudgeIsSupportedVideo(mime *mimetype.MIME) bool {
	for _, MIME := range SupportedMimeVideo {
		if mime.String() == MIME {
			return true
		}
	}
	return false
}

// GetCurrentAbPathByCaller 获取当前执行文件绝对路径（go run）
func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// ByteToKB 文件大小转化
func ByteToKB(sizeBytes int64) int64 {
	sizeKB := sizeBytes / 1024
	return int64(sizeKB)
}
