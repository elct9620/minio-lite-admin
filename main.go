package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/olivere/vite"
)

//go:embed all:dist
var distFS embed.FS

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

func main() {
	var (
		addr  = flag.String("addr", ":8080", "HTTP server address")
		isDev = flag.Bool("dev", false, "run in development mode")
	)
	flag.Parse()

	// Set up routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/server-info", serverInfoHandler)

	// Serve frontend
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			// Serve static assets
			if *isDev {
				// In dev mode, serve from filesystem
				http.FileServer(http.Dir(".")).ServeHTTP(w, r)
			} else {
				// In prod mode, serve from embedded dist
				sub, err := fs.Sub(distFS, "dist")
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				http.FileServer(http.FS(sub)).ServeHTTP(w, r)
			}
			return
		}

		// Create Vite fragment
		var viteFragment *vite.Fragment
		var err error

		if *isDev {
			viteFragment, err = vite.HTMLFragment(vite.Config{
				FS:        os.DirFS("."),
				IsDev:     true,
				ViteURL:   "http://localhost:5173",
				ViteEntry: "/src/main.ts",
			})
		} else {
			sub, err := fs.Sub(distFS, "dist")
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
	})

	// Start server
	log.Printf("Server starting on %s", *addr)
	if *isDev {
		log.Println("Running in development mode")
		log.Println("Make sure to run 'npm run dev' for the Vite dev server")
	} else {
		log.Println("Running in production mode")
	}

	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","service":"minio-lite-admin"}`)
}

func serverInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"version":"0.1.0","name":"MinIO Lite Admin"}`)
}
