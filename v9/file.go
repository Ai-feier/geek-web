//go:build v9

package web

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FileUploader struct {
	// FileField 对应于文件在表单中的字段名字
	FileField string
	// DstPathFunc 用于计算目标路径
	DstPathFunc func(fh *multipart.FileHeader) string
}

func (f *FileUploader) Handle() HandleFunc {
	// 这里可以额外做一些检测
	// if f.FileField == "" {
	// 	// 这种方案默认值我其实不是很喜欢
	// 	// 因为我们需要教会用户说，这个 file 是指什么意思
	// 	f.FileField = "file"
	// }
	return func(ctx *Context) {
		src, srcHeader, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败，未找到数据")
			log.Fatalln(err)
			return
		}
		defer src.Close()
		// 打开目标文件
		dst, err := os.OpenFile(f.DstPathFunc(srcHeader),
			os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败")
			log.Fatalln(err)
			return
		}

		// 将请求中的数据读取到目标文件中
		_, err = io.CopyBuffer(dst, src, nil)
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("上传失败")
			log.Fatalln(err)
			return
		}

		ctx.RespData = []byte("上传成功")
	}
}

// FileDownloader 直接操作了 http.ResponseWriter
// 所以在 Middleware 里面将不能使用 RespData
// 因为没有赋值
type FileDownloader struct {
	Dir string
}

func (fd *FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		// 获取 url 中文件名
		req, _ := ctx.QueryValue("file").String()
		path := filepath.Join(fd.Dir, filepath.Clean(req))
		fn := filepath.Base(path)
		header := ctx.Resp.Header()                                  // 直接操作 resp
		header.Set("Content-Disposition", "attachment;filename="+fn) // * attachment: 保存到本地
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")
		http.ServeFile(ctx.Resp, ctx.Req, path)
	}
}

type StaticResourceHandlerOption func(s *StaticResourceHandler)

type StaticResourceHandler struct {
	// 目录
	dir string
	// 路径前缀
	pathPrefix string
	// contentType 类型
	extensionContentTypeMap map[string]string

	// 缓存静态资源的限制
	cache       *lru.Cache
	maxFileSize int
}

// 缓存的数据类型
type fileCacheItem struct {
	fileName    string
	fileSize    int
	contentType string
	data        []byte
}

func NewStaticResourceHandler(dir string, pathPrefix string,
	opts ...StaticResourceHandlerOption) *StaticResourceHandler {
	res := &StaticResourceHandler{
		dir:        dir,
		pathPrefix: pathPrefix,
		extensionContentTypeMap: map[string]string{
			// 这里根据自己的需要不断添加
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (h *StaticResourceHandler) Handle(ctx *Context) {
	// 获取文件名字
	req, _ := ctx.PathValue("file").String()
	if item, ok := h.readFileFromData(req); ok {
		log.Printf("从缓存中读取数据...")
		h.writeItemAsResponse(item, ctx.Resp)
		return
	}
	path := filepath.Join(h.dir, req)
	f, err := os.Open(path)
	if err != nil {
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 获取文件后缀名找到相应的 contentType
	//ext := filepath.Ext(f.Name())
	ext := getFileExt(f.Name())
	// 查看是否支持当前 contentType
	t, ok := h.extensionContentTypeMap[ext]
	if !ok {
		ctx.Resp.WriteHeader(http.StatusBadRequest)
		return
	}
	
	// 读取文件
	data, err := io.ReadAll(f)
	if err != nil {
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 构造 lru item
	item := &fileCacheItem{
		fileName:    req,
		fileSize:    len(data),
		contentType: t,
		data:        data,
	}

	h.cacheFile(item)
	h.writeItemAsResponse(item, ctx.Resp)
}

// 将 item 放入 lru 缓存
func (h *StaticResourceHandler) cacheFile(item *fileCacheItem) {
	if h.cache != nil && item.fileSize < h.maxFileSize {
		h.cache.Add(item.fileName, item)
	}	
}

// 从缓存中读取数据
func (h *StaticResourceHandler) readFileFromData(fileName string) (*fileCacheItem, bool) {
	if h.cache != nil {
		if file, ok := h.cache.Get(fileName); ok {
			return file.(*fileCacheItem), ok
		}
	}
	return nil, false
}

func (h *StaticResourceHandler) writeItemAsResponse(item *fileCacheItem, writer http.ResponseWriter) {
	// 直接将缓存中的数据写回到 response
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", item.contentType)
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", item.fileSize))
	_, _ = writer.Write(item.data)
}

func getFileExt(name string) string {
	idx := strings.LastIndex(name, ".")
	if idx == len(name)-1 {
		return ""
	}
	return name[idx+1:]
}
