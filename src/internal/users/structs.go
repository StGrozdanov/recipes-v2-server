package users

type BaseUserData struct {
	Username  string `json:"username" db:"username"`
	AvatarURL string `json:"avatarURL" db:"avatar_url"`
}

type OwnerData struct {
	Username string `db:"owner_name" json:"username" `
	Id       int    `db:"ownerid" json:"id"`
}

type User struct {
	Username            string `json:"username" db:"username"`
	AvatarURL           string `json:"avatarURL" db:"avatar_url"`
	CoverPhotoURL       string `json:"coverPhotoURL" db:"cover_photo_url"`
	Email               string `json:"email" db:"email"`
	CreatedRecipesCount int    `json:"createdRecipesCount" db:"created_recipes_count"`
}
