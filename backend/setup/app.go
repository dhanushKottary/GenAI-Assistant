package setup

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"sample-project/handlers"
)

type App struct {
	mux *http.ServeMux
}

func NewApp() *App {

	// Load environment variables
	_ = godotenv.Load()

	// Initialize Gemini client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatalf("failed to create Gemini client: %v", err)
	}

	// Register handlers
	h := handlers.NewHandler(client)
	mux := http.NewServeMux()

	mux.HandleFunc("/api/explain-code", h.CodeExplainer)
	mux.HandleFunc("/api/devops-assistant", h.ReviewDevOps)
	mux.HandleFunc("/api/interview", h.HandleInterview)

	return &App{mux: mux}
}

func (a *App) Run() {
	fmt.Println("🚀 Server started on http://localhost:8085")
	log.Fatal(http.ListenAndServe(":8085", a.mux))
}
