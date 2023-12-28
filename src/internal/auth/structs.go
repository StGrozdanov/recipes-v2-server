package auth

type UsernameData struct {
	Username string `db:"username" json:"username" valid:"required,minstringlength(3)"`
}

type EmailData struct {
	Email string `db:"email" json:"email" valid:"required,email"`
}

type UserAuthData struct {
	Username string `db:"username" json:"username" valid:"required,minstringlength(3)"`
	Password string `db:"password" json:"password" valid:"required,minstringlength(3)"`
}

type UserRegistrationData struct {
	Username string `db:"username" json:"username" valid:"required,minstringlength(3)"`
	Password string `db:"password" json:"password" valid:"required,minstringlength(3)"`
	Email    string `db:"email" json:"email" valid:"required,email"`
}

type VerificationCodeData struct {
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Code     string `json:"code"`
}

type UserAuthDataResult struct {
	Username        string `db:"username" json:"username"`
	Role            string `db:"role" json:"-"`
	Id              int    `db:"id" json:"id"`
	AvatarURL       string `db:"avatar_url" json:"avatarURL"`
	CoverPhotoURL   string `db:"cover_photo_url" json:"coverPhotoURL"`
	Email           string `db:"email" json:"email"`
	IsAdministrator bool   `db:"is_administrator" json:"isAdministrator"`
	IsModerator     bool   `db:"is_moderator" json:"isModerator"`
	SessionToken    string `json:"sessionToken"`
	NewPassword     string `db:"new_password" json:"-"`
}
