package test

import (
	"bytes"
	"github.com/Ai-feier/geek-web"
	"html/template"
	"mime/multipart"
	"path"
	"testing"
)

func TestFileUploader_Handle(t *testing.T) {
	s := web.NewHTTPServer()
	s.GET("/upload_page", func(ctx *web.Context) {
		tpl := template.New("upload")
		tpl, err := tpl.Parse(`
<html>
<body>
	<form action="/upload" method="post" enctype="multipart/form-data">
		 <input type="file" name="myfile" />
		 <button type="submit">上传</button>
	</form>
</body>
<html>
`)
		if err != nil {
			t.Fatal(err)
		}
		page := &bytes.Buffer{}
		tpl.Execute(page, nil)
		ctx.RespStatusCode = 200
		ctx.RespData = page.Bytes()
	})
	
	s.POST("/upload", (&web.FileUploader{
		FileField:   "myfile",
		DstPathFunc: func(fh *multipart.FileHeader) string {
			return path.Join("testdata", "upload", fh.Filename)
		},
	}).Handle())
	
	s.Start(":8081")
}

func TestFileDownload_Handle(t *testing.T) {
	s := web.NewHTTPServer()
	fd := &web.FileDownloader{
		Dir: "./testdata/download",
	}
	s.GET("/download", fd.Handle())
	s.Start(":8081")
}

func TestStaticResourceHandler_Handle(t *testing.T) {
	s := web.NewHTTPServer()
	handler := web.NewStaticResourceHandler("./testdata/img", "/img")
	s.GET("/img/:file", handler.Handle)
	s.Start(":8081")
}