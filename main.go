package main

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/vecmezoni/gomeet/jira"
	"time"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/vecmezoni/gomeet/meetup"
)

type jwtCustomClaims struct {
	Temporary *oauth.Credentials `json:"temporary"`
	Permanent *oauth.Credentials `json:"permanent"`
  Name string `json:"name"`
  DisplayName string `json:"display_name"`
  Avatar string `json:"avatar"`
	jwt.StandardClaims
}


type CustomContext struct {
	echo.Context
	unauthorizedClient *jira.Client
	client *jira.AuthorizedClient
}

func (c *CustomContext) SaveToken(claims *jwtCustomClaims) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	c.SetCookie(
		&http.Cookie{
			Name: "token",
			Value: t,
			Expires: time.Now().Add(time.Hour * 72),
      Path: "/",
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
      "",
      "",
      "",
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

  loaded := false

  cache := meetup.Init()

	g.Use(AuthorizationMiddleware)
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
		TokenLookup: "cookie:token",
		Claims:     &jwtCustomClaims{},
	}))
	g.Use(ClientMiddleWare)

  g.GET("/", func(c echo.Context) error {
    return c.HTML(http.StatusOK, "HUY!")
  })

	g.GET("/oauth/callback", func(c echo.Context) error {
		context := c.(CustomContext)
		user := context.Get("user").(*jwt.Token)
		claims := user.Claims.(*jwtCustomClaims)
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

    myself, err := jira.NewAuthorizedClient(context.unauthorizedClient, claims.Permanent).GetMyself()

    if err != nil {
      return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching user information, "+err.Error())
    }

    claims.Name = myself.Name
    claims.DisplayName = myself.DisplayName
    claims.Avatar = myself.AvatarUrls.Big

		if err := context.SaveToken(claims); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error saving session, "+err.Error())
		}

		return context.Redirect(http.StatusFound, "/")
	})

  g.GET("/api/me", func(c echo.Context) error {
    context := c.(CustomContext)
    user := context.Get("user").(*jwt.Token)
    claims := user.Claims.(*jwtCustomClaims)

    result := new(meetup.User)

    result.Name = claims.Name
    result.DisplayName = claims.DisplayName
    result.Avatar = claims.Avatar

    return c.JSON(http.StatusOK, result)
  })

	g.GET("/api/talks", func(c echo.Context) error {
    context := c.(CustomContext)

    if !loaded {
      talks, err := context.client.GetTalks(`project = "R&D :: Meetups"`)
      if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err)
      }
      cache.Fill(meetup.TalksFromJira(talks))
      loaded = true
    }

    result, err := cache.GetTalks()

    if err != nil {
      return echo.NewHTTPError(http.StatusInternalServerError, err)
    }

		return c.JSON(http.StatusOK, result)
	})

  g.GET("/api/:key/like", func(c echo.Context) error {
    context := c.(CustomContext)

    key := context.Param("key")

    err := context.client.Like(key)

    if err != nil {
      return echo.NewHTTPError(http.StatusInternalServerError, err)
    }

    updatedJiraTalk, err := context.client.GetTalk(key)

    talk := meetup.TalkFromJira(updatedJiraTalk)

    err = cache.Update(talk)

    if err != nil {
      return echo.NewHTTPError(http.StatusInternalServerError, err)
    }

    cachedTalk, err := cache.GetTalk(key)

    if err != nil {
      return echo.NewHTTPError(http.StatusInternalServerError, err)
    }

    return c.JSON(http.StatusOK, cachedTalk)
  })

	e.Logger.Fatal(e.Start(":1323"))
}
