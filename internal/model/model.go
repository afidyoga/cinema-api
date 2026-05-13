package model

import (
	"time"
)

type User struct {
	ID           string    `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Phone        string    `db:"phone" json:"phone"`
	Role         string    `db:"role" json:"role"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type Cinema struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	City      string    `db:"city" json:"city"`
	Address   string    `db:"address" json:"address"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Studio struct {
	ID        string    `db:"id" json:"id"`
	CinemaID  string    `db:"cinema_id" json:"cinema_id"`
	Name      string    `db:"name" json:"name"`
	Capacity  int       `db:"capacity" json:"capacity"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Movie struct {
	ID              string    `db:"id" json:"id"`
	Title           string    `db:"title" json:"title"`
	DurationMinutes int       `db:"duration_minutes" json:"duration_minutes"`
	Genre           string    `db:"genre" json:"genre"`
	Rating          string    `db:"rating" json:"rating"`
	Synopsis        string    `db:"synopsis" json:"synopsis"`
	PosterURL       string    `db:"poster_url" json:"poster_url"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

type Schedule struct {
	ID        string    `db:"id" json:"id"`
	MovieID   string    `db:"movie_id" json:"movie_id"`
	StudioID  string    `db:"studio_id" json:"studio_id"`
	ShowTime  time.Time `db:"show_time" json:"show_time"`
	Price     float64   `db:"price" json:"price"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	MovieTitle  string `db:"movie_title" json:"movie_title,omitempty"`
	StudioName  string `db:"studio_name" json:"studio_name,omitempty"`
	CinemaName  string `db:"cinema_name" json:"cinema_name,omitempty"`
}

type CreateScheduleRequest struct {
	MovieID  string    `json:"movie_id" binding:"required"`
	StudioID string    `json:"studio_id" binding:"required"`
	ShowTime time.Time `json:"show_time" binding:"required"`
	Price    float64   `json:"price" binding:"required,gt=0"`
}

type UpdateScheduleRequest struct {
	ShowTime *time.Time `json:"show_time"`
	Price    *float64   `json:"price"`
	Status   *string    `json:"status"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
