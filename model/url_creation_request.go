package model

type UrlCreationRequest struct {
	LongUrl string `json:"longUrl" binding:"required"`
}
