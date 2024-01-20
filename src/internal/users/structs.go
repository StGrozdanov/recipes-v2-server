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

type UserAdminData struct {
	User
	Id        int    `db:"user_id" json:"id"`
	Role      string `db:"role" json:"role"`
	IsBlocked bool   `db:"is_blocked" json:"isBlocked"`
}

type UserImages struct {
	AvatarURL     string `json:"avatarURL" db:"avatar_url"`
	CoverPhotoURL string `json:"coverPhotoURL" db:"cover_photo_url"`
}

type UserChangeRoleData struct {
	UserId int    `db:"user_id" json:"userId" valid:"required"`
	Role   string `db:"role" json:"role" valid:"required"`
}

type BlockUserData struct {
	UserId int    `db:"user_id" json:"userId" valid:"required"`
	Reason string `db:"reason" json:"reason" valid:"required"`
}
