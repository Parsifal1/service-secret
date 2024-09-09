package cleaner

import (
	"context"
	"database/sql"
	"log"
	"secret/store"
	"sync"
	"time"
)

func CheckAndCleanupSecrets(ctx context.Context) {
	//Подключение к БД
	db, err := store.PostConn()
	if err != nil {
		log.Println(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)

	//Проверка истечения 72 часов
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := store.DeleteSecret(); err != nil {
					log.Println(err)
				}
			}
		}
	}()
	wg.Wait()
}
