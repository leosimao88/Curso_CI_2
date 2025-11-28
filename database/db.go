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

	// retry logic: try for up to ~90 seconds (30 tries * 3s)
	for i := 0; i < 30; i++ {
		DB, err = gorm.Open(postgres.Open(stringDeConexao))
		if err == nil {
			// try to ping the underlying sql DB
			sqlDB, derr := DB.DB()
			if derr == nil {
				if perr := sqlDB.Ping(); perr != nil {
					err = perr
				} else {
					err = nil
					break
				}
			} else {
				err = derr
			}
		}
		// sleep a bit and try again
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		// include the error message to help CI debugging
		log.Panicf("Erro ao conectar com banco de dados: %v", err)
	}

	DB.AutoMigrate(&models.Aluno{})
}
