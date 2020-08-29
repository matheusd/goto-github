# GOTO Github

Go app to redirect ["vanity"](#rant) package import paths to a github owner. For example,
this redirects requests for package `matheusd.com/mymod/pkg/sub` to 
`github.com/matheusd/mymode` using the rules for [remote import paths](https://golang.org/cmd/go/#hdr-Remote_import_paths).

This is meant to be used in conjunction with an existing reverse proxy for your
domain, so it plays nice with the rest of your content.

## Installing

```shell
$ go get matheusd.com/goto-github
```

## Using

```shell
$ goto-github -addr 127.0.0.1:9292 -host matheusd.com matheusd
```

### Configuring nginx

```
server {
	// ...

	error_page 418 = @goget;
	
	if ($args = "go-get=1") {
		return 418;
	}

	// ...

	location @goget {
		proxy_pass http://127.0.0.1:9292;
		add_header Cache-control: max-age 86400;
	}
}
```

# Thanks

Largely based on the original work by rsc: https://github.com/rsc/go-import-redirector

# Rant

I dislike the term "vanity" imports. The author/maintainer of package should 
have **authoritative** control over its contents. Thus, calling this style of
redirecting (`example.com/module` => `github.com/user/module`) and specifically the
non-content-serving `example.com/module` a "vanity" import implies that the
real content comes from `github.com` when this is exactly backwards from the PoV 
of the author.

For the sake of illustration, if an author decides to move his main repository
from Github to Gitlab, you'd expect users to _follow_ them (by updating their
import paths to the new provider). Thus it is the _author_ of the package that
is the ultimate source for it, **not** Github.

And so, such an author deciding to name their package after a non-content-serving
domain under their control (while providing the actual contents on some third party
content distribution platform) is **not** in fact being vain: they are in fact
asserting **their** (the author's) authority over that package and trying to
make life easier for their users on the inevitable future where that the content
distribution platform ceases to exist.
