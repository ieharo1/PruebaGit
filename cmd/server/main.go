package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"callflowmanager/internal/config"
	"callflowmanager/internal/handlers"
	"callflowmanager/internal/middlewares"
	"callflowmanager/internal/repositories"
	"callflowmanager/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	wd, _ := os.Getwd()

	if err := godotenv.Load(filepath.Join(wd, ".env")); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	if _, err := config.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repositories.NewAgentRepository().CreateIndexes(ctx)
	repositories.NewCustomerRepository().CreateIndexes(ctx)
	repositories.NewCallRepository().CreateIndexes(ctx)
	repositories.NewUserRepository().CreateIndexes(ctx)

	authService := services.NewAuthService()
	agentService := services.NewAgentService()
	customerService := services.NewCustomerService()
	callService := services.NewCallService()

	authHandler := handlers.NewAuthHandler(authService, cfg)
	agentHandler := handlers.NewAgentHandler(agentService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	callHandler := handlers.NewCallHandler(callService)

	authMiddleware := middlewares.NewAuthMiddleware(cfg)
	rateLimiter := middlewares.NewRateLimiter(100, 1)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.CORSMiddleware())

	r.LoadHTMLGlob(filepath.Join(wd, "web/templates/*"))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "CallFlowManager - Sistema de Call Center"})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{"title": "Login - CallFlowManager"})
	})

	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{"title": "Registro - CallFlowManager"})
	})

	api := r.Group("/api")
	api.Use(middlewares.RateLimit(rateLimiter))
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			protected.GET("/agents", agentHandler.GetAgents)
			protected.POST("/agents", agentHandler.CreateAgent)

			protected.GET("/customers", customerHandler.GetCustomers)
			protected.POST("/customers", customerHandler.CreateCustomer)

			protected.GET("/calls", callHandler.GetCalls)
			protected.POST("/calls", callHandler.CreateCall)
			protected.PUT("/calls/:id/status", callHandler.UpdateStatus)
			protected.GET("/calls/stats", callHandler.GetStats)
		}
	}

	web := r.Group("/dashboard")
	web.Use(authMiddleware.RequireAuth())
	{
		web.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, "dashboard.html", gin.H{"title": "Dashboard - CallFlowManager"})
		})
		web.GET("/calls", func(c *gin.Context) {
			c.HTML(http.StatusOK, "calls.html", gin.H{"title": "Llamadas - CallFlowManager"})
		})
		web.GET("/customers", func(c *gin.Context) {
			c.HTML(http.StatusOK, "customers.html", gin.H{"title": "Clientes - CallFlowManager"})
		})
		web.GET("/agents", func(c *gin.Context) {
			c.HTML(http.StatusOK, "agents.html", gin.H{"title": "Agentes - CallFlowManager"})
		})
	}

	r.Static("/static", filepath.Join(wd, "web/static"))

	port := cfg.Port
	log.Printf("Server starting on port %s", port)

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Server exited")
}
