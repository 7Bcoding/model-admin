package templates

import (
	"html/template"
	"io"
)

type Layout struct {
	Base     string
	Content  string
	Data     map[string]interface{}
	Template *template.Template
}

func NewLayout(base, content string, data map[string]interface{}) *Layout {
	return &Layout{
		Base:    base,
		Content: content,
		Data:    data,
	}
}

func (l *Layout) Render(w io.Writer) error {
	return l.Template.ExecuteTemplate(w, l.Base, l.Data)
}
