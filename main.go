package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const apiEndpoint = "https://api.openai.com/v1/chat/completions"

func main() {
	// Define a flag for user input
	var userContent string
	flag.StringVar(&userContent, "c", "", "Content for the request")
	flag.Parse()

	// Check if user content is provided
	if userContent == "" {
		log.Fatalf("No content provided. Use -c flag to specify the content, e.g., go run main.go -c \"your request here\"")
	}

	// Get API key and model configuration
	apiKey, model := getConfig()

	// Prepare the request body using user-provided content
	requestBody := map[string]interface{}{
		"model": model, // The model name from env var or default
		"store": true,  // Optional parameter (if supported by the API)
		"messages": []map[string]string{
			{"role": "user", "content": userContent},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Error marshaling request body: %v", err)
	}

	// Create an HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making API request: %v", err)
	}
	defer resp.Body.Close()

	// Read and process the response using io.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API returned an error: %s\nResponse: %s", resp.Status, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Error unmarshaling response JSON: %v", err)
	}

	// Extract and print only the content field
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		firstChoice := choices[0].(map[string]interface{})
		if message, ok := firstChoice["message"].(map[string]interface{}); ok {
			if content, ok := message["content"].(string); ok {
				fmt.Printf("Response from ChatGPT:\n%s\n", content)
				return
			}
		}
	}

	log.Fatalf("Unexpected response structure: %s", string(body))
}

// getConfig retrieves API key and model name from environment variables or defaults.
func getConfig() (string, string) {
	apiKeySource := "environment variable"

	// Load environment variables from .env file (optional)
	err := godotenv.Load()
	if err == nil {
		apiKeySource = ".env file" // If .env is loaded successfully, assume key might come from there
	} else {
		log.Println("No .env file found, proceeding with existing environment variables")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatalf("API key not found in environment variables. Please set OPENAI_API_KEY.")
	}

	maskedKey := fmt.Sprintf("***%s", apiKey[len(apiKey)-4:])
	fmt.Printf("Got API key %s from %s\n", maskedKey, apiKeySource)

	model := os.Getenv("OPENAI_API_MODEL")
	if model == "" {
		model = "gpt-4o-mini" // Default model if none is specified
		fmt.Println("Using default model: gpt-4o-mini (no model found in .env file or environment variable)")
	} else {
		fmt.Printf("Using model: %s\n", model)
	}

	return apiKey, model
}
