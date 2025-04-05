package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"movie-api/internal/auth"
	"movie-api/internal/database"
	"movie-api/internal/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
    // Initialize once
    testDB = database.InitMockDB()
    
    // Run tests
    code := m.Run()
    
    // Optional: Add cleanup if needed
    database.ClearMockDB(testDB)
    
    os.Exit(code)
}


func setupRouter() *gin.Engine {
	db := testDB
    r := gin.Default()
    r.SetTrustedProxies(nil)

    r.POST("/api/token", auth.LoginHandler(db))
    r.POST("/api/users", auth.CreateUser(db))
    r.GET("/api/movies", handlers.GetMovies(db))
    r.GET("/api/movies/:id/", handlers.GetMovieDetails(db))
    r.GET("/api/reviews", handlers.GetReviews(db))
    r.GET("/api/reviews/:id/", handlers.GetReviewDetails(db))
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authGroup := r.Group("/")
    authGroup.Use(auth.JWTAuthMiddleware(db))
    {
        authGroup.POST("/api/reviews", handlers.CreateReview(db))
        authGroup.POST("/api/movies", handlers.CreateMovie(db))
    }

    return r
}

func TestPublicRoutes(t *testing.T) {
    router := setupRouter()

    t.Run("GET /api/movies", func(t *testing.T) {
        req, _ := http.NewRequest("GET", "/api/movies", nil)
        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        if resp.Code != http.StatusOK {
            t.Errorf("Expected status %d but got %d", http.StatusOK, resp.Code)
        }
    })
}

func TestAuthenticatedRoutes(t *testing.T) {
    router := setupRouter()

    t.Run("POST /api/reviews (unauthenticated)", func(t *testing.T) {
        req, _ := http.NewRequest("POST", "/api/reviews", nil)
        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        if resp.Code != http.StatusUnauthorized {
            t.Errorf("Expected status %d but got %d", http.StatusUnauthorized, resp.Code)
        }
    })

    t.Run("POST /api/movies (unauthenticated)", func(t *testing.T) {
        req, _ := http.NewRequest("POST", "/api/movies", nil)
        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        if resp.Code != http.StatusUnauthorized {
            t.Errorf("Expected status %d but got %d", http.StatusUnauthorized, resp.Code)
        }
    })
}

func TestAuthenticatedEndpoints(t *testing.T) {
    router := setupRouter()

    t.Run("POST /api/users (create user)", func(t *testing.T) {
        reqBody := `{"password": "testuser1", "username": "testpassword"}`
        req, _ := http.NewRequest("POST", "/api/users", strings.NewReader(reqBody))
        req.Header.Set("Content-Type", "application/json")
        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        if resp.Code != http.StatusCreated {
            t.Fatalf("Failed to create user. Expected status %d but got %d", http.StatusOK, resp.Code)
        }
    })

    var token string
    t.Run("POST /api/token (generate token)", func(t *testing.T) {
        reqBody := `{"password": "testuser1", "username": "testpassword"}`
        req, _ := http.NewRequest("POST", "/api/token", strings.NewReader(reqBody))
        req.Header.Set("Content-Type", "application/json")
        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        if resp.Code != http.StatusOK {
            t.Fatalf("Failed to generate token. Expected status %d but got %d", http.StatusOK, resp.Code)
        }

        // Extract token from response body
        var tokenResponse struct {
			Token string `json:"token"`
		}
		json.Unmarshal(resp.Body.Bytes(), &tokenResponse)
		token = tokenResponse.Token
    })

    t.Run("POST /api/movies (authenticated)", func(t *testing.T) {
        reqBody := `{"title": "Inception", "director": "Christopher Nolan", "year": 2010}`
        req, _ := http.NewRequest("POST", "/api/movies", strings.NewReader(reqBody))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", token)
        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        if resp.Code != http.StatusCreated {
            t.Errorf("Expected status %d but got %d", http.StatusOK, resp.Code)
        }
    })

    t.Run("POST /api/review (authenticated)", func(t *testing.T) {
        reqBody := `{"movie_id": 1, "text": "Amazing movie!", "user_rating": 5}`
        req, _ := http.NewRequest("POST", "/api/reviews", strings.NewReader(reqBody))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", token)
        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        if resp.Code != http.StatusCreated {
            t.Errorf("Expected status %d but got %d", http.StatusOK, resp.Code)
        }
    })
}

func TestLogin(t *testing.T) {
    router := setupRouter()

    t.Run("POST /api/users", func(t *testing.T) {
        reqBody := `{"username": "newuser", "password": "testpassword"}`
        req, _ := http.NewRequest("POST", "/api/users", strings.NewReader(reqBody))
        req.Header.Set("Content-Type", "application/json")
        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        if resp.Code != http.StatusCreated {
            t.Errorf("Expected status %d but got %d", http.StatusOK, resp.Code)
        }
    })

    t.Run("POST /api/token", func(t *testing.T) {
        reqBody := `{"username": "newuser", "password": "testpassword"}`
        req, _ := http.NewRequest("POST", "/api/token", strings.NewReader(reqBody))
        req.Header.Set("Content-Type", "application/json")
        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        if resp.Code != http.StatusOK {
            t.Errorf("Expected status %d but got %d", http.StatusOK, resp.Code)
        }
    })
}