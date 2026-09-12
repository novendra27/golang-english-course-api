package locales

import "embed"

// FS menampung file-file kamus JSON secara tersemat (embedded) ke dalam binary Go
//
//go:embed *.json
var FS embed.FS
