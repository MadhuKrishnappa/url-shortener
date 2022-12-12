package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"url-shortener/db"
	"url-shortener/model"
	"url-shortener/shortener"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func CreateShortUrl(c *gin.Context) {
	var creationRequest model.UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrlHash := shortener.GenerateShortLink(creationRequest.LongUrl)
	LoadEnv()
	host := os.Getenv("HOST")

	if IsShortUrlCreated(shortUrlHash, creationRequest.LongUrl, c, host) {
		return
	}

	if IsShortURLCountLimitReached(c, creationRequest.LongUrl) {
		return
	}

	shortUrlMapping := db.SaveUrlMapping(shortUrlHash, creationRequest.LongUrl)

	c.JSON(200, gin.H{
		"message":      "short url created successfully",
		"shortUrl":     host + shortUrlHash,
		"shortUrlHash": shortUrlHash,
		"longUrl":      creationRequest.LongUrl,
		"expiryDate":   shortUrlMapping.ExpiryDate,
	})

}

func IsShortURLCountLimitReached(c *gin.Context, longUrl string) bool {
	LoadEnv()
	maxShortUrlCount := os.Getenv("MAX_SHORT_URL_COUNT")

	fmt.Println("maxShortUrlCount : ", maxShortUrlCount)

	maxCount, err := strconv.ParseInt(maxShortUrlCount, 6, 12)
	if err != nil {
		fmt.Println(err)
	}
	shortUrlCount := db.GetActiveShortUrlCount()
	fmt.Println("shortUrlCount : ", shortUrlCount)
	if shortUrlCount >= maxCount {
		c.JSON(500, gin.H{
			"message": "Short url creation limit reached, Max : " + maxShortUrlCount,
			"longUrl": longUrl,
		})
		return true
	}

	return false
}

func IsShortUrlCreated(shortUrlHash string, longUrl string, c *gin.Context, host string) bool {
	responseUrlMapping := db.RetrieveInitialUrl(shortUrlHash)

	time := GetCurrentTime()

	if (responseUrlMapping != model.FetchUrlResponse{} && !responseUrlMapping.ExpiryDate.IsZero()) {
		if time.After(responseUrlMapping.ExpiryDate) {
			fmt.Println()
			fmt.Println("ShortUrl is expired")
			fmt.Println()
			c.JSON(500, gin.H{
				"message":      "short url is expired",
				"shortUrl":     host + shortUrlHash,
				"shortUrlHash": shortUrlHash,
				"longUrl":      longUrl,
				"createdAt":    responseUrlMapping.CreatedAt,
				"expiryDate":   responseUrlMapping.ExpiryDate,
			})
			return true
		} else {
			fmt.Println()
			fmt.Println("ShortUrl is already created")
			fmt.Println()
			c.JSON(500, gin.H{
				"message":      "short url is already created",
				"shortUrl":     host + shortUrlHash,
				"shortUrlHash": shortUrlHash,
				"longUrl":      longUrl,
				"createdDate":  responseUrlMapping.CreatedAt,
				"expiryDate":   responseUrlMapping.ExpiryDate,
			})
			return true
		}
	}
	return false
}

func GetCurrentTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	return time.Now().In(loc)
}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrlHash := c.Param("shortUrl")
	fmt.Println(shortUrlHash)
	responseUrlMapping := db.RetrieveInitialUrl(shortUrlHash)

	time := GetCurrentTime()
	LoadEnv()
	host := os.Getenv("HOST")

	if (responseUrlMapping != model.FetchUrlResponse{} && !responseUrlMapping.ExpiryDate.IsZero()) {
		if time.After(responseUrlMapping.ExpiryDate) {
			fmt.Println()
			fmt.Println("ShortUrl is expired")
			fmt.Println()
			c.JSON(500, gin.H{
				"message":      "short url is expired",
				"shortUrl":     host + responseUrlMapping.ShortUrl,
				"shortUrlHash": shortUrlHash,
				"longUrl":      responseUrlMapping.LongUrl,
				"createdAt":    responseUrlMapping.CreatedAt,
				"expiryDate":   responseUrlMapping.ExpiryDate,
			})
			return
		}
	}

	if (responseUrlMapping == model.FetchUrlResponse{}) {
		c.JSON(500, gin.H{
			"message": "Invalid Short URL",
		})
	} else {
		c.JSON(200, gin.H{
			"message":     "Fetched url successfully",
			"longUrl":     responseUrlMapping.LongUrl,
			"createdDate": responseUrlMapping.CreatedAt,
			"expiryDate":  responseUrlMapping.ExpiryDate,
		})
	}
}
