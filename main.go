package main

import (
	"fmt"
	_ "github.com/appleboy/gin-jwt/v2"
	"github.com/drbear95/gonotter-server/auth"
	"github.com/drbear95/gonotter-server/ghql"
	_ "github.com/drbear95/gonotter-server/persistance"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"golang.org/x/net/context"
	"log"
	_ "net/http"
	"time"
)

func executeGraphql() gin.HandlerFunc {
	h := handler.New(&handler.Config{
		Schema:     ghql.GetSchema(),
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	return func(c *gin.Context) {
		authDetails, err := auth.ExtractTokenMetadata(c.Request)

		if err != nil {
			log.Fatal(err)
		}

		ctx := context.WithValue(c.Request.Context(), "auth_details", authDetails)

		h.ContextHandler(ctx, c.Writer, c.Request)
	}
}

func getLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

func main() {
	r := gin.New()
	r.Use(getLogger())
	r.Use(gin.Recovery())

	apiV1 := r.Group("/api/v1")

	apiV1.POST("/signIn", auth.SignIn)
	apiV1.POST("/signUp", auth.SignUp)

	apiV1.Use(auth.TokenAuthMiddleware())
	{
		apiV1.POST("", executeGraphql())
	}

	log.Fatal(r.Run(":8080"))
}
