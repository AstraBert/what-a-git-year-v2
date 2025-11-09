package auth

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthorizeGetFail(t *testing.T) {
	if _, ok := os.LookupEnv("POSTGRES_CONNECTION_STRING"); !ok {
		t.Skip()
	} else {
		c := fiber.Ctx{}
		c.Cookie(&fiber.Cookie{Name: "session_token", Value: "noSession"})
		err := AuthorizeGet(&c)
		if err == nil {
			t.Error("Expected an error, got none")
		}
	}
}

func TestAuthorizePostFail(t *testing.T) {
	if _, ok := os.LookupEnv("POSTGRES_CONNECTION_STRING"); !ok {
		t.Skip()
	} else {
		c := fiber.Ctx{}
		c.Cookie(&fiber.Cookie{Name: "session_token", Value: "noSession"})
		c.Cookie(&fiber.Cookie{Name: "csrf_token", Value: "noCSRF"})
		_, err := AuthorizePost(&c)
		if err == nil {
			t.Error("Expected an error, got none")
		}
	}
}
