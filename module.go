package caddyembed

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// embedded is what will contain your static files. The go command
// will automatically embed the files subfolder into this virtual
// file system. You can optionally change the go:embed directive
// to embed other files or folders.
//
//go:embed files
var embedded embed.FS

// files is the actual, more generic file system to be utilized.
var files fs.FS = embedded

// topFolder is the name of the top folder of the virtual
// file system. go:embed does not let us add the contents
// of a folder to the root of a virtual file system, so
// if we want to trim that root folder prefix, we need to
// also specify it in code as a string. Otherwise the
// user would need to add configuration or code to trim
// this root prefix from all filenames, e.g. specifying
// "root files" in their file_server config.
//
// It is NOT REQUIRED to change this if changing the
// go:embed directive; it is just for convenience in
// the default case.
const topFolder = "files"

func init() {
	caddy.RegisterModule(FS{})
	stripFolderPrefix()
}

// stripFolderPrefix opens the root of the file system. If it
// contains only 1 file, being a directory with the same
// name as the topFolder const, then the file system will
// be fs.Sub()'ed so the contents of the top folder can be
// accessed as if they were in the root of the file system.
// This is a convenience so most users don't have to add
// additional configuration or prefix their filenames
// unnecessarily.
func stripFolderPrefix() error {
	if f, err := files.Open("."); err == nil {
		defer f.Close()

		if dir, ok := f.(fs.ReadDirFile); ok {
			entries, err := dir.ReadDir(2)
			if err == nil &&
				len(entries) == 1 &&
				entries[0].IsDir() &&
				entries[0].Name() == topFolder {
				if sub, err := fs.Sub(embedded, topFolder); err == nil {
					files = sub
				}
			}
		}
	}
	return nil
}

// FS implements a Caddy module and fs.FS for an embedded
// file system provided by an unexported package variable.
//
// To use, simply put your files in a subfolder called
// "files", then build Caddy with your local copy of this
// plugin. Your site's files will be embedded directly
// into the binary.
//
// If the embedded file system contains only one file in
// its root which is a folder named "files", this module
// will strip that folder prefix using fs.Sub(), so that
// the contents of the folder can be accessed by name as
// if they were in the actual root of the file system.
// In other words, before: files/foo.txt, after: foo.txt.
type FS struct{}

// CaddyModule returns the Caddy module information.
func (FS) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.fs.embedded",
		New: func() caddy.Module { return new(FS) },
	}
}

func (FS) Open(name string) (fs.File, error) {
	// TODO: the file server doesn't clean up leading and trailing slashes, but embed.FS is particular so we remove them here; I wonder if the file server should be tidy in the first place
	name = strings.Trim(name, "/")
	return files.Open(name)
}

// UnmarshalCaddyfile exists so this module can be used in
// the Caddyfile, but there is nothing to unmarshal.
func (FS) UnmarshalCaddyfile(d *caddyfile.Dispenser) error { return nil }

// Interface guards
var (
	_ fs.FS                 = (*FS)(nil)
	_ caddyfile.Unmarshaler = (*FS)(nil)
)
