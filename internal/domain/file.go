package domain


import (
	"time"
	"github.com/google/uuid"
)

type File struct {
	ID               uuid.UUID      `json:"id"`
	Name      		string    `json:"name"`     
	MimeType  		string    `json:"mimeType"` 
	Size      		int64     `json:"size"`
	Path      		string    `json:"path"`
	Base64			string 			`json:"base64"`
	CreatedAt 		time.Time `json:"createdAt"`
}
