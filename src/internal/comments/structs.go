package comments

import (
	"recipes-v2-server/internal/users"
	"time"
)

type Comment struct {
	Content            string    `json:"content" db:"content"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
	RecipeName         string    `json:"recipeName" db:"recipe_name"`
	users.BaseUserData `json:"owner"`
}
