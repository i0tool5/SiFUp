package templates

import "embed"

//go:embed html/*
var Temlpates embed.FS

//go:embed fonts/*
var Fonts embed.FS

type Template struct {
	templates map[string]string
}

func (t *Template) Form() string {
	return t.templates["form"]
}

func (t *Template) Genos() string {
	return t.templates["genos"]
}

// New pre-creates template instances.
func New() (*Template, error) {
	t := new(Template)
	templates := make(map[string]string)

	template, err := Temlpates.ReadFile("html/formtemplate.html")
	if err != nil {
		return nil, err
	}
	templates["form"] = string(template)

	template, err = Fonts.ReadFile("fonts/Genos.ttf")
	if err != nil {
		return nil, err
	}
	templates["genos"] = string(template)

	t.templates = templates

	return t, nil
}
