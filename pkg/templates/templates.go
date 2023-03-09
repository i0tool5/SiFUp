package templates

import "embed"

//go:embed html/*
var Temlpates embed.FS

//go:embed fonts/*
var Fonts embed.FS
