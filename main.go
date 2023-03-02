package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sausheong/openai"
)

var env string
var openAIApiKey string
var openAIOrganization string

func init() {
	env = os.Getenv("ENV")
	if env != "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Printf("Failed to load the env vars: %v", err)
		}
	}
	openAIApiKey = os.Getenv("OPENAI_API_KEY")
	openAIOrganization = os.Getenv("OPENAI_ORGANIZATION")
}

func main() {
	addr := os.Getenv("PORT")
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/ask", a)
	mux.HandleFunc("/gen", g)
	server := &http.Server{
		Addr:    ":" + addr,
		Handler: mux,
	}
	server.ListenAndServe()
}

// handler for /ask
func a(w http.ResponseWriter, r *http.Request) {
	if env != "prod" {
		enableCors(&w)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read the body: %v", err)
	}
	data := string(body)
	text, err := ask(data)
	if err != nil {
		log.Printf("Failed to talk to OpenAI: %v", err)
	}
	w.Write([]byte(text))
}

// handler for /gen
func g(w http.ResponseWriter, r *http.Request) {
	if env != "prod" {
		enableCors(&w)
	}
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

// enable CORS for the API
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// ask OpenAI to generate text
func ask(prompt string) (string, error) {
	oaClient := openai.NewClient(openAIApiKey, openAIOrganization)
	request := make(openai.CompletionRequest)
	request.SetModel(openai.TEXT_GPT_35_TURBO)
	request.SetPrompt(prompt + " {}")
	request["temperature"] = 0.75
	request["max_tokens"] = 4096 - len(prompt)
	request["stop"] = "{}"

	cr, err := oaClient.Complete(request)
	if err != nil {
		log.Println("Completion request failed:", err)
	}
	return cr.Text(), err
}

// ask OpenAI to generate an image
func generate(prompt string) (string, error) {
	oaClient := openai.NewClient(openAIApiKey, openAIOrganization)
	request := make(openai.ImageRequest)
	request.SetPrompt(prompt)
	request.SetFormat("b64_json")
	request.SetSize("512x512")

	cr, err := oaClient.GenerateImage(request)
	if err != nil {
		log.Println("Completion request failed:", err)
	}
	return cr.ImageBase64(), err
}
