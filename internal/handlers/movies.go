package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"movie-api/internal/models"
)

type MovieResponse struct {
	ID            uint    `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	AverageRating float64 `json:"average_rating"`
}

type MovieDetailResponse struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ReleaseDate   string    `json:"release_date"`
	AverageRating float64   `json:"average_rating"`
}

// GetMovies godoc
// @Summary Get list of movies
// @Description Get all movies with average ratings
// @Tags movies
// @Produce json
// @Success 200 {array} MovieResponse
// @Failure 500 {object} map[string]string
// @Router /movies [get]
func GetMovies(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var movies []MovieResponse

		result := db.Model(&models.Movie{}).
			Select("movies.id, movies.title, movies.description, "+
				"COALESCE(AVG(reviews.rating), 0) as average_rating").
			Joins("LEFT JOIN reviews ON reviews.movie_id = movies.id").
			Group("movies.id").
			Scan(&movies)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
			return
		}

		c.JSON(http.StatusOK, movies)
	}
}

// GetMovieDetails godoc
// @Summary Get movie details
// @Description Get detailed information about a specific movie
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} MovieDetailResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /movies/{id} [get]
func GetMovieDetails(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
			return
		}

		var movie MovieDetailResponse
		result := db.Model(&models.Movie{}).
			Select("movies.id, movies.title, movies.description, "+
				"movies.release_date, COALESCE(AVG(reviews.rating), 0) as average_rating").
			Joins("LEFT JOIN reviews ON reviews.movie_id = movies.id").
			Where("movies.id = ?", id).
			Group("movies.id").
			First(&movie)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}