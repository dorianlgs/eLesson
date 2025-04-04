package ui

import (
	"embed"
	"io/fs"
)

//go:generate npm install --force
//go:generate npm run build
//go:embed all:dist
var distDir embed.FS

var DistDirFS, _ = fs.Sub(distDir, "dist")
