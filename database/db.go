package database

import (
	"log"
	"os"
	"time"
	"github.com/guilhermeonrails/api-go-gin/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func ConectaComBancoDeDados() {
	// Use explicit DB_* env variables to avoid clobbering system vars (e.g. USER)
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "root"
	}
	if password == "" {
		password = "root"
	}
	if dbname == "" {
		dbname = "root"
	}

	stringDeConexao := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"

	// retry logic: try for up to ~30 seconds
	for i := 0; i < 15; i++ {
		DB, err = gorm.Open(postgres.Open(stringDeConexao))
		if err == nil {
			break
		}
		// sleep a bit and try again
		// Use time.Sleep without import collision
		// import time
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Panic("Erro ao conectar com banco de dados")
	}

	DB.AutoMigrate(&models.Aluno{})
}
