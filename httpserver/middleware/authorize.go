package middleware

import (
	"context"
	"fmt"
	"github.com/200Lab-Education/go-sdk/logger"
	"github.com/200Lab-Education/go-sdk/plugin/oauthclient"
	"github.com/200Lab-Education/go-sdk/sdkcm"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ServiceContext interface {
	Logger(prefix string) logger.Logger
	Get(prefix string) (interface{}, bool)
	MustGet(prefix string) interface{}
}

type CurrentUserProvider interface {
	GetCurrentUser(ctx context.Context, oauthID string) (sdkcm.User, error)
	ServiceContext
}

type Tracker interface {
	TrackApiCall(userId uint32, url string) error
}

type Caching interface {
	GetCurrentUser(ctx context.Context, sig string) (sdkcm.Requester, error)
	WriteCurrentUser(ctx context.Context, sig string, u sdkcm.Requester) error
}

func Authorize(cup CurrentUserProvider, tracker Tracker, isRequired ...bool) gin.HandlerFunc {
	required := len(isRequired) == 0

	return func(c *gin.Context) {
		token := accessTokenFromRequest(c.Request)

		if token == "" {
			if required {
				panic(sdkcm.ErrUnauthorized(nil, sdkcm.ErrAccessTokenInvalid))
			} else {
				c.Set("current_user", guest{})
				c.Next()
				return
			}

		}

		tc := cup.MustGet("oauth").(oauthclient.TrustedClient)
		tokenInfo, err := tc.Introspect(c.Request.Context(), token)

		if err != nil {
			panic(sdkcm.ErrUnauthorized(err, sdkcm.ErrAccessTokenInactivated))
		}

		if !tokenInfo.Active {
			panic(sdkcm.ErrUnauthorized(sdkcm.ErrAccessTokenInactivated, sdkcm.ErrAccessTokenInactivated))
		}

		// Fetch user info from db
		u, err := cup.GetCurrentUser(c.Request.Context(), tokenInfo.UserId)

		if err != nil {
			panic(sdkcm.ErrUnauthorized(err, sdkcm.ErrUserNotFound))
		}

		go func(uid uint32, url string) {
			_ = tracker.TrackApiCall(u.UserID(), c.Request.URL.String())
		}(u.UserID(), c.Request.URL.String())

		c.Set("current_user", sdkcm.CurrentUser(tokenInfo, u))
	}
}

func AuthorizeWithCache(cup CurrentUserProvider, cache Caching, tracker Tracker, isRequired ...bool) gin.HandlerFunc {
	required := len(isRequired) == 0

	return func(c *gin.Context) {
		token := accessTokenFromRequest(c.Request)
		ctx := c.Request.Context()

		if token == "" {
			if required {
				panic(sdkcm.ErrUnauthorized(nil, sdkcm.ErrAccessTokenInvalid))
			} else {
				c.Set("current_user", guest{})
				c.Next()
				return
			}
		}

		tokenComps := strings.Split(token, ".")
		sig := tokenComps[len(tokenComps)-1]

		if cache != nil {
			if cacheUser, err := cache.GetCurrentUser(ctx, sig); err == nil {
				c.Set("current_user", cacheUser)

				go func(uid uint32, url string) {
					_ = tracker.TrackApiCall(cacheUser.UserID(), c.Request.URL.String())
				}(cacheUser.UserID(), c.Request.URL.String())
				return
			}
		}

		tc := cup.MustGet("oauth").(oauthclient.TrustedClient)
		tokenInfo, err := tc.Introspect(c.Request.Context(), token)

		if err != nil {
			panic(sdkcm.ErrUnauthorized(err, sdkcm.ErrAccessTokenInactivated))
		}

		if !tokenInfo.Active {
			panic(sdkcm.ErrUnauthorized(sdkcm.ErrAccessTokenInactivated, sdkcm.ErrAccessTokenInactivated))
		}

		// Fetch user info from db
		u, err := cup.GetCurrentUser(c.Request.Context(), tokenInfo.UserId)

		if err != nil {
			panic(sdkcm.ErrUnauthorized(err, sdkcm.ErrUserNotFound))
		}

		go func(uid uint32, url string) {
			_ = tracker.TrackApiCall(u.UserID(), c.Request.URL.String())
		}(u.UserID(), c.Request.URL.String())

		if cache != nil {
			_ = cache.WriteCurrentUser(ctx, sig, sdkcm.CurrentUser(tokenInfo, u))
		}

		c.Set("current_user", sdkcm.CurrentUser(tokenInfo, u))
	}
}

func RequireRoles(roles ...fmt.Stringer) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, ok := c.Get("current_user")

		if !ok {
			panic(sdkcm.ErrUnauthorized(sdkcm.ErrNoPermission, sdkcm.ErrNoPermission))
		}

		requester := r.(sdkcm.Requester)
		reqRole := sdkcm.ParseSystemRole(requester.GetSystemRole())

		for _, v := range roles {
			if v.String() == reqRole.String() {
				c.Next()
				return
			}
		}

		panic(sdkcm.ErrUnauthorized(nil, sdkcm.ErrNoPermission))
	}
}

func accessTokenFromRequest(req *http.Request) string {
	// According to https://tools.ietf.org/html/rfc6750 you can pass tokens through:
	// - Form-Encoded Body Parameter. Recommended, more likely to appear. e.g.: Authorization: Bearer mytoken123
	// - URI Query Parameter e.g. access_token=mytoken123

	auth := req.Header.Get("Authorization")
	split := strings.SplitN(auth, " ", 2)
	if len(split) != 2 || !strings.EqualFold(split[0], "bearer") {
		// Nothing in Authorization header, try access_token
		// Empty string returned if there's no such parameter
		if err := req.ParseMultipartForm(1 << 20); err != nil && err != http.ErrNotMultipart {
			return ""
		}
		return req.Form.Get("access_token")
	}

	return split[1]
}

type guest struct{}

func (g guest) OAuthID() string       { return "" }
func (g guest) UserID() uint32        { return 0 }
func (g guest) GetSystemRole() string { return sdkcm.SysRoleGuest.String() }
