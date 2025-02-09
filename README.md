*This is a sample Golang code demonstrating the use of an OpenAI API Key

***Install the requirements

go get github.com/joho/godotenv

go mod tidy

***Set up environment variables in the .env file:

OPENAI_API_KEY=<your_openai_api_key>
OPENAI_API_MODEL=gpt-4o-mini

***Run the code
go run main.go -c "your query here"
