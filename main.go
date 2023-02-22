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
	mux.HandleFunc("/g", generate)
	server := &http.Server{
		Addr:    ":" + addr,
		Handler: mux,
	}
	server.ListenAndServe()
}

func generate(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	data := string(body)
	log.Println(data)

	text, err := talkToOpenAI("tell me more about ", data)
	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte(text))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// talk to OpenAI API
func talkToOpenAI(header string, prompt string) (string, error) {
	openAIApiKey := os.Getenv("OPENAI_API_KEY")
	openAIOrganization := os.Getenv("OPENAI_ORGANIZATION")

	oaClient := openai.NewClient(openAIApiKey, openAIOrganization)
	request := make(openai.CompletionRequest)
	request.SetUser("sausheong")
	request.SetModel(openai.TEXT_DAVINCI_002)
	request.SetPrompt(fmt.Sprintf("%s:%s", header, prompt))
	request["temperature"] = 0.75
	request["max_tokens"] = 50

	cr, err := oaClient.Complete(request)
	if err != nil {
		log.Println("Completion request failed:", err)
	}
	return cr.Text(), err
}
