package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// ChatRequest represents the request structure for the chat API
type ChatRequest struct {
	Message string `json:"message"`
}

// DeepSeekRequest represents the request structure for the DeepSeek API
type DeepSeekRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message represents a single message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepSeekResponse represents the response structure from the DeepSeek API
type DeepSeekResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a choice in the DeepSeek API response
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents the token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamResponse represents a streaming response
type StreamResponse struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	// Get server port from environment variables
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // default port
	}

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./static")

	// Define routes
	router.GET("/", indexHandler)
	router.POST("/chat", chatHandler)

	// Start server
	log.Printf("Server starting on port %s\n", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func indexHandler(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/static/index.html")
}

func chatHandler(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get API key and URL from environment variables
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	apiURL := os.Getenv("DEEPSEEK_API_URL")

	if apiKey == "" || apiURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API key or URL not configured"})
		return
	}

	// Prepare request for DeepSeek API
	deepSeekReq := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{Role: "user", Content: req.Message},
		},
	}

	// Convert request to JSON
	requestBody, err := json.Marshal(deepSeekReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare request"})
		return
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", apiURL+"/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to DeepSeek API"})
		return
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from DeepSeek API"})
		return
	}

	// Parse response
	var deepSeekResp DeepSeekResponse
	if err := json.Unmarshal(body, &deepSeekResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response from DeepSeek API"})
		return
	}

	// Extract the response message
	if len(deepSeekResp.Choices) > 0 {
		responseText := deepSeekResp.Choices[0].Message.Content
		c.JSON(http.StatusOK, gin.H{"response": responseText})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No response from DeepSeek API"})
	}
}
