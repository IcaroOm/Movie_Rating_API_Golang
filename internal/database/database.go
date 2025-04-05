package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"movie-api/internal/models"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("movies.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.SetupJoinTable(&models.Movie{}, "Genres", &models.MovieGenre{})
    db.SetupJoinTable(&models.Movie{}, "Directors", &models.MovieDirector{})
    db.SetupJoinTable(&models.Movie{}, "Writers", &models.MovieWriter{})
    db.SetupJoinTable(&models.Movie{}, "Actors", &models.MovieActor{})
	
	db.AutoMigrate(
		&models.User{},
		&models.Movie{},
		&models.Genre{},
		&models.Person{},
		&models.Role{},
		&models.Country{},
		&models.Language{},
		&models.Review{},
	)
	return db
}

func InitMockDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
        SkipDefaultTransaction: true,
    })
    if err != nil {
        panic("failed to connect database")
    }

    db.SetupJoinTable(&models.Movie{}, "Genres", &models.MovieGenre{})
    db.SetupJoinTable(&models.Movie{}, "Directors", &models.MovieDirector{})
    db.SetupJoinTable(&models.Movie{}, "Writers", &models.MovieWriter{})
    db.SetupJoinTable(&models.Movie{}, "Actors", &models.MovieActor{})
    
    db.AutoMigrate(
        &models.User{},
        &models.Movie{},
        &models.Genre{},
        &models.Person{},
        &models.Role{},
        &models.Country{},
        &models.Language{},
        &models.Review{},
        &models.MovieGenre{},
        &models.MovieDirector{},
        &models.MovieWriter{},
        &models.MovieActor{},
    )

    return db
}

func ClearMockDB(db *gorm.DB) {
    db.Exec("DELETE FROM users")
    db.Exec("DELETE FROM movies")
    db.Exec("DELETE FROM reviews")
    db.Exec("DELETE FROM genres")
    db.Exec("DELETE FROM people")
    db.Exec("DELETE FROM movie_genres")
    db.Exec("DELETE FROM movie_directors")
    db.Exec("DELETE FROM movie_writers")
    db.Exec("DELETE FROM movie_actors")
}