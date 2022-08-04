package caddyembed

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// The variable below is what will contain your static files.
// The go command will automatically embed the contents of the
// files subfolder into this virtual file system. You can
// optionally change the go:embed directive to embed other files
// or folders.

//go:embed files
var embedded embed.FS

func init() {
	caddy.RegisterModule(FS{})
}

// FS implements a Caddy module and fs.StatFS for an embedded
// file system provided by an unexported package variable.
//
// Simply put your files in a subfolder called `files` then
// build Caddy with your local copy of this plugin. Your
// site's files will be embedded directly into the binary.
type FS struct{}

// CaddyModule returns the Caddy module information.
func (FS) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.fs.embedded",
		New: func() caddy.Module { return new(FS) },
	}
}

func (FS) Open(name string) (fs.File, error) {
	// TODO: the file server doesn't clean up leading and trailing slashes, but embed.FS is particular so we remove them here; I wonder if the file server should be tidy in the first place (see also Stat below)
	name = strings.Trim(name, "/")
	return embedded.Open(name)
}

func (FS) Stat(name string) (fs.FileInfo, error) {
	name = strings.Trim(name, "/")
	file, err := embedded.Open(name)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}
	defer file.Close()
	return file.Stat()
}

func (FS) UnmarshalCaddyfile(d *caddyfile.Dispenser) error { return nil }

// Interface guards
var (
	_ fs.StatFS             = (*FS)(nil)
	_ caddyfile.Unmarshaler = (*FS)(nil)
)
