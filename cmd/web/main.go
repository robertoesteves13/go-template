package main

import (
	"fmt"
	"net/http"
	"os"

	model "github.com/robertoesteves13/go-template"
	"github.com/robertoesteves13/go-template/cmd/web/services"
	"github.com/robertoesteves13/go-template/internal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	session_manager, err := services.NewSessionManager[model.User](os.Getenv("MEMCACHE_URL"))
	if err != nil {
		fmt.Printf("Failed to initialize session manager: %v", err)
		os.Exit(1)
	}

	asset_handler, err := services.NewAssetHandler(nil)
	r.Use(middleware.Logger, session_manager.Authenticate)

	if err != nil {
		fmt.Printf("Failed to initialize asset handler: %v", err)
		os.Exit(1)
	}
	r.Get("/assets/{filename}", asset_handler.HandleFunc)

	session_manager.LoginRoute(r, func(r *http.Request) (*model.User, error) {
		err := r.ParseForm()
		if err != nil {
			return nil, fmt.Errorf("failed to parse form")
		}
		email := r.Form.Get("email")
		pw := r.Form.Get("password")

		conn, err := internal.GetConnection(r.Context())
		if err != nil {
			return nil, fmt.Errorf("failed to get connection: %v", err)
		}
		defer conn.Release()

		user, err := model.UserFromDB(r.Context(), conn, email)
		if err != nil {
			return nil, fmt.Errorf("failed to get user from db: %v", err)
		}
		is_same_password := user.ValidatePassword(pw)

		if email == user.Email && is_same_password {
			return user, nil
		} else {
			return nil, fmt.Errorf("invalid credentials")
		}
	})
	RegisterRoutes(r)

	err = internal.ConnectDatabase()
	if err != nil {
		fmt.Printf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer internal.CloseConn()

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
