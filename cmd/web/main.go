package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/robertoesteves13/go-template/internal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	asset_handler, err := NewAssetHandler(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize asset handler: %v", err)
		os.Exit(1)
	}

	r.Get("/assets/{filename}", asset_handler.HandleFunc)

	RegisterRoutes(r)

	err = internal.ConnectDatabase()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer internal.CloseConn()

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start server: %v", err)
		os.Exit(1)
	}
}
