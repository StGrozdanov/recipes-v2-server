package comments

import (
	"recipes-v2-server/internal/users"
	"time"
)

type Comment struct {
	Content            string    `json:"content" db:"content"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
	RecipeName         string    `json:"recipeName" db:"recipe_name"`
	Id                 int       `json:"id" db:"id"`
	users.BaseUserData `json:"owner"`
}

type CommentData struct {
	Content         string `json:"content" db:"content" valid:"required"`
	RecipeName      string `json:"recipeName" db:"recipe_name" valid:"required"`
	users.OwnerData `json:"owner"`
}

type CommentIdData struct {
	Id int `json:"id" db:"id" valid:"required"`
}

type CommentEditData struct {
	CommentIdData
	Content string `json:"content" db:"content" valid:"required"`
}
