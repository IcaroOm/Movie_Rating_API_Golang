package handlers

import (
	"fmt"
	"movie-api/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// CreateMovie godoc
// @Summary Create a new movie
// @Description Create movie with relationships
// @Tags movies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param movie body CreateMovieRequest true "Movie data"
// @Success 201 {object} models.Movie
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /movies [post]
func CreateMovie(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req CreateMovieRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if req.Year < 1888 || req.Year > time.Now().Year()+5 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
            return
        }

        movie := models.Movie{
            Title:     req.Title,
            Year:      req.Year,
            Runtime:   req.Runtime,
            Description:      req.Description,
            Tagline:   req.Tagline,
            Budget:    req.Budget,
            Gross:     req.Gross,
        }

        tx := db.Begin()
        defer func() {
            if r := recover(); r != nil {
                tx.Rollback()
            }
        }()

        if err := handleRelationships(tx, &movie, req); err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := tx.Create(&movie).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create movie"})
            return
        }

        tx.Commit()
        c.JSON(http.StatusCreated, movie)
    }
}

func handleRelationships(tx *gorm.DB, movie *models.Movie, req CreateMovieRequest) error {
    if len(req.GenreIDs) > 0 {
        var genres []models.Genre
        if err := tx.Find(&genres, "id IN ?", req.GenreIDs).Error; err != nil {
            return err
        }
        movie.Genres = genres
    }

    if len(req.DirectorIDs) > 0 {
        var directors []models.Person
        if err := tx.Find(&directors, "id IN ?", req.DirectorIDs).Error; err != nil {
            return err
        }
        movie.Directors = directors
    }

	if len(req.WriterIDs) > 0 {
		var writers []models.Person
		if err := tx.Find(&writers, "in IN ?", req.WriterIDs).Error; err != nil {
            return err 
        }
        movie.Writers = writers
	}

    if len(req.WriterIDs) > 0 {
        var actors []models.Person
        if err := tx.Find(&actors, "in IN ?", req.ActorIDs).Error; err != nil {
            return err
        }
        movie.Actors = actors
    }

    if len(req.LanguageIDs) > 0 {
        var languages []models.Language
        if err := tx.Find(&languages, "in IN ?", req.ActorIDs).Error; err != nil {
            return err
        }
        movie.Languages = languages
    }

	if req.CountryID != nil { 
		var country models.Country
		if err := tx.First(&country, *req.CountryID).Error; err != nil {
			return fmt.Errorf("invalid country ID")
		}
		movie.Country = country
	}

    return nil
}

type CreateMovieRequest struct {
    Title     string  `json:"title" binding:"required"`
    Year      int     `json:"year" binding:"required"`
    Runtime   *int    `json:"runtime,omitempty"`
    Description      *string `json:"plot,omitempty"`
    Tagline   *string `json:"tagline,omitempty"`
    GenreIDs  []uint  `json:"genre_ids,omitempty"`
    DirectorIDs []uint `json:"director_ids,omitempty"`
    WriterIDs []uint  `json:"writer_ids,omitempty"`
    ActorIDs  []uint  `json:"actor_ids,omitempty"`
    CountryID *uint    `json:"country_id,omitempty"`
    LanguageIDs []uint `json:"language_ids,omitempty"`
    Budget    *int64  `json:"budget,omitempty"`
    Gross     *int64  `json:"gross,omitempty"`
}