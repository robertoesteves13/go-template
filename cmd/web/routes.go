package main

import (
	"context"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oklog/ulid/v2"
	"github.com/robertoesteves13/go-template"
	"github.com/robertoesteves13/go-template/cmd/web/templates"
	"github.com/robertoesteves13/go-template/internal"
	"github.com/robertoesteves13/go-template/internal/database"
)

func RegisterRoutes(r chi.Router) {
	r.Get("/", postsFeed)
	r.Get("/post/{id}", postPage)
	r.Get("/posts/create", postCreate)
}

func postCreate(w http.ResponseWriter, r *http.Request) {
	conn, err := internal.GetConnection(r.Context())
	if err != nil {
		io.WriteString(w, "internal server error")
		w.WriteHeader(500)
		return
	}
	defer conn.Release()

	post := go_template.NewPost()
	err = post.InsertDB(r.Context(), conn)
	if err != nil {
		io.WriteString(w, "internal server error")
		w.WriteHeader(500)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func postsFeed(w http.ResponseWriter, r *http.Request) {
	conn, err := internal.GetConnection(r.Context())
	if err != nil {
		io.WriteString(w, "internal server error")
		w.WriteHeader(500)
		return
	}
	defer conn.Release()

	db := database.New(conn)
	db_posts, err := db.ListPosts(r.Context())
	if err != nil {
		io.WriteString(w, "internal server error")
		w.WriteHeader(500)
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
		io.WriteString(w, "internal server error")
		w.WriteHeader(500)
		return
	}
	defer conn.Release()

	db := database.New(conn)
	db_post, err := db.GetPost(r.Context(), pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		io.WriteString(w, "internal server error")
		w.WriteHeader(500)
		return
	}

	post := go_template.PostFromDB(db_post)

	ctx := context.WithValue(r.Context(), templates.TemplateTitle, post.Title())
	ctx = context.WithValue(ctx, templates.TemplateDescription, post.Subtitle())

	templates.Post(post).Render(ctx, w)
}
