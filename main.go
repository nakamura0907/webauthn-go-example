package main

import (
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	inmemory     = make(map[string]*webauthn.SessionData)
	inmemoryUser = make(map[string]*User)
)

func main() {
	waConfig := &webauthn.Config{
		RPDisplayName: "WebAuthn Go Example",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:8080"},
	}

	webAuthn, err := webauthn.New(waConfig)
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/webauthn/register/begin", func(c echo.Context) error {
		// Find or create the new user
		user := &User{
			ID:          []byte("1234"),
			Name:        "test",
			DisplayName: "Test User",
		}
		options, session, err := webAuthn.BeginRegistration(user)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		// store the sessionData values
		inmemory["1234"] = session
		inmemoryUser["1234"] = user

		return c.JSON(http.StatusOK, options)
		// options.publicKey contain our registration options
	})
	e.POST("/webauthn/register/finish", func(c echo.Context) error {
		// Get the user
		user := inmemoryUser["1234"]

		// Get the session data stored from the function above
		session := inmemory["1234"]

		credential, err := webAuthn.FinishRegistration(user, *session, c.Request())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		// If creation was successful, store the credential object
		// Pseudocode to add the user credential.
		user.Credentials = append(user.Credentials, *credential)
		inmemoryUser["1234"] = user

		return c.JSON(http.StatusOK, "Registration complete!")
	})

	e.POST("/webauthn/login/begin", func(c echo.Context) error {
		// Find the user
		user := inmemoryUser["1234"]

		options, session, err := webAuthn.BeginLogin(user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		// store the sessionData values
		inmemory["1234"] = session

		return c.JSON(http.StatusOK, options)
		// options.publicKey contain our registration options
	})
	e.POST("/webauthn/login/finish", func(c echo.Context) error {
		// Get the user
		user := inmemoryUser["1234"]

		// Get the session data stored from the function above
		session := inmemory["1234"]

		_, err := webAuthn.FinishLogin(user, *session, c.Request())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		// Handle credential.Authenticator.CloneWarning

		// If login was successful, update the credential object
		// Pseudocode to update the user credential.
		// TODO: Update the user credential

		return c.JSON(http.StatusOK, "Login complete!")
	})

	e.Static("/static/js/", "./static/js")
	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
