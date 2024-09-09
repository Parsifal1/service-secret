package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"secret/models"
	"secret/store"
	"unicode/utf8"
)

const (
	keyLetter = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// GenIdentifier Генерирует случайный идентификатор для получения секрета
func GenIdentifier() string {
	bytesSlice := make([]byte, 16)
	for i := range bytesSlice {
		bytesSlice[i] = keyLetter[rand.Int63()%int64(len(keyLetter))]
	}
	return string(bytesSlice)
}

// GetSecret функция получения секрета
func GetSecret(w http.ResponseWriter, r *http.Request) error {

	// Получение секрета
	vars := mux.Vars(r)
	identificator := vars["identificator"]
	secretData, err := store.SelectSecret(identificator)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
	}
	if secretData.ID == 0 {
		http.Error(w, "Секрет с таким идентификатором не найден.", http.StatusNotFound)
		return err
	}

	// Обновление счетчика получения секрета
	secretData.Counter += 1
	if err := store.UpdateSecret(secretData); err != nil {
		http.Error(w, "Произошла ошибка при обновлении данных.", http.StatusInternalServerError)
		return err
	}

	// Проверка счетчика
	if secretData.Counter <= 3 {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(fmt.Sprintf("Ваш секрет: %s \nПолучено раз: %d", secretData.TextSecret, secretData.Counter)))
		if err != nil {
			log.Printf("Ошибка записи ответа: %v", err)
			http.Error(w, "Ошибка при записи ответа.", http.StatusInternalServerError)
			return err
		}
	} else {
		http.Error(w, "Секрет недоступен для получения.", http.StatusForbidden)
	}
	return nil
}

// SaveSecret функция сохранения секрета
func SaveSecret(w http.ResponseWriter, r *http.Request) error {
	// Проверка текста секрета и сохранение
	var secret models.Secret
	err := json.NewDecoder(r.Body).Decode(&secret.TextSecret)
	if err != nil {
		http.Error(w, "Ошибка декодирования данных", http.StatusBadRequest)
		return err
	}
	if utf8.RuneCountInString(secret.TextSecret) != 500 {
		http.Error(w, "Количество символов секрета не равно 500.", http.StatusBadRequest)
		return err
	}

	// Подключение к базе данных
	db, err := store.PostConn()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}(db)

	identifier := GenIdentifier()
	if err := store.InsertSecret(identifier, secret.TextSecret); err != nil {
		http.Error(w, "Произошла ошибка при сохранении секрета.", http.StatusInternalServerError)
		return err
	}

	responseMessage := fmt.Sprintf("Идентификатор: %s\nТекст секрета: %s", identifier, secret.TextSecret)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(responseMessage))
	if err != nil {
		log.Printf("Ошибка записи ответа: %v", err)
		http.Error(w, "Ошибка при записи ответа.", http.StatusInternalServerError)
		return err
	}
	return nil
}
