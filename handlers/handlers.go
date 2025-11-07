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

func HandleSearchGateway(c *fiber.Ctx) error {
	searchType := c.FormValue("search-type")
	if searchType == "org" {
		return HandleOrgSearch(c)
	} else {
		return HandleUserSearch(c)
	}
}

func HandleUserSearch(c *fiber.Ctx) error {
	uniqueSearchId, _ := auth.GenerateToken(16)
	phMonitor, ghClient := configurePosthogAndGitHub()
	value := c.FormValue("search-input")
	start := time.Now()
	stats, err := ghClient.GetUserStats(value)
	end := time.Now()
	latency := end.Sub(start).Milliseconds()
	c.Set("Content-Type", "text/html")
	if err != nil {
		errPh := phMonitor.SendEvent(uniqueSearchId, "ghSearch", "userSearch", latency, true, err.Error())
		if errPh != nil {
			log.Println("PostHog failed to record event")
		}
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	errPh := phMonitor.SendEvent(uniqueSearchId, "ghSearch", "userSearch", latency, false, "")
	if errPh != nil {
		log.Println("PostHog failed to record event")
	}
	return templates.UserStatsDisplay(*stats).Render(c.Context(), c.Response().BodyWriter())
}

func HandleOrgSearch(c *fiber.Ctx) error {
	_, err := auth.AuthorizePost(c)
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	uniqueSearchId, _ := auth.GenerateToken(16)
	value := c.FormValue("search-input")
	phMonitor, ghClient := configurePosthogAndGitHub()
	start := time.Now()
	stats, err := ghClient.GetOrgStats(value)
	end := time.Now()
	latency := end.Sub(start).Milliseconds()
	c.Set("Content-Type", "text/html")
	if err != nil {
		errPh := phMonitor.SendEvent(uniqueSearchId, "ghSearch", "orgSearch", latency, true, err.Error())
		if errPh != nil {
			log.Println("PostHog failed to record event")
		}
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	errPh := phMonitor.SendEvent(uniqueSearchId, "ghSearch", "orgSearch", latency, false, "")
	if errPh != nil {
		log.Println("PostHog failed to record event")
	}
	return templates.OrgStatsDisplay(*stats).Render(c.Context(), c.Response().BodyWriter())
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

func SearchRoute(c *fiber.Ctx) error {
	err := auth.AuthorizeGet(c)
	c.Set("Content-Type", "text/html")
	return templates.SearchInterface(err == nil).Render(c.Context(), c.Response().BodyWriter())
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
		phMonitor.SendEvent(username, "userAuth", "signUp", time.Since(start).Milliseconds(), true, err.Error())
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	queries := db.New(sqlDb)
	_, err = queries.GetUser(ctx, username)
	secondPoint := time.Now()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			hashed_psw, err := auth.HashPassword(password)
			thirdPoint := time.Now()
			if err != nil {
				phMonitor.SendEvent(username, "userAuth", "signUp", thirdPoint.Sub(start).Milliseconds(), true, err.Error())
				return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
			}
			_, err = queries.CreateUser(ctx, db.CreateUserParams{Username: username, HashedPassword: hashed_psw})
			end := time.Now()
			if err != nil {
				phMonitor.SendEvent(username, "userAuth", "signUp", end.Sub(start).Milliseconds(), true, err.Error())
				return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
			} else {
				phMonitor.SendEvent(username, "userAuth", "signUp", end.Sub(start).Milliseconds(), false, "")
				return templates.StatusBanner(nil).Render(c.Context(), c.Response().BodyWriter())
			}
		} else {
			phMonitor.SendEvent(username, "userAuth", "signUp", secondPoint.Sub(start).Milliseconds(), true, err.Error())
			return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
		}
	} else {
		phMonitor.SendEvent(username, "userAuth", "signUp", secondPoint.Sub(start).Milliseconds(), true, "user already exists")
		return templates.StatusBanner(errors.New("user already exists")).Render(c.Context(), c.Response().BodyWriter())
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
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	queries := db.New(sqlDb)
	user, err := queries.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			phMonitor.SendEvent(username, "userAuth", "login", time.Since(start).Milliseconds(), true, err.Error())
			return templates.StatusBanner(errors.New("there is no user with this username")).Render(c.Context(), c.Response().BodyWriter())
		} else {
			phMonitor.SendEvent(username, "userAuth", "login", time.Since(start).Milliseconds(), true, err.Error())
			return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
		}
	}
	if !auth.CompareHashToPassword(password, user.HashedPassword) {
		return templates.StatusBanner(errors.New("wrong username or password")).Render(c.Context(), c.Response().BodyWriter())
	} else {
		sess_token, errSes := auth.GenerateToken(32)
		csrf_token, errCsrf := auth.GenerateToken(32)
		if errSes != nil || errCsrf != nil {
			phMonitor.SendEvent(username, "userAuth", "login", time.Since(start).Milliseconds(), true, "an error occurred while generating auth credentials")
			return templates.StatusBanner(errors.New("an error occurred while generating your authentication credentials")).Render(c.Context(), c.Response().BodyWriter())
		}
		err = queries.UpdateUserTokensLogin(ctx, db.UpdateUserTokensLoginParams{SessionToken: pgtype.Text{String: sess_token, Valid: true}, CsrfToken: pgtype.Text{String: csrf_token, Valid: true}, Username: username})
		if err != nil {
			phMonitor.SendEvent(username, "userAuth", "login", time.Since(start).Milliseconds(), true, err.Error())
			return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
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
		return c.Status(500).JSON(fiber.Map{"message": "An error occurred: " + err.Error()})
	} else {
		ctx := context.Background()
		sqlDb, err := auth.CreateNewDb()
		if err != nil {
			phMonitor.SendEvent(user.Username, "userAuth", "logout", time.Since(start).Milliseconds(), true, err.Error())
			return c.Status(500).JSON(fiber.Map{"message": "An error occurred: " + err.Error()})
		}
		queries := db.New(sqlDb)
		st := c.Cookies("session_token", "")
		csrf := c.Cookies("csrf_token", "")
		err = queries.UpdateUserTokensLogout(ctx, db.UpdateUserTokensLogoutParams{SessionToken: pgtype.Text{String: st, Valid: true}, CsrfToken: pgtype.Text{String: csrf, Valid: true}})
		if err != nil {
			phMonitor.SendEvent(user.Username, "userAuth", "logout", time.Since(start).Milliseconds(), true, err.Error())
			return c.Status(500).JSON(fiber.Map{"message": "An error occurred: " + err.Error()})
		}
		phMonitor.SendEvent(user.Username, "userAuth", "logout", time.Since(start).Milliseconds(), false, "")
		c.Set("HX-Redirect", "/signin")
		return c.SendStatus(fiber.StatusOK)
	}
}
