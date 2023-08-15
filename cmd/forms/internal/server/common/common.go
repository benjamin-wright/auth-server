package common

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
)

type CommonData struct {
	Prefix       string
	Domain       string
	Title        string
	LoginLink    bool
	RegisterLink bool
}

type RenderData struct {
	Common  CommonData
	Context interface{}
}

//go:embed common.html
var commonContent string

type CommonTemplate struct {
	t *template.Template
}

func New(content string) (*CommonTemplate, error) {
	t := template.New("common")
	template.Must(t.Parse(commonContent))

	_, err := t.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %+v", err)
	}

	return &CommonTemplate{
		t: t,
	}, nil
}

func (t *CommonTemplate) Execute(w io.Writer, data RenderData) error {
	return t.t.Execute(w, data)
}
