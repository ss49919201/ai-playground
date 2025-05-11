package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ss49919201/ai-kata/dress/backend/auth"
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
		if _, err := auth.Signup(
			c.PostForm("email"),
			c.PostForm("password"),
		); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"msg": "ok",
		})
	})

	// signin
	apiV1Group.POST("/signin", func(c *gin.Context) {
		if _, err := auth.Signin(
			c.PostForm("email"),
			c.PostForm("password"),
		); err != nil {
			c.JSON(500, gin.H{
				"error": "failed to signin",
			})
			return
		}
		c.JSON(200, gin.H{
			"msg": "ok",
		})
	})

	// signout
	apiV1Group.POST("/signout", func(c *gin.Context) {
		if err := auth.Signout(
			c.PostForm("token"),
		); err != nil {
			c.JSON(500, gin.H{
				"error": "failed to signout",
			})
			return
		}
		c.JSON(200, gin.H{
			"msg": "ok",
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
	threadsApi := apiGroupWithAuth.Group("/threads")
	threadsApi.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "ok",
		})
	})
	threadsApi.POST("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"id": "1",
		})
	})
	threadsApi.GET("/:id/posts", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"id": c.Param("id"),
			"posts": []string{
				"post1",
				"post2",
			},
		})
	})
	threadsApi.POST("/:id/posts", func(c *gin.Context) {
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
