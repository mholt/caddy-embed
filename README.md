Caddy embedded file system
===========================

This Caddy plugin embeds your site directly into your web server's binary.

**NOTE:** This plugin requires building Caddy from source on your own machine because you need to add your own content before compiling. The `go` command is _required_ and I recommend using [`xcaddy`](https://github.com/caddyserver/xcaddy) to build. You cannot download this plugin from the Caddy website, for example, and expect it to have your site embedded within it.

## Instructions

1. Clone this repo: `git clone https://github.com/mholt/caddy-embed.git && cd caddy-embed`
2. Replace the contents of the `files` subfolder with your site.
3. Build Caddy with your locally-cloned copy of this plugin: `xcaddy build --with github.com/mholt/caddy-embed=.`

Now wherever your server goes, your site goes with it. Serve it up like this:

```
example.com

file_server {
	fs embedded
}
```

You can customize the `//go:embed` directive in the source before building if you want to choose other files or folders to embed. See the [Go `embed` package docs](https://pkg.go.dev/embed).

## Site root

Somewhat annoyingly, there doesn't seem to be a way to embed a folder's contents into the root of the file system: you _have_ to prefix filenames with the name of the folder. That means your homepage would be at `https://example.com/files/index.html` by default.

To work around this limitation, specify a root of `files`:

```
file_server {
	fs embedded
	root files
}
```
