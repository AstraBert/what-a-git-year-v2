package auth

import (
	"context"
	"errors"

	"github.com/AstraBert/what-a-git-year-v2/db"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrUnauthorized = errors.New("unauthorized")

func AuthorizePost(c *fiber.Ctx) error {
	sqlDb, err := CreateNewDb()
	if err != nil {
		return ErrUnauthorized
	}
	st := c.Cookies("session_token", "")
	if st == "" {
		return ErrUnauthorized
	}
	queries := db.New(sqlDb)
	ctx := context.Background()
	user, err := queries.GetUserBySessionToken(ctx, pgtype.Text{String: st, Valid: true})
	if err != nil {
		return ErrUnauthorized
	}
	csrf := c.Cookies("csrf_token", "")
	if csrf == "" {
		return ErrUnauthorized
	}
	if csrf != user.CsrfToken.String {
		return ErrUnauthorized
	}
	return nil
}

func AuthorizeGet(c *fiber.Ctx) error {
	sqlDb, err := CreateNewDb()
	if err != nil {
		return ErrUnauthorized
	}
	st := c.Cookies("session_token", "")
	if st == "" {
		return ErrUnauthorized
	}
	queries := db.New(sqlDb)
	ctx := context.Background()
	_, err = queries.GetUserBySessionToken(ctx, pgtype.Text{String: st, Valid: true})
	if err != nil {
		return ErrUnauthorized
	}
	return nil
}
