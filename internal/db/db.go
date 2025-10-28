package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func Init(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Ошибка инициализации соединения: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("БД недоступна: %v", err)
	}

	if err := db.AutoMigrate(
		&User{},
		&Order{},
		&Withdrawal{},
	); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	log.Println("Подключение к БД успешно")
	return db
}

func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}
