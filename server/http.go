package server

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

type pattern string

type method string

type routeRules struct {
	methods map[method]http.Handler
}

type router struct {
	routes map[pattern]routeRules
}

func Router() *router {
	return &router{routes: make(map[pattern]routeRules)}
}

func (r *router) HandleFunc(m method, p pattern, f func(w http.ResponseWriter, req *http.Request)) {
	rules, exists := r.routes[p]

	if !exists {
		rules = routeRules{methods: make(map[method]http.Handler)}
		r.routes[p] = rules
	}

	rules.methods[m] = http.HandlerFunc(f)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if extension := filepath.Ext(path); extension != "" {
		path = filepath.Dir(path)
	}

	p, exists := r.routes[pattern(path)]
	if !exists {
		http.NotFound(w, req)
		return
	}

	h, exists := p.methods[method(req.Method)]
	if !exists {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	h.ServeHTTP(w, req)
}

func (r *router) RegisterFS(pathPrefix string, dirPrefix string, files fs.FS) {
	html, _ := fs.Sub(files, dirPrefix)
	fileHandler := http.FileServer(http.FS(html)).ServeHTTP

	fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			p := normalizePath(pathPrefix, dirPrefix, path)
			r.HandleFunc(http.MethodGet, pattern(p), fileHandler)
		}
		return nil
	})
}

func normalizePath(pathPrefix, prefix, path string) string {
	p := strings.TrimPrefix(path, prefix)
	if p == "" || p == "." {
		p = string(filepath.Separator)
	}

	return pathPrefix + p
}
