package handlers

import (
	"fmt"
	"movie-api/internal/models"
	"movie-api/internal/auth"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewResponse struct {
	ID        uint   `json:"id"`
	UserName  string `json:"user_name"`
	MovieTitle string `json:"movie_title"`
	Rating	  float64	 `json:"user_rating"`
	Text      string `json:"text"`
}

type CreateReviewRequest struct {
	MovieID uint   `json:"movie_id" binding:"required"`
	Rating	float64	`json:"user_rating,omitempty"`
	Text    string `json:"text,omitempty"`
}

func (r *CreateReviewRequest) Validate() error {
    if r.Rating == 0 && r.Text == "" {
        return fmt.Errorf("either rating or text must be provided, but not both empty")
    }
    return nil
}

// GetReviews godoc
// @Summary Get all reviews
// @Description Get a list of all reviews with user and movie information.
// @Tags review
// @Produce json
// @Success 200 {array} ReviewResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /reviews [get]
func GetReviews(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reviews []ReviewResponse

		result := db.Model(&models.Review{}).
			Select(`reviews.id, users.username as user_name, movies.title as movie_title, reviews.text, 
			CASE WHEN reviews.rating IS NOT NULL THEN reviews.rating ELSE NULL END as rating`).
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

// GetReviewDetails godoc
// @Summary Get review details by ID
// @Description Get details of a specific review by its ID.
// @Tags review
// @Param id path int true "Review ID"
// @Produce json
// @Success 200 {object} ReviewResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /reviews/{id} [get]
func GetReviewDetails(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
			return
		}

		var review ReviewResponse
		result := db.Model(&models.Review{}).
			Select(`reviews.id, users.username as user_name, movies.title as movie_title, reviews.text, 
			CASE WHEN reviews.rating IS NOT NULL THEN reviews.rating ELSE NULL END as rating`).
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

// CreateReview godoc
// @Summary Create a new review
// @Description Create review
// @Tags review
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param movie body CreateReviewRequest true "Review data"
// @Success 201 {object} models.Review
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /reviews [post]
func CreateReview(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := auth.GetUserIDFromToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
			return
		}

		var req CreateReviewRequest

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
		
		var movie models.Movie
		if result := db.First(&movie, req.MovieID); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie not found"})
			return
		}

		if err := req.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		review := models.Review{
			MovieID: req.MovieID,
			UserID:  userID,
			Text:    req.Text,
			Rating:  req.Rating,
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