package go_template

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/robertoesteves13/go-template/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       ulid.ULID
	Name     string
	Email    string
	password []byte
}

func NewUser(name string, email string, password string) (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	return &User{
		Id:       ulid.Make(),
		Name:     name,
		Email:    email,
		password: hashed,
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.password, []byte(password))
	return err == nil
}

func UserFromDB(ctx context.Context, conn *pgxpool.Conn, email string) (*User, error) {
	db := database.New(conn)
	dbusr, err := db.GetUserByEmail(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	user := &User{
		Id:       dbusr.ID.Bytes,
		Name:     dbusr.Name.String,
		Email:    dbusr.Email.String,
		password: dbusr.Password,
	}

	return user, nil
}

func (u *User) InsertDB(ctx context.Context, conn *pgxpool.Conn) error {
	db := database.New(conn)

	return db.InsertUser(ctx, database.InsertUserParams{
		ID:        pgtype.UUID{Bytes: u.Id, Valid: true},
		Name:     pgtype.Text{String: u.Name, Valid: true},
		Email:  pgtype.Text{String: u.Email, Valid: true},
		Password:   u.password,
	})
}
