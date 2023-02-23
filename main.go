package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sausheong/openai"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load the env vars: %v", err)
	}
}

func main() {

	server()
}

func server() {
	addr := os.Getenv("PORT")
	mux := http.NewServeMux()
	mux.HandleFunc("/ask", a)
	mux.HandleFunc("/gen", g)
	server := &http.Server{
		Addr:    ":" + addr,
		Handler: mux,
	}
	server.ListenAndServe()
}

func a(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read the body: %v", err)
	}
	data := string(body)
	text, err := ask("", data)
	if err != nil {
		log.Printf("Failed to talk to OpenAI: %v", err)
	}
	w.Write([]byte(text))
}

func g(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read the body: %v", err)
	}
	data := string(body)
	img, err := generate(data)
	if err != nil {
		log.Printf("Failed to talk to OpenAI: %v", err)
	}
	w.Write([]byte(img))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// ask OpenAI to generate text
func ask(header string, prompt string) (string, error) {
	openAIApiKey := os.Getenv("OPENAI_API_KEY")
	openAIOrganization := os.Getenv("OPENAI_ORGANIZATION")

	oaClient := openai.NewClient(openAIApiKey, openAIOrganization)
	request := make(openai.CompletionRequest)
	request.SetUser("sausheong")
	request.SetModel(openai.TEXT_DAVINCI_003)
	request.SetPrompt(fmt.Sprintf("%s:%s", header, prompt))
	request["temperature"] = 0.75
	request["max_tokens"] = 4000

	cr, err := oaClient.Complete(request)
	if err != nil {
		log.Println("Completion request failed:", err)
	}
	return cr.Text(), err
}

func generate(prompt string) (string, error) {
	openAIApiKey := os.Getenv("OPENAI_API_KEY")
	openAIOrganization := os.Getenv("OPENAI_ORGANIZATION")

	oaClient := openai.NewClient(openAIApiKey, openAIOrganization)
	request := make(openai.ImageRequest)
	request.SetUser("sausheong")
	request.SetPrompt(prompt)
	request.SetFormat("b64_json")
	request.SetSize("512x512")

	cr, err := oaClient.GenerateImage(request)
	if err != nil {
		log.Println("Completion request failed:", err)
	}
	return cr.ImageBase64(), err
}
