package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oklog/ulid/v2"
	"github.com/robertoesteves13/go-template"
	"github.com/robertoesteves13/go-template/cmd/web/templates"
	"github.com/robertoesteves13/go-template/internal"
	"github.com/robertoesteves13/go-template/internal/database"
)

// All your routes should be written here. You could also transform this file
// into a folder if you feel this got large enough.
func RegisterRoutes(r chi.Router) {
	r.Get("/", postsFeed)
	r.Post("/", postsFeed)
	r.Get("/post/{id}", postPage)
	r.Get("/posts/create", postCreate)

	r.Get("/login", loginPage)
	r.Get("/register", registerPage)
	r.Post("/register", registerUser)
}

func postCreate(w http.ResponseWriter, r *http.Request) {
	conn, err := internal.GetConnection(r.Context())
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	post := go_template.NewPost()
	err = post.InsertDB(r.Context(), conn)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func postsFeed(w http.ResponseWriter, r *http.Request) {
	conn, err := internal.GetConnection(r.Context())
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	db := database.New(conn)
	db_posts, err := db.ListPosts(r.Context())
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	posts := go_template.PostFromDBSlice(db_posts)

	ctx := context.WithValue(r.Context(), templates.TemplateTitle, "Posts")
	ctx = context.WithValue(ctx, templates.TemplateDescription, "List of all posts of the website")

	templates.PostsFeed(posts).Render(ctx, w)
}

func postPage(w http.ResponseWriter, r *http.Request) {
	id, err := ulid.ParseStrict(chi.URLParam(r, "id"))
	if err != nil {
		io.WriteString(w, "invalid page")
		w.WriteHeader(404)
		return
	}

	conn, err := internal.GetConnection(r.Context())
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	db := database.New(conn)
	db_post, err := db.GetPost(r.Context(), pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	post := go_template.PostFromDB(db_post)

	ctx := context.WithValue(r.Context(), templates.TemplateTitle, post.Title())
	ctx = context.WithValue(ctx, templates.TemplateDescription, post.Subtitle())

	templates.Post(post).Render(ctx, w)
}

func registerPage(w http.ResponseWriter, r *http.Request) {
	templates.RegisterPage().Render(r.Context(), w)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := go_template.NewUser(username, email, password)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	conn, err := internal.GetConnection(r.Context())
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	err = user.InsertDB(r.Context(), conn)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	templates.LoginPage().Render(r.Context(), w)
}
