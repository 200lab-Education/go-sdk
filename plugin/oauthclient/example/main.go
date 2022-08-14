/*
 * @author           Viet Tran <viettranx@gmail.com>
 * @copyright       2019 Viet Tran <viettranx@gmail.com>
 * @license           Apache-2.0
 */

package main

import (
	"github.com/200Lab-Education/go-sdk"
	"github.com/200Lab-Education/go-sdk/plugin/oauthclient"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
)

var appClientConf = clientcredentials.Config{
	ClientID:     "200lab",
	ClientSecret: "secret-cannot-tell",
	Scopes:       []string{"root"},
	TokenURL:     "http://localhost:3000/oauth2/token",
}

func main() {
	service := goservice.New(
		goservice.WithName("demo"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(oauthclient.New("oauth", appClientConf)),
	)

	_ = service.Init()

	service.HTTPServer().AddHandler(func(engine *gin.Engine) {
		engine.GET("/login", func(context *gin.Context) {
			tc := service.MustGet("oauth").(oauthclient.TrustedClient)
			t, err := tc.PasswordCredentialsToken("admin", "Admin@2019")
			context.JSON(http.StatusOK, gin.H{"t": t, "err": err})
		})

		engine.GET("/introspect", func(context *gin.Context) {
			tc := service.MustGet("oauth").(oauthclient.TrustedClient)
			t, err := tc.Introspect("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOltdLCJlbWFpbCI6ImNvcmVAMjAwbGFiLmlvIiwiZXhwIjoxNTU2OTAyMzQ2LCJpYXQiOjE1NTQzMTAzNDYsImlzcyI6IiIsImp0aSI6IjI4ZmY1ZTI1LTY5NTgtNDhkZC05MGU5LWNiMDFjZjM5YjgwYSIsIm5iZiI6MTU1NDMxMDM0Niwic2NwIjpbIm9mZmxpbmUiXSwic3ViIjoiNWM5YjJhYzg3MzE3MWExN2Q1NzMyOGU4IiwidXNlcl9pZCI6IjVjOWIyYWM4NzMxNzFhMTdkNTczMjhlOCIsInVzZXJuYW1lIjoiYWRtaW4ifQ.YKrgfMyZ9Hs-RpUR6mTlENDcFAKrT2Pu7JrfE38bmSRFRMxleC48gEJArxy-1casJEQW_yW3Df9V-wKwGqK365VzV9T1aBdfxzOpU3GBRCq6YjaEx1d1SYPttZD02uOmVRu3zka-jm3225YkxYf6TXMRQ0xQbEg_RXUujE333sc")
			context.JSON(http.StatusOK, gin.H{"t": t, "err": err})
		})

		engine.GET("/create_user", func(context *gin.Context) {
			tc := service.MustGet("oauth").(oauthclient.TrustedClient)
			// t, err := tc.CreateUser(fmt.Sprintf("a%d", rand.Int()), "123456")
			username := "viet"
			email := "viet@200lab.io"
			password := "123456"

			t, err := tc.CreateUser(&oauthclient.OAuthUserCreate{
				Username: &username,
				Email:    &email,
				Password: &password,
			})

			context.JSON(http.StatusOK, gin.H{"t": t, "err": err})
		})

		engine.GET("/update_user", func(context *gin.Context) {
			tc := service.MustGet("oauth").(oauthclient.TrustedClient)
			// t, err := tc.CreateUser(fmt.Sprintf("a%d", rand.Int()), "123456")
			username := "viet"
			email := "viet@200lab.io"
			password := "123456"

			err := tc.UpdateUser("1", &oauthclient.OAuthUserUpdate{
				Username: &username,
				Email:    &email,
				Password: &password,
			})

			if err != nil {
				context.JSON(http.StatusUnprocessableEntity, gin.H{"err": err})
				return
			}

			context.JSON(http.StatusOK, gin.H{"success": true})

		})

		engine.PUT("/update_user", func(c *gin.Context) {
			tc := service.MustGet("oauth").(oauthclient.TrustedClient)
			username := "viet"
			email := "viet@200lab.io"
			password := "123456"

			err := tc.UpdateUser("1", &oauthclient.OAuthUserUpdate{
				Username: &username,
				Email:    &email,
				Password: &password,
			})

			if err != nil {
				c.JSON(http.StatusOK, gin.H{"err": err})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		engine.GET("/create_user/gmail", func(context *gin.Context) {
			tc := service.MustGet("oauth").(oauthclient.TrustedClient)
			t, err := tc.CreateUserWithEmail("viet@200lab.io")
			context.JSON(http.StatusOK, gin.H{"t": t, "err": err})
		})

		engine.GET("/create_user/facebook", func(context *gin.Context) {
			tc := service.MustGet("oauth").(oauthclient.TrustedClient)
			t, err := tc.CreateUserWithFacebook("12345", "test@200lab.io")
			context.JSON(http.StatusOK, gin.H{"t": t, "err": err})
		})

		engine.GET("/create_user/apple", func(context *gin.Context) {
			tc := service.MustGet("oauth").(oauthclient.TrustedClient)
			t, err := tc.CreateUserWithApple("12345", "test@200lab.io")
			context.JSON(http.StatusOK, gin.H{"t": t, "err": err})
		})

		engine.GET("/change_password", func(context *gin.Context) {
			tc := service.MustGet("oauth").(oauthclient.TrustedClient)
			err := tc.ChangePassword("5ca5b98f73171a20237053a3", "123456", "123456")
			context.JSON(http.StatusOK, gin.H{"err": err})
		})

		engine.GET("/delete_user", func(context *gin.Context) {
			tc := service.MustGet("oauth").(oauthclient.TrustedClient)
			err := tc.DeleteUser("5ca5b98f73171a20237053a3")
			context.JSON(http.StatusOK, gin.H{"err": err})
		})
	})

	_ = service.Start()
}
