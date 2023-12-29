package users

type BaseUserData struct {
	Username  string `json:"username" db:"username"`
	AvatarURL string `json:"avatarURL" db:"avatar_url"`
}

type OwnerData struct {
	Username string `db:"owner_name" json:"username" `
	Id       int    `db:"ownerid" json:"id"`
}
