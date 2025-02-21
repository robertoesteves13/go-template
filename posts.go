package go_template

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/robertoesteves13/go-template/internal/database"
)

type Post struct {
	id         ulid.ULID
	title      string
	subtitle   string
	content    string
	created_at time.Time
	updated_at time.Time
}

func NewPost() *Post {
	return &Post{
		ulid.Make(),
		"New Post",
		"This is a new post, what's inside it?",
		"Surprise! it's just a bunch of useless text :P.",
		time.Now(),
		time.Now(),
	}
}

func (p *Post) Id() ulid.ULID {
	return p.id
}

func (p *Post) Title() string {
	return p.title
}

func (p *Post) Subtitle() string {
	return p.subtitle
}

func (p *Post) Content() string {
	return p.content
}

func (p *Post) CreatedAt() time.Time {
	return p.created_at
}

func (p *Post) UpdatedAt() time.Time {
	return p.updated_at
}

func (p *Post) SetTitle(title string) {
	p.title = title
	p.updated_at = time.Now()
}

func (p *Post) SetSubtitle(subtitle string) {
	p.subtitle = subtitle
	p.updated_at = time.Now()
}

func (p *Post) SetContent(content string) {
	p.content = content
	p.updated_at = time.Now()
}

func (p *Post) URL() string {
	return fmt.Sprintf("/post/%s", p.Id())
}

func (p *Post) UpdateDB(ctx context.Context, conn *pgxpool.Conn) error {
	db := database.New(conn)
	var updated_at pgtype.Timestamp
	updated_at.Time = p.updated_at

	return db.UpdatePost(ctx, database.UpdatePostParams{
		Title:     pgtype.Text{String: p.title, Valid: true},
		Subtitle:  pgtype.Text{String: p.subtitle, Valid: true},
		Content:   pgtype.Text{String: p.content, Valid: true},
		UpdatedAt: updated_at,
		ID:        pgtype.UUID{Bytes: p.id, Valid: true},
	})
}

func (p *Post) InsertDB(ctx context.Context, conn *pgxpool.Conn) error {
	db := database.New(conn)
	var updated_at pgtype.Timestamp
	var created_at pgtype.Timestamp

	updated_at.Time = p.updated_at
	created_at.Time = p.created_at


	return db.InsertPost(ctx, database.InsertPostParams{
		Title:     pgtype.Text{String: p.title, Valid: true},
		Subtitle:  pgtype.Text{String: p.subtitle, Valid: true},
		Content:   pgtype.Text{String: p.content, Valid: true},
		UpdatedAt: updated_at,
		CreatedAt: created_at,
		ID:        pgtype.UUID{Bytes: p.id, Valid: true},
	})
}

func PostFromDB(post database.Post) *Post {
	return &Post{
		id:         post.ID.Bytes,
		title:      post.Title.String,
		subtitle:   post.Subtitle.String,
		content:    post.Content.String,
		created_at: post.CreatedAt.Time,
		updated_at: post.UpdatedAt.Time,
	}
}

func PostFromDBSlice(posts []database.Post) []Post {
	ps := make([]Post, 0, len(posts))
	for i := range posts {
		ps = append(ps, *PostFromDB(posts[i]))
	}

	return ps
}
