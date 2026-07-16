package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ioarintoko/sv-backend-test/internal/config"
	"github.com/ioarintoko/sv-backend-test/internal/database"
	"github.com/ioarintoko/sv-backend-test/internal/handlers"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)
	defer db.Close()

	postHandler := handlers.NewPostHandler(db)

	router := gin.Default()

	article := router.Group("/article")
	{
		article.POST("/", postHandler.CreatePost)
		article.GET("/:id/:offset", postHandler.GetPosts)
		article.GET("/:id", postHandler.GetPostByID)

		article.PUT("/:id", postHandler.UpdatePost)
		article.PATCH("/:id", postHandler.UpdatePost)

		article.DELETE("/:id", postHandler.DeletePost)
	}

	log.Printf("Server running on port %s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}