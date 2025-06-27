package http

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"os"

	"github.com/elct9620/minio-lite-admin/internal/config"
	"github.com/olivere/vite"
)

const indexTemplate = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>MinIO Lite Admin</title>
    {{ .Vite.Tags }}
  </head>
  <body>
    <div id="app"></div>
  </body>
</html>`

// FrontendHandler handles frontend requests with Vite integration
type FrontendHandler struct {
	config *config.Config
	distFS embed.FS
}

// NewFrontendHandler creates a new frontend handler
func NewFrontendHandler(cfg *config.Config, distFS embed.FS) *FrontendHandler {
	return &FrontendHandler{
		config: cfg,
		distFS: distFS,
	}
}

// ServeHTTP handles frontend requests
func (h *FrontendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		h.serveStaticAssets(w, r)
		return
	}

	h.serveIndex(w, r)
}

// serveStaticAssets serves static assets (CSS, JS, images, etc.)
func (h *FrontendHandler) serveStaticAssets(w http.ResponseWriter, r *http.Request) {
	if h.config.Server.Dev {
		// In dev mode, serve from filesystem
		http.FileServer(http.Dir(".")).ServeHTTP(w, r)
	} else {
		// In prod mode, serve from embedded dist
		sub, err := fs.Sub(h.distFS, "dist")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.FileServer(http.FS(sub)).ServeHTTP(w, r)
	}
}

// serveIndex serves the main index.html with Vite integration
func (h *FrontendHandler) serveIndex(w http.ResponseWriter, r *http.Request) {
	// Create Vite fragment
	var viteFragment *vite.Fragment
	var err error

	if h.config.Server.Dev {
		viteFragment, err = vite.HTMLFragment(vite.Config{
			FS:        os.DirFS("."),
			IsDev:     true,
			ViteURL:   h.config.Vite.URL,
			ViteEntry: h.config.Vite.Entry,
		})
	} else {
		sub, err := fs.Sub(h.distFS, "dist")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		viteFragment, err = vite.HTMLFragment(vite.Config{
			FS:    sub,
			IsDev: false,
		})
	}

	if err != nil {
		http.Error(w, "Error creating Vite fragment", http.StatusInternalServerError)
		return
	}

	// Parse and execute template
	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, map[string]interface{}{
		"Vite": viteFragment,
	}); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
