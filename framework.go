package framework

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type App struct {
	Name    string
	Version string
	Port    int
	router  *Router
}

func NewApp() *App {
	return &App{
		Name:    "Gor Framewirk",
		Version: "0.1.0",
		Port:    81,
		router:  NewRouter(),
	}
}

// LoadEnv loads environment variables from file
func (g *App) LoadEnv(filenames ...string) {
	// By default, look for .env
	if len(filenames) == 0 {
		filenames = []string{".env"}
	}

	// Load .env
	fmt.Println("Load env files", filenames)
	if err := godotenv.Load(filenames...); err != nil {
		// Ignore error if file doesn't exist (this is normal)
		fmt.Printf("Warning: .env file not found or error loading: %v\n", err)
	}

	// After loading .env, update port from environment variable
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			g.Port = p
			fmt.Printf("Loaded PORT=%d from .env file\n", g.Port)
		} else {
			fmt.Printf("Warning: invalid PORT value in environment: %s\n", envPort)
		}
	}
}

// Router returns the router instance
func (g *App) Router() *Router {
	return g.router
}

// Run initializes and starts the framework
func (g *App) Run() error {
	addr := fmt.Sprintf(":%d", g.Port)
	fmt.Printf("Starting web server on port %d...\n", g.Port)
	return http.ListenAndServe(addr, g.router)
}
