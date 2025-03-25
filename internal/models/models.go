package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique"`
    Password string
    Token    string `gorm:"unique"`
}

type Movie struct {
    gorm.Model
    Title       string     `gorm:"size:200"`
    Year        time.Time
    Runtime     *int       `gorm:"default:null"`
    Rating      *float64   `gorm:"default:null"`
    Votes       *int       `gorm:"default:null"`
    Metascore   *int       `gorm:"default:null"`
    Plot        *string    `gorm:"type:text;default:null"`
    Tagline     *string    `gorm:"size:200;default:null"`
    Genres    	[]Genre    `gorm:"many2many:movie_genres;"`
    Directors 	[]Person   `gorm:"many2many:movie_directors;"`
    Writers   	[]Person   `gorm:"many2many:movie_writers;"`
    Actors    	[]Person   `gorm:"many2many:movie_actors;"`
    CountryID   uint
    Country     Country
    Languages   []Language `gorm:"many2many:movie_languages;"`
    Budget      *int64     `gorm:"default:null"`
    Gross       *int64     `gorm:"default:null"`
    Roles     	[]Role	   `gorm:"foreignKey:MovieID"`
}

type Genre struct {
    gorm.Model
    Name string `gorm:"size:50;unique"`
}

type Person struct {
    gorm.Model
    Name string `gorm:"size:100"`
}

type Role struct {
    gorm.Model
    MovieID   uint    `gorm:"index"`  // Foreign key for Movie
    PersonID  uint    `gorm:"index"`  // Foreign key for Person
    Character string  `gorm:"size:100"`
    
    // Explicit relationship definitions
    Movie     Movie   `gorm:"foreignKey:MovieID;references:ID"`
    Actor     Person  `gorm:"foreignKey:PersonID;references:ID"`
}	

type Country struct {
    gorm.Model
    Name string `gorm:"size:50;unique"`
}

type Language struct {
    gorm.Model
    Name string `gorm:"size:50;unique"`
}

type Review struct {
    gorm.Model
    MovieID    uint       `gorm:"index:,unique,composite:user_movie"`
    UserID     uint       `gorm:"index:,unique,composite:user_movie"`
    Rating     *float64   `gorm:"default:null"`
    Text       *string    `gorm:"type:text;default:null"` 
    Movie      Movie
    User       User
}

type MovieGenre struct {
    MovieID uint `gorm:"primaryKey"`
    GenreID uint `gorm:"primaryKey"`
}

type MovieDirector struct {
    MovieID  uint `gorm:"primaryKey"`
    PersonID uint `gorm:"primaryKey"`
}

type MovieWriter struct {
    MovieID  uint `gorm:"primaryKey"`
    PersonID uint `gorm:"primaryKey"`
}

type MovieActor struct {
    MovieID  uint `gorm:"primaryKey"`
    PersonID uint `gorm:"primaryKey"`
}