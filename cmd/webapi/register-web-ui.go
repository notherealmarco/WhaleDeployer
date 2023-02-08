//go:build webui

package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/notherealmarco/WhaleDeployer/frontend"
)

func registerWebUI(hdl http.Handler) (http.Handler, error) {
	distDirectory, err := fs.Sub(frontend.Dist, "dist")
	if err != nil {
		return nil, fmt.Errorf("error embedding WebUI dist/ directory: %w", err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/dashboard/") {
			http.StripPrefix("/dashboard/", http.FileServer(http.FS(distDirectory))).ServeHTTP(w, r)
			return
		} else if r.RequestURI == "/" {
			// Redirect to dashboard
			http.Redirect(w, r, "/dashboard/", http.StatusTemporaryRedirect)
			return
		}
		hdl.ServeHTTP(w, r)
	}), nil
}
