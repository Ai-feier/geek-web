//go:build v7
package web

import (
	"bytes"
	"context"
	"html/template"
	"io/fs"
)

type TemplateEngine interface {
	// Render 渲染页面
	// data 是渲染页面所需要的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}

type GoTemplateEngine struct {
	T *template.Template
	// 也可以考虑设计为 map[string]*template.Template
	// 但是其实没太大必要，因为 template.Template 本身就提供了按名索引的功能
}

func (g *GoTemplateEngine) Render(ctx context.Context,
	tplName string, data any) ([]byte, error) {
	
	bs := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(bs, tplName, data)
	return bs.Bytes(), err
}

func (g *GoTemplateEngine) LoadFromGlob(pattern string) error {
	var err error
	g.T, err = template.ParseGlob(pattern)
	return err
}

func (g *GoTemplateEngine) LoadFs(fs fs.FS, pattern string) error {
	var err error
	g.T, err = template.ParseFS(fs, pattern)
	return err
}

func (g *GoTemplateEngine) LoadFromFs(files ...string) error {
	var err error
	g.T, err = template.ParseFiles(files...)
	return err
}