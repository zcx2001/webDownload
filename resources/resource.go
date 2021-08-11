package resources

import "embed"

//go:embed static/*
var Content embed.FS
