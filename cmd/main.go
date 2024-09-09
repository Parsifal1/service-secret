package main

import (
	"context"
	"log"
	"net/http"
	"secret/cleaner"
	"secret/handler"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	go cleaner.CheckAndCleanupSecrets(ctx)

	handler.InitRouter(router)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
