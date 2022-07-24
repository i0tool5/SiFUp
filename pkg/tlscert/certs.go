package tlscert

import "embed"

//go:embed files/*
var CertFiles embed.FS
