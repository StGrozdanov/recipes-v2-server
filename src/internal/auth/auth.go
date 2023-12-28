package auth

import (
	"github.com/nleeper/goment"
	"golang.org/x/crypto/bcrypt"
	"recipes-v2-server/database"
	"recipes-v2-server/utils"
	"strconv"
	"time"
)

var saltRounds int

// Login accepts username and password validates it and if such user exists - returns JWT auth token
func Login(loginData UserAuthData) (userData UserAuthDataResult, err error) {
	err = database.GetSingleRecordNamedQuery(
		&userData,
		`SELECT users.id,
					   username,
					   COALESCE(avatar_url, '') AS avatar_url,
       				   COALESCE(cover_photo_url, '') AS cover_photo_url,
					   email,
					   password,
					   COALESCE(new_password, '') AS new_password,
					   role,
					   CASE role
						   WHEN 'ADMINISTRATOR' THEN true
						   ELSE false
						   END AS is_administrator,
					   CASE role
						   WHEN 'MODERATOR' THEN true
						   ELSE false
						   END AS is_moderator
				FROM users
						 JOIN users_roles ON users.id = users_roles.user_entity_id
						 JOIN roles ON users_roles.roles_id = roles.id
				WHERE username = :username`,
		loginData,
	)
	if err != nil {
		return
	}

	if userData.NewPassword == "" {
		_, err = updateUserPassword(loginData, &userData)
		if err != nil {
			return
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userData.NewPassword), []byte(loginData.Password)); err != nil {
		return
	}

	jwtToken, err := utils.GenerateJWT(utils.GenerateJWTParams{
		Role:       userData.Role,
		Expiration: 1 * time.Hour,
	})
	if err != nil {
		return
	}

	userData.SessionToken = jwtToken
	return
}

// GetSaltRounds retrieves the bcrypt salt rounds from the config and stores them in memory
func GetSaltRounds(salt string) {
	asNumber, err := strconv.Atoi(salt)
	if err != nil {
		return
	}
	saltRounds = asNumber
}

func updateUserPassword(loginData UserAuthData, userData *UserAuthDataResult) (user UserAuthDataResult, err error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(loginData.Password), saltRounds)

	_, err = database.ExecuteNamedQuery(
		`UPDATE users SET new_password = :password WHERE id = :user_id`,
		map[string]interface{}{"password": hashedPassword, "user_id": userData.Id},
	)
	if err != nil {
		return
	}

	userData.NewPassword = string(hashedPassword)
	return
}

// UsernameExists checks if username exists in the database and returns a boolean value
func UsernameExists(usernameData UsernameData) (exists bool, err error) {
	err = database.GetSingleRecordNamedQuery(
		&exists,
		`SELECT EXISTS(SELECT username FROM users WHERE username = :username);`,
		usernameData,
	)
	return !exists, err
}

// EmailExists checks if email exists in the database and returns a boolean value
func EmailExists(emailData EmailData) (exists bool, err error) {
	err = database.GetSingleRecordNamedQuery(
		&exists,
		`SELECT EXISTS(SELECT email FROM users WHERE email = :email);`,
		emailData,
	)
	return !exists, err
}

// Register user registration handler
func Register(registrationData UserRegistrationData) (userData UserAuthDataResult, err error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(registrationData.Password), saltRounds)
	registrationData.Password = string(hashedPassword)

	err = database.GetSingleRecordNamedQuery(
		&userData,
		`WITH user_data AS (
					INSERT INTO users (is_blocked, email, password, username, new_password)
						VALUES (false, :email, :password, :username, :password)
						RETURNING *),
				
					 role_data AS (
						 INSERT
							 INTO users_roles (user_entity_id, roles_id)
								 VALUES ((SELECT id FROM user_data), 3)
								 RETURNING *)
				
				SELECT user_data.id,
					   username,
					   COALESCE(avatar_url, '')      AS avatar_url,
					   COALESCE(cover_photo_url, '') AS cover_photo_url,
					   email,
					   role,
					   CASE role
						   WHEN 'ADMINISTRATOR' THEN true
						   ELSE false
						   END                       AS is_administrator,
					   CASE role
						   WHEN 'MODERATOR' THEN true
						   ELSE false
						   END                       AS is_moderator
				FROM user_data
						 JOIN role_data ON role_data.user_entity_id = user_data.id
						 JOIN roles ON role_data.roles_id = roles.id;`,
		registrationData,
	)
	if err != nil {
		return
	}

	jwtToken, err := utils.GenerateJWT(utils.GenerateJWTParams{
		Role:       userData.Role,
		Expiration: 1 * time.Hour,
	})
	if err != nil {
		return
	}

	userData.SessionToken = jwtToken
	return
}

// RequestVerificationCode generates a JWT token with reduced expiry used as verification code for the user
// with the given email
func RequestVerificationCode(emailData EmailData) (response VerificationCodeData, err error) {
	jwtToken, jwtErr := utils.GenerateJWT(utils.GenerateJWTParams{
		Role:       "",
		Expiration: 20 * time.Minute,
	})
	if jwtErr != nil {
		return response, jwtErr
	}

	dateNow, dateErr := goment.New(time.Now())
	if err != nil {
		return response, dateErr
	}

	err = database.GetSingleRecordNamedQuery(
		&response,
		`WITH user_data AS (SELECT id, email, username, :code AS code
                   FROM users
                   WHERE email = :email),

					 insert_data AS (INSERT
						 INTO password_requests (code, issued_at, issued_by_user, publication_status_enum)
							 VALUES (:code, :issued_at, (SELECT id FROM user_data), 'PENDING'))
				
				SELECT email,
					   username,
					   code
				FROM user_data;`,
		map[string]interface{}{"email": emailData.Email, "code": jwtToken, "issued_at": dateNow.UTC().Format("YYYY-MM-DD")},
	)
	if err != nil {
		return
	}

	return
}
