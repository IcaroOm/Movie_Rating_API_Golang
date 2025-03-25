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