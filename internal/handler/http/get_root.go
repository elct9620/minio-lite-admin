package http

import (
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

// GetRootHandler handles frontend requests with Vite integration
// This serves the SPA index.html for all non-API routes to support client-side routing
func (s *Service) GetRootHandler(w http.ResponseWriter, r *http.Request) {
	// Check if this is a static asset request (has file extension)
	if s.isStaticAsset(r.URL.Path) {
		s.serveStaticAssets(w, r)
		return
	}

	// For all other routes (/, /dashboard, /access-keys, etc.), serve the SPA index.html
	s.serveIndex(w, r)
}

// serveStaticAssets serves static assets (CSS, JS, images, etc.)
func (s *Service) serveStaticAssets(w http.ResponseWriter, r *http.Request) {
	if s.config.Server.Dev {
		// In dev mode, serve from filesystem
		http.FileServer(http.Dir(".")).ServeHTTP(w, r)
	} else {
		// In prod mode, serve from embedded dist
		sub, err := fs.Sub(s.distFS, "dist")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.FileServer(http.FS(sub)).ServeHTTP(w, r)
	}
}

// serveIndex serves the main index.html with Vite integration
func (s *Service) serveIndex(w http.ResponseWriter, r *http.Request) {
	// Create Vite fragment
	var viteFragment *vite.Fragment
	var err error

	if s.config.Server.Dev {
		viteFragment, err = vite.HTMLFragment(vite.Config{
			FS:        os.DirFS("."),
			IsDev:     true,
			ViteURL:   s.config.Vite.URL,
			ViteEntry: s.config.Vite.Entry,
		})
	} else {
		sub, err := fs.Sub(s.distFS, "dist")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		viteFragment, err = vite.HTMLFragment(vite.Config{
			FS:    sub,
			IsDev: false,
		})
		if err != nil {
			http.Error(w, "Error creating Vite fragment", http.StatusInternalServerError)
			return
		}
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

// isStaticAsset determines if a path is for a static asset (has file extension)
func (s *Service) isStaticAsset(path string) bool {
	// Get file extension
	ext := strings.ToLower(filepath.Ext(path))

	// List of static asset extensions
	staticExtensions := []string{
		".js", ".css", ".map", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico",
		".woff", ".woff2", ".ttf", ".eot", ".otf", ".mp4", ".webm", ".mp3",
		".json", ".xml", ".txt", ".html", ".htm",
	}

	for _, staticExt := range staticExtensions {
		if ext == staticExt {
			return true
		}
	}

	return false
}
