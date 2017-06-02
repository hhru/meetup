package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/hhru/meetup/jira"
	"github.com/hhru/meetup/talk"
	"github.com/hhru/meetup/user"
	"net/http"
	"time"
)

type jwtCustomClaims struct {
	Temporary *oauth.Credentials `json:"temporary"`
	Permanent *oauth.Credentials `json:"permanent"`
	user.User
	jwt.StandardClaims
}

type CustomContext struct {
	echo.Context
	unauthorizedClient *jira.Client
	client             *jira.AuthorizedClient
}

func (c *CustomContext) SaveToken(claims *jwtCustomClaims) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	c.SetCookie(
		&http.Cookie{
			Name:    "token",
			Value:   t,
			Expires: time.Now().Add(time.Hour * 72),
			Path:    "/",
		},
	)

	return nil
}

func AuthorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if _, err := c.Request().Cookie("token"); err != nil {
			return c.Redirect(http.StatusFound, "/oauth/login")
		}
		return next(c)
	}
}

func ClientMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*jwtCustomClaims)

		myCtx := c.(CustomContext)
		myCtx.client = jira.NewAuthorizedClient(myCtx.unauthorizedClient, claims.Permanent)

		return next(myCtx)
	}
}

func bindApp() echo.MiddlewareFunc {
	jiraClient, _ := jira.NewClient("jira.hh.ru")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			myCtx := CustomContext{c, jiraClient, nil}
			return next(myCtx)
		}
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(bindApp())

	e.Debug = true

	e.GET("/oauth/login", func(c echo.Context) error {
		context := c.(CustomContext)

		e.Logger.Debug(context.Request().Host)
		callback := "http://" + context.Request().Host + "/oauth/callback"
		tempCred, err := context.unauthorizedClient.RequestTemporaryCredentials(nil, callback, nil)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error getting temp cred, "+err.Error())
		}

		claims := &jwtCustomClaims{
			tempCred,
			nil,
			user.User{Name: "", DisplayName: "", Avatar: ""},
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		}

		if err := context.SaveToken(claims); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error saving session , "+err.Error())
		}

		return context.Redirect(http.StatusFound, context.unauthorizedClient.AuthorizationURL(tempCred, nil))
	})

	g := e.Group("")

	g.Use(AuthorizationMiddleware)
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte("secret"),
		TokenLookup: "cookie:token",
		Claims:      &jwtCustomClaims{},
	}))
	g.Use(ClientMiddleWare)

	g.Static("/", "dist")

	g.GET("/oauth/callback", func(c echo.Context) error {
		context := c.(CustomContext)
		storage := context.Get("user").(*jwt.Token)
		claims := storage.Claims.(*jwtCustomClaims)
		tempCred := claims.Temporary
		if tempCred == nil || tempCred.Token != context.FormValue("oauth_token") {
			return echo.NewHTTPError(http.StatusInternalServerError, "Unknown oauth_token.")
		}
		tokenCred, _, err := context.unauthorizedClient.RequestToken(nil, tempCred, context.FormValue("oauth_verifier"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error getting request token, "+err.Error())
		}
		claims.Temporary = nil
		claims.Permanent = tokenCred

		myself, err := user.GetMyself(jira.NewAuthorizedClient(context.unauthorizedClient, claims.Permanent))

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching storage information, "+err.Error())
		}

		claims.User = *myself

		if err := context.SaveToken(claims); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error saving session, "+err.Error())
		}

		return context.Redirect(http.StatusFound, "/")
	})

	g.GET("/api/me", func(c echo.Context) error {
		context := c.(CustomContext)
		storage := context.Get("user").(*jwt.Token)
		claims := storage.Claims.(*jwtCustomClaims)
		if claims.User.Name == "" {
			return echo.NewHTTPError(http.StatusForbidden, "Unauthorized")
		}
		return c.JSON(http.StatusOK, claims.User)
	})

	g.GET("/api/talks", func(c echo.Context) error {
		context := c.(CustomContext)
		storage := context.Get("user").(*jwt.Token)
		claims := storage.Claims.(*jwtCustomClaims)

		result, err := talk.GetAllTalks(&claims.User, context.client)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, result)
	})

	g.GET("/api/:key/like", func(c echo.Context) error {
		context := c.(CustomContext)

		key := context.Param("key")
		storage := context.Get("user").(*jwt.Token)
		claims := storage.Claims.(*jwtCustomClaims)

		result, err := talk.Like(&claims.User, context.client, key)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, result)
	})

	g.GET("/api/:key/dislike", func(c echo.Context) error {
		context := c.(CustomContext)

		key := context.Param("key")
		storage := context.Get("user").(*jwt.Token)
		claims := storage.Claims.(*jwtCustomClaims)

		result, err := talk.Dislike(&claims.User, context.client, key)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, result)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
