// @title Movie API
// @version 1.0
// @description This is a sample movie rating and review API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@movieapi.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @host localhost:8000
// @BasePath /api
package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "movie-api/docs"
	"movie-api/internal/auth"
	"movie-api/internal/database"
	"movie-api/internal/handlers"
)

func main() {
	// Initialize database
	db := database.InitDB()
	
	// Setup router
	r := gin.Default()

	// Public routes
	r.POST("/api/token", auth.LoginHandler(db))
	r.POST("/api/users", handlers.CreateUser(db)) 
	r.GET("/api/movies", handlers.GetMovies(db))
	r.GET("/api/movies/:id/", handlers.GetMovieDetails(db))
	r.GET("/api/reviews", handlers.GetReviews(db))
	r.GET("/api/reviews/:id/", handlers.GetReviewDetails(db))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Protected routes
	authGroup := r.Group("/")
	authGroup.Use(auth.JWTAuthMiddleware(db))
	{
		authGroup.POST("/api/reviews", handlers.CreateReview(db))
		authGroup.POST("/api/movies", handlers.CreateMovie(db))
	}

	r.Run(":8000")
}