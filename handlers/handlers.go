package handlers

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/AstraBert/what-a-git-year-v2/auth"
	"github.com/AstraBert/what-a-git-year-v2/db"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func HandleUserSearch() {

}

func HandleOrgSearch() {

}

func HomeRoute() {

}

func LoginRoute() {

}

func SignUpRoute() {

}

func LogoutRoute() {

}

func HandleSignUp(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	ctx := context.Background()
	sqlDb, err := auth.CreateNewDb()
	if err != nil {
		// banners := templates.SingupBanner(err)
		// return banners.Render(c.Context(), c.Response().BodyWriter())
		return c.SendStatus(500)
	}
	queries := db.New(sqlDb)
	_, err = queries.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			hashed_psw, err := auth.HashPassword(password)
			if err != nil {
				// banners := templates.SingupBanner(err)
				// return banners.Render(c.Context(), c.Response().BodyWriter())
				return c.SendStatus(500)
			}
			_, err = queries.CreateUser(ctx, db.CreateUserParams{Username: username, HashedPassword: hashed_psw})
			if err != nil {
				// banners := templates.SingupBanner(err)
				// return banners.Render(c.Context(), c.Response().BodyWriter())
				return c.SendStatus(500)
			} else {
				// banners := templates.SingupBanner(nil)
				// return banners.Render(c.Context(), c.Response().BodyWriter())
				return c.SendStatus(200)
			}
		} else {
			// banners := templates.SingupBanner(err)
			// return banners.Render(c.Context(), c.Response().BodyWriter())
			return c.SendStatus(500)
		}
	} else {
		// banners := templates.SingupBanner(errors.New("user already exists"))
		// return banners.Render(c.Context(), c.Response().BodyWriter())
		return c.SendStatus(fiber.StatusConflict)
	}
}

func HandleLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	ctx := context.Background()
	sqlDb, err := auth.CreateNewDb()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
	}
	queries := db.New(sqlDb)
	user, err := queries.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// banners := templates.SingupBanner(errors.New("there is no user with this username"))
			// return banners.Render(c.Context(), c.Response().BodyWriter())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
		} else {
			// banners := templates.SingupBanner(err)
			// return banners.Render(c.Context(), c.Response().BodyWriter())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
		}
	}
	if !auth.CompareHashToPassword(password, user.HashedPassword) {
		// banners := templates.SingupBanner(errors.New("wrong username or password"))
		// return banners.Render(c.Context(), c.Response().BodyWriter())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	} else {
		sess_token, errSes := auth.GenerateToken(32)
		csrf_token, errCsrf := auth.GenerateToken(32)
		if errSes != nil || errCsrf != nil {
			// banners := templates.SingupBanner(errors.New("an error occurred while generating your authentication credentials"))
			// return banners.Render(c.Context(), c.Response().BodyWriter())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + errors.New("an error occurred while generating your authentication credentials").Error()})
		}
		err = queries.UpdateUserTokensLogin(ctx, db.UpdateUserTokensLoginParams{SessionToken: pgtype.Text{String: sess_token, Valid: true}, CsrfToken: pgtype.Text{String: csrf_token, Valid: true}, Username: username})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
			// banners := templates.SingupBanner(err)
			// return banners.Render(c.Context(), c.Response().BodyWriter())
		} else {
			c.Cookie(&fiber.Cookie{
				Name:     "session_token",
				Value:    sess_token,
				Expires:  time.Now().Add(24 * time.Hour),
				HTTPOnly: true,
			})
			c.Cookie(&fiber.Cookie{
				Name:     "csrf_token",
				Value:    csrf_token,
				Expires:  time.Now().Add(24 * time.Hour),
				HTTPOnly: false,
			})
			c.Set("HX-Redirect", "/")
			return c.SendStatus(fiber.StatusOK)
		}
	}
}

func HandleLogout(c *fiber.Ctx) error {
	err := auth.AuthorizePost(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
	} else {
		ctx := context.Background()
		sqlDb, err := auth.CreateNewDb()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
		}
		queries := db.New(sqlDb)
		st := c.Cookies("session_token", "")
		csrf := c.Cookies("csrf_token", "")
		queries.UpdateUserTokensLogout(ctx, db.UpdateUserTokensLogoutParams{SessionToken: pgtype.Text{String: st, Valid: true}, CsrfToken: pgtype.Text{String: csrf, Valid: true}})
		c.Set("HX-Redirect", "/signin")
		return c.SendStatus(fiber.StatusOK)
	}
}
