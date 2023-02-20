package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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
		Addr:    addr,
		Handler: mux,
	}
	server.ListenAndServe()
}

func generate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	data := string(body)
	log.Println(data)

	w.Write([]byte(data))
}
