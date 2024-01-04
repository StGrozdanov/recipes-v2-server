package users

type BaseUserData struct {
	Username  string `json:"username" db:"username"`
	AvatarURL string `json:"avatarURL" db:"avatar_url"`
}

type OwnerData struct {
	Username string `db:"owner_name" json:"username" valid:"required"`
	Id       int    `db:"owner_id" json:"id" valid:"required"`
}

type User struct {
	Username            string `json:"username" db:"username" valid:"required,minstringlength(3)"`
	AvatarURL           string `json:"avatarURL" db:"avatar_url"`
	CoverPhotoURL       string `json:"coverPhotoURL" db:"cover_photo_url"`
	Email               string `json:"email" db:"email" valid:"required,email"`
	CreatedRecipesCount int    `json:"createdRecipesCount" db:"created_recipes_count"`
}
