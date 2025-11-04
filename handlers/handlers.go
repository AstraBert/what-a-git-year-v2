package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/AstraBert/what-a-git-year-v2/auth"
	"github.com/AstraBert/what-a-git-year-v2/db"
	"github.com/AstraBert/what-a-git-year-v2/gh"
	"github.com/AstraBert/what-a-git-year-v2/monitoring"
	"github.com/AstraBert/what-a-git-year-v2/templates"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func configurePosthog() *monitoring.PosthogMonitor {
	phApiKey := os.Getenv("POSTHOG_API_KEY")
	phEndpoint := os.Getenv("POSTHOG_ENDPOINT")
	return monitoring.NewPosthogMonitor(phApiKey, phEndpoint)

}

func configurePosthogAndGitHub() (*monitoring.PosthogMonitor, *gh.GitYearClient) {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	ghClient := gh.NewGitYearClient(token)
	phApiKey := os.Getenv("POSTHOG_API_KEY")
	phEndpoint := os.Getenv("POSTHOG_ENDPOINT")
	phMonitor := monitoring.NewPosthogMonitor(phApiKey, phEndpoint)
	return phMonitor, ghClient
}

func HandleUserSearch(c *fiber.Ctx) error {
	uniqueSearchId, _ := auth.GenerateToken(16)
	user := c.FormValue("user")
	phMonitor, ghClient := configurePosthogAndGitHub()
	start := time.Now()
	stats, err := ghClient.GetUserStats(user)
	end := time.Now()
	latency := end.Sub(start).Milliseconds()
	if err != nil {
		errPh := phMonitor.SendEvent(uniqueSearchId, "ghSearch", "userSearch", latency, true, err.Error())
		if errPh != nil {
			log.Println("PostHog failed to record event")
		}
		return c.SendStatus(500)
	}
	errPh := phMonitor.SendEvent(uniqueSearchId, "ghSearch", "userSearch", latency, false, "")
	if errPh != nil {
		log.Println("PostHog failed to record event")
	}
	return c.Status(200).JSON(fiber.Map{"stats": stats})
}

func HandleOrgSearch(c *fiber.Ctx) error {
	_, err := auth.AuthorizePost(c)
	if err != nil {
		return c.SendStatus(401)
	}
	uniqueSearchId, _ := auth.GenerateToken(16)
	organization := c.FormValue("organization")
	phMonitor, ghClient := configurePosthogAndGitHub()
	start := time.Now()
	stats, err := ghClient.GetOrgStats(organization)
	end := time.Now()
	latency := end.Sub(start).Milliseconds()
	if err != nil {
		errPh := phMonitor.SendEvent(uniqueSearchId, "ghSearch", "userSearch", latency, true, err.Error())
		if errPh != nil {
			log.Println("PostHog failed to record event")
		}
		return c.SendStatus(500)
	}
	errPh := phMonitor.SendEvent(uniqueSearchId, "ghSearch", "userSearch", latency, false, "")
	if errPh != nil {
		log.Println("PostHog failed to record event")
	}
	return c.Status(200).JSON(fiber.Map{"stats": stats})
}

func HomeRoute(c *fiber.Ctx) error {
	err := auth.AuthorizeGet(c)
	c.Set("Content-Type", "text/html")
	return templates.Home(err == nil).Render(c.Context(), c.Response().BodyWriter())
}

func LoginRoute(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html")
	return templates.SignIn().Render(c.Context(), c.Response().BodyWriter())
}

func SignUpRoute(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html")
	return templates.SignUp().Render(c.Context(), c.Response().BodyWriter())
}

func PageDoesNotExistRoute(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html")
	return templates.Page404().Render(c.Context(), c.Response().BodyWriter())
}

