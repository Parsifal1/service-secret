package handler

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"secret/service"
)

func InitRouter(router *mux.Router) {
	// Метод получения секрета
	router.HandleFunc("/secret/{identificator}", GetSecretHandler).Methods("GET")
	// Метод сохранения секрета
	router.HandleFunc("/secret", SaveSecretHandler).Methods("POST")
}

func GetSecretHandler(w http.ResponseWriter, r *http.Request) {
	err := service.GetSecret(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error retrieving secret: %v", err)
		return
	}
}

func SaveSecretHandler(w http.ResponseWriter, r *http.Request) {
	err := service.SaveSecret(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error retrieving secret: %v", err)
		return
	}
}
