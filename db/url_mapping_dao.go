package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"url-shortener/entity"
	"url-shortener/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var Instance *sql.DB

func Connect() {

	LoadEnv()
	dbConnectionURL := os.Getenv("DB_COONECTION_URL")
	database, dbError := sql.Open("mysql", dbConnectionURL)

	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to DB")
	}

	Instance = database
	log.Println("Connected to Database")
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func SaveUrlMapping(shortURL string, longURL string) entity.ShortUrlMapping {

	loc, _ := time.LoadLocation("Asia/Kolkata")
	LoadEnv()
	expiryDuration, err := strconv.Atoi(os.Getenv("SHORT_URL_EXPIRY_TIME_DIFF"))
	expiryInHR := time.Duration(expiryDuration)
	expiryDate := time.Now().In(loc).Add(time.Hour * expiryInHR)
	query := "INSERT INTO short_url_mapping (`short_url`, `long_url`, `expiry_date`) VALUES(?, ?, ?);"
	inserResult, err := Instance.ExecContext(context.Background(), query, shortURL, longURL, expiryDate)

	shortURLMapping := entity.ShortUrlMapping{}

	if err != nil {
		fmt.Println(err)
		return shortURLMapping
	}

	insertedId, err := inserResult.LastInsertId()

	if err != nil {
		fmt.Println(err)
		return shortURLMapping
	}
	shortURLMapping.Id = insertedId
	shortURLMapping.LongUrl = longURL
	shortURLMapping.ShortUrl = shortURL
	shortURLMapping.ExpiryDate = expiryDate

	return shortURLMapping
}

func RetrieveInitialUrl(shortURL string) model.FetchUrlResponse {

	var urlMapping entity.ShortUrlMapping
	err := Instance.QueryRow("SELECT * FROM short_url_mapping where short_url = ?", shortURL).Scan(&urlMapping.Id, &urlMapping.ShortUrl, &urlMapping.LongUrl, &urlMapping.ExpiryDate,
		&urlMapping.CreatedAt)

	urlResponse := model.FetchUrlResponse{}
	if err != nil {
		fmt.Println("Error : ", err)
		if err == sql.ErrNoRows {
			return urlResponse
		}
		return urlResponse
	}

	urlResponse.Id = urlMapping.Id
	urlResponse.ShortUrl = urlMapping.ShortUrl
	urlResponse.LongUrl = urlMapping.LongUrl
	urlResponse.ExpiryDate = urlMapping.ExpiryDate
	urlResponse.CreatedAt = urlMapping.CreatedAt

	fmt.Println("urlResponse : ", urlResponse)

	return urlResponse
}

func GetActiveShortUrlCount() int64 {

	var count int64

	loc, _ := time.LoadLocation("Asia/Kolkata")
	currentDate := time.Now().In(loc)

	err := Instance.QueryRow("SELECT COUNT(*) FROM short_url_mapping where expiry_date > ? ", currentDate).Scan(&count)
	if err != nil {
		fmt.Println("Error : ", err)
	}
	return count

}
