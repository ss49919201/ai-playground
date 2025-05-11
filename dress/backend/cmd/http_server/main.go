package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	server := NewServer()

	if err := http.ListenAndServe(":8080", server); err != nil {
		slog.Error("failed to run http server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

type Sever struct {
	ginEngine *gin.Engine
}

func NewServer() *Sever {
	engine := gin.New()

	apiGroup := engine.Group("/api")
	apiV1Group := apiGroup.Group("/v1")

	// signup
	apiV1Group.POST("/signup", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "signup",
		})
	})

	// login
	apiV1Group.POST("/login", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "login",
		})
	})

	// logout
	apiV1Group.POST("/logout", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "logout",
		})
	})

	apiGroupWithAuth := apiV1Group.Group("/")
	apiGroupWithAuth.Use(
		// TODO: implement auth middleware
		func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "auth",
			})
		},
	)

	// threads
	threads := apiGroupWithAuth.Group("/threads")
	threads.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "threads",
		})
	})
	threads.POST("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"id": "1",
		})
	})
	threads.GET("/:id/posts", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"id": c.Param("id"),
			"posts": []string{
				"post1",
				"post2",
			},
		})
	})
	threads.POST("/:id/posts", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"id":   c.Param("id"),
			"post": "post",
		})
	})

	return &Sever{
		ginEngine: engine,
	}
}

func (s *Sever) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.ginEngine.ServeHTTP(w, r)
}
