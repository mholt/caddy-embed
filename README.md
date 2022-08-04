Caddy embedded file system
===========================

This Caddy plugin embeds your site directly into your web server's binary.

**NOTE:** This plugin requires building Caddy from source on your own machine because you need to add your own content to the `files` directory before compiling. The `go` command is _required_ and I recommend using [`xcaddy`](https://github.com/caddyserver/xcaddy) to build. You cannot download this plugin from the Caddy website, for example, and expect it to have your site embedded within it.

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

Somewhat annoyingly, when you use `go:embed` to add a folder, Go embeds the folder to the root of the virtual file system, without a way to configure Go to add its _contents_ to the root. Because our directive is `go:embed files`, that means all filenames have to be prefixed with `files/`. This is unintuitive as you would expect your site root, for example, to be at `index.html`, not `files/index.html`.

To counter this behavior, this module automatically ["subs the FS"](https://pkg.go.dev/io/fs#Sub) to trim that top-level folder prefix as long as the embedded directory is named `files`, and it is not moved or renamed by you.

I would recommend simply doing as the instructions say, and putting your content into the `files` folder. You can put multiple folders in there if you want more than one. But you are always welcome to do your own thing and change the go:embed directive, etc. If you do that, the automatic prefix stripping won't work for you.

If you are using `file_server` _and if you change the go:embed directive_, you can still strip the root folder name like this:

```
file_server {
	fs embedded
	root myfolder
}
```

where `myfolder` is the name of the folder you added to the `go:embed` directive. **This is only required if you are customizing the go:embed directive.**

In summary: For most people, this module will "just work" and you do not need to set or change the site root.
