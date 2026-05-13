package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/afidyoga/cinema-api/internal/config"
	"github.com/afidyoga/cinema-api/internal/handler"
	"github.com/afidyoga/cinema-api/internal/middleware"
	"github.com/afidyoga/cinema-api/internal/repository"
	"github.com/afidyoga/cinema-api/internal/service"
)

func main() {
	cfg := config.Load()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("database connected")

	userRepo := repository.NewUserRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)

	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	scheduleSvc := service.NewScheduleService(scheduleRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	scheduleHandler := handler.NewScheduleHandler(scheduleSvc)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", middleware.Auth(authSvc), authHandler.Me)
	}

	schedules := v1.Group("/schedules")
	schedules.Use(middleware.Auth(authSvc))
	{
		schedules.GET("", scheduleHandler.GetAll)
		schedules.GET("/:id", scheduleHandler.GetByID)
		schedules.POST("", middleware.AdminOnly(), scheduleHandler.Create)
		schedules.PUT("/:id", middleware.AdminOnly(), scheduleHandler.Update)
		schedules.DELETE("/:id", middleware.AdminOnly(), scheduleHandler.Delete)
	}

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
