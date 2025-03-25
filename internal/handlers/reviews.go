package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"movie-api/internal/models"
)

type ReviewResponse struct {
	ID        uint   `json:"id"`
	UserName  string `json:"user_name"`
	MovieTitle string `json:"movie_title"`
	Text      string `json:"text"`
}

type CreateReviewRequest struct {
	MovieID uint   `json:"movie_id" binding:"required"`
	Text    string `json:"text" binding:"required"`
}

func GetReviews(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reviews []ReviewResponse

		result := db.Model(&models.Review{}).
			Select("reviews.id, users.username as user_name, movies.title as movie_title, reviews.text").
			Joins("JOIN users ON users.id = reviews.user_id").
			Joins("JOIN movies ON movies.id = reviews.movie_id").
			Scan(&reviews)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
			return
		}

		c.JSON(http.StatusOK, reviews)
	}
}

func GetReviewDetails(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
			return
		}

		var review ReviewResponse
		result := db.Model(&models.Review{}).
			Select("reviews.id, users.username as user_name, movies.title as movie_title, reviews.text").
			Joins("JOIN users ON users.id = reviews.user_id").
			Joins("JOIN movies ON movies.id = reviews.movie_id").
			Where("reviews.id = ?", id).
			First(&review)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
			return
		}

		c.JSON(http.StatusOK, review)
	}
}

func CreateReview(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.MustGet("userID").(uint)
		var req CreateReviewRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Check if movie exists
		var movie models.Movie
		if result := db.First(&movie, req.MovieID); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie not found"})
			return
		}

		review := models.Review{
			MovieID: req.MovieID,
			UserID:  userID,
			Text:    &req.Text,
		}

		if result := db.Create(&review); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":         review.ID,
			"message":    "Review created successfully",
		})
	}
}