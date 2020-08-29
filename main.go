package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	addr    = flag.String("addr", ":http", "serve http on `address`")
	host    = flag.String("host", "", "host to use (requested host if empty)")
	appName = filepath.Base(os.Args[0])

	destUser string
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s <github username>\n", appName)
	fmt.Fprintf(os.Stderr, "options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "examples:\n")
	fmt.Fprintf(os.Stderr, "\t%s matheusd\n", appName)
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
	}
	destUser = flag.Arg(0)
	http.HandleFunc("/", redirect)
	log.Printf("Listening on address %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var tmpl = template.Must(template.New("main").Parse(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{.ImportRoot}} {{.VCS}} {{.VCSRoot}}">
<meta http-equiv="refresh" content="0; url=https://pkg.go.dev/{{.ImportRoot}}{{.Suffix}}">
</head>
<body>
Redirecting to docs at <a href="https://pkg.go.dev/{{.ImportRoot}}{{.Suffix}}">pkg.go.dev/{{.ImportRoot}}{{.Suffix}}</a>...
</body>
</html>
`))

type data struct {
	ImportRoot string
	VCSRoot    string
	VCS        string
	Suffix     string
}

func redirect(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/")
	split := strings.SplitN(path, "/", 2)
	repo := ""
	if len(split) > 0 {
		repo = split[0]
	}

	importRoot := *host
	if importRoot == "" {
		importRoot = req.Host
	}
	importRoot = importRoot + "/" + path

	d := data{
		ImportRoot: importRoot,
		VCS:        "git",
		VCSRoot:    "https://github.com/" + destUser + "/" + repo,
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, d)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(buf.Bytes())
}
