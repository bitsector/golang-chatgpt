package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const apiEndpoint = "https://api.openai.com/v1/chat/completions"

func main() {
	// Load environment variables from .env file (optional)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with existing environment variables")
	}

	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatalf("API key not found in environment variables. Please set OPENAI_API_KEY.")
	}

	// Prepare the request body
	requestBody := map[string]interface{}{
		"model": "gpt-4o-mini", // The model name
		"store": true,          // Optional parameter (if supported by the API)
		"messages": []map[string]string{
			{"role": "user", "content": "write a haiku about ai"},
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

	// Read and process the response
	body, err := ioutil.ReadAll(resp.Body)
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

	fmt.Printf("Response from ChatGPT:\n%s\n", result)
}