func HandleSignUp(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	passwordR := c.FormValue("passwordRepeat")
	if password != passwordR {
		return c.SendStatus(400)
	}
	ctx := context.Background()
	phMonitor := configurePosthog()
	start := time.Now()
	sqlDb, err := auth.CreateNewDb()
	if err != nil {
		// banners := templates.SingupBanner(err)
		// return banners.Render(c.Context(), c.Response().BodyWriter())
		phMonitor.SendEvent(username, "userAuth", "signUp", time.Since(start).Milliseconds(), true, err.Error())
		return c.SendStatus(500)
	}
	queries := db.New(sqlDb)
	_, err = queries.GetUser(ctx, username)
	secondPoint := time.Now()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			hashed_psw, err := auth.HashPassword(password)
			thirdPoint := time.Now()
			if err != nil {
				// banners := templates.SingupBanner(err)
				// return banners.Render(c.Context(), c.Response().BodyWriter())
				phMonitor.SendEvent(username, "userAuth", "signUp", thirdPoint.Sub(start).Milliseconds(), true, err.Error())
				return c.SendStatus(500)
			}
			_, err = queries.CreateUser(ctx, db.CreateUserParams{Username: username, HashedPassword: hashed_psw})
			end := time.Now()
			if err != nil {
				// banners := templates.SingupBanner(err)
				// return banners.Render(c.Context(), c.Response().BodyWriter())
				phMonitor.SendEvent(username, "userAuth", "signUp", end.Sub(start).Milliseconds(), true, err.Error())
				return c.SendStatus(500)
			} else {
				// banners := templates.SingupBanner(nil)
				// return banners.Render(c.Context(), c.Response().BodyWriter())
				phMonitor.SendEvent(username, "userAuth", "signUp", end.Sub(start).Milliseconds(), false, "")
				return c.SendStatus(200)
			}
		} else {
			// banners := templates.SingupBanner(err)
			// return banners.Render(c.Context(), c.Response().BodyWriter())
			phMonitor.SendEvent(username, "userAuth", "signUp", secondPoint.Sub(start).Milliseconds(), true, err.Error())
			return c.SendStatus(500)
		}
	} else {
		// banners := templates.SingupBanner(errors.New("user already exists"))
		// return banners.Render(c.Context(), c.Response().BodyWriter())
		phMonitor.SendEvent(username, "userAuth", "signUp", secondPoint.Sub(start).Milliseconds(), true, "user already exists")
		return c.SendStatus(fiber.StatusConflict)
	}
}

func HandleLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	ctx := context.Background()
	phMonitor := configurePosthog()
	start := time.Now()
	sqlDb, err := auth.CreateNewDb()
	if err != nil {
		phMonitor.SendEvent(username, "userAuth", "login", time.Since(start).Milliseconds(), true, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
	}
	queries := db.New(sqlDb)
	user, err := queries.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// banners := templates.SingupBanner(errors.New("there is no user with this username"))
			// return banners.Render(c.Context(), c.Response().BodyWriter())
			phMonitor.SendEvent(username, "userAuth", "login", time.Since(start).Milliseconds(), true, err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
		} else {
			// banners := templates.SingupBanner(err)
			// return banners.Render(c.Context(), c.Response().BodyWriter())
			phMonitor.SendEvent(username, "userAuth", "login", time.Since(start).Milliseconds(), true, err.Error())
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
			phMonitor.SendEvent(username, "userAuth", "login", time.Since(start).Milliseconds(), true, "an error occurred while generating auth credentials")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + errors.New("an error occurred while generating your authentication credentials").Error()})
		}
		err = queries.UpdateUserTokensLogin(ctx, db.UpdateUserTokensLoginParams{SessionToken: pgtype.Text{String: sess_token, Valid: true}, CsrfToken: pgtype.Text{String: csrf_token, Valid: true}, Username: username})
		if err != nil {
			phMonitor.SendEvent(username, "userAuth", "login", time.Since(start).Milliseconds(), true, err.Error())
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
			phMonitor.SendEvent(username, "userAuth", "login", time.Since(start).Milliseconds(), false, "")
			return c.SendStatus(fiber.StatusOK)
		}
	}
}

func HandleLogout(c *fiber.Ctx) error {
	phMonitor := configurePosthog()
	start := time.Now()
	user, err := auth.AuthorizePost(c)
	if err != nil {
		phMonitor.SendEvent(user.Username, "userAuth", "logout", time.Since(start).Milliseconds(), true, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
	} else {
		ctx := context.Background()
		sqlDb, err := auth.CreateNewDb()
		if err != nil {
			phMonitor.SendEvent(user.Username, "userAuth", "logout", time.Since(start).Milliseconds(), true, err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
		}
		queries := db.New(sqlDb)
		st := c.Cookies("session_token", "")
		csrf := c.Cookies("csrf_token", "")
		err = queries.UpdateUserTokensLogout(ctx, db.UpdateUserTokensLogoutParams{SessionToken: pgtype.Text{String: st, Valid: true}, CsrfToken: pgtype.Text{String: csrf, Valid: true}})
		if err != nil {
			phMonitor.SendEvent(user.Username, "userAuth", "logout", time.Since(start).Milliseconds(), true, err.Error())
			c.SendStatus(500)
		}
		phMonitor.SendEvent(user.Username, "userAuth", "logout", time.Since(start).Milliseconds(), false, "")
		c.Set("HX-Redirect", "/signin")
		return c.SendStatus(fiber.StatusOK)
	}
}
