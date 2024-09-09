package store

import (
	"database/sql"
	"log"
	model "secret/models"

	_ "github.com/lib/pq"
)

// Подключение к БД
func PostConn() (*sql.DB, error) {
	connect := "user=postgres password=1234567890 dbname=secret sslmode=disable"
	db, err := sql.Open("postgres", connect)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return db, nil
}

func SelectSecret(identifier string) (model.Secret, error) {
	// Подключение к PostgreSQL
	db, err := PostConn()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return model.Secret{}, err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}(db)

	var result model.Secret
	db.QueryRow(`SELECT id, identifier, text_secret, counter, save_date FROM secret_shema.secret where ("identifier") = $1`, identifier).Scan(&result.ID, &result.Identifier, &result.TextSecret, &result.Counter, &result.SaveDate)
	return result, err
}

func InsertSecret(identifier string, textSecret string) error {
	// Подключение к PostgreSQL
	db, err := PostConn()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}(db)
	res, err := db.Prepare(`INSERT INTO secret_shema.secret ("identifier", "text_secret", "counter","save_date") VALUES ($1, $2, 0 , current_timestamp)`)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = res.Exec(identifier, textSecret)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteSecret() error {
	// Подключение к PostgreSQL
	db, err := PostConn()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}(db)

	res, err := db.Prepare(`DELETE FROM secret_shema.secret WHERE (save_date + interval  '72 hours') < now();`)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = res.Exec()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func UpdateSecret(secretData model.Secret) error {
	// Подключение к PostgreSQL
	db, err := PostConn()
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}(db)
	res, err := db.Prepare(`UPDATE secret_shema.secret SET counter = $1 WHERE ("id") = $2`)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = res.Exec(secretData.Counter, secretData.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
