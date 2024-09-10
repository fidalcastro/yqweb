package main

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fidalcastro/yqweb/handlers"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
		panic(err)
	}
	listen_address := os.Getenv("LISTEN_ADDRESS")
	slog.Info("Starting server on " + listen_address)

	router := chi.NewMux()

	// Create a route along /files that will serve contents from
	// the ./data/ folder.
	workdir, _ := os.Getwd()
	assetsDir := filepath.Join(workdir, "assets")
	filesDir := http.Dir(assetsDir)
	slog.Info("Serving static files from " + assetsDir)
	FileServer(router, "/assets", filesDir)

	router.Get("/", handlers.GetFoo)
	router.Get("/filter", handlers.GetFilter)
	router.Post("/filter", handlers.PostFoo)

	http.ListenAndServe(":"+listen_address, router)
	slog.Info("Server is running on port " + listen_address)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
