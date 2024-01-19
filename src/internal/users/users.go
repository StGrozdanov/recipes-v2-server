package users

import (
	"bytes"
	"mime/multipart"
	"recipes-v2-server/database"
	"recipes-v2-server/utils"
)

// GetUser gets the user details
func GetUser(username string) (user User, err error) {
	err = database.GetSingleRecordNamedQuery(
		&user,
		`SELECT email,
					   username,
					   COALESCE(avatar_url, '')      AS avatar_url,
					   COALESCE(cover_photo_url, '') AS cover_photo_url,
					   COUNT(recipes.id)             AS created_recipes_count
				FROM users
						 LEFT JOIN recipes ON recipes.owner_id = users.id
				WHERE username = :username
				GROUP BY avatar_url, cover_photo_url, email, username;`,
		map[string]interface{}{"username": username},
	)
	return
}

// UploadCoverImage uploads a new cover image for the user
func UploadCoverImage(file *multipart.FileHeader, fileKey, username string) (imageURL string, err error) {
	contentType := file.Header.Get("Content-Type")

	fileContent, _ := file.Open()
	buffer := make([]byte, file.Size)
	_, _ = fileContent.Read(buffer)
	fileBytes := bytes.NewReader(buffer)

	err = utils.UploadToS3(fileBytes, fileKey, contentType)
	if err != nil {
		return
	}

	imageURL = utils.GetTheFullS3BucketURL() + "/" + fileKey

	_, err = database.ExecuteNamedQuery(
		`UPDATE users SET cover_photo_url = :image_url WHERE username = :username;`,
		map[string]interface{}{"username": username, "image_url": imageURL},
	)
	return
}

// UploadAvatarImage uploads a new avatar image for the user
func UploadAvatarImage(file *multipart.FileHeader, fileKey, username string) (imageURL string, err error) {
	contentType := file.Header.Get("Content-Type")

	fileContent, _ := file.Open()
	buffer := make([]byte, file.Size)
	_, _ = fileContent.Read(buffer)
	fileBytes := bytes.NewReader(buffer)

	err = utils.UploadToS3(fileBytes, fileKey, contentType)
	if err != nil {
		return
	}

	imageURL = utils.GetTheFullS3BucketURL() + "/" + fileKey

	_, err = database.ExecuteNamedQuery(
		`UPDATE users SET avatar_url = :image_url WHERE username = :username;`,
		map[string]interface{}{"username": username, "image_url": imageURL},
	)
	return
}

// EditData uploads a new avatar image for the user
func EditData(oldUsername string, data User) (response User, err error) {
	err = database.GetSingleRecordNamedQuery(
		&response,
		`UPDATE users SET username = :username, email = :email 
             	WHERE username = :old_username 
             	RETURNING username, email, avatar_url, cover_photo_url;`,
		map[string]interface{}{"username": data.Username, "email": data.Email, "old_username": oldUsername},
	)
	return
}

// Count retrieves the total count of the users
func Count() (count int, err error) {
	err = database.GetSingleRecord(&count, `SELECT COUNT(id) FROM users;`)
	return
}

// GetAllUsers gets the user details for all users
func GetAllUsers() (user []UserAdminData, err error) {
	err = database.GetMultipleRecords(
		&user,
		`SELECT email,
					   users.id 						AS user_id,
					   role,
					   username,
					   blacklist.ip_address IS NOT NULL AS is_blocked,
					   COALESCE(avatar_url, '')         AS avatar_url,
					   COALESCE(cover_photo_url, '')    AS cover_photo_url,
					   COUNT(recipes.id)                AS created_recipes_count
				FROM users
						 LEFT JOIN recipes ON recipes.owner_id = users.id
						 LEFT JOIN users_roles ON users.id = users_roles.user_entity_id
						 LEFT JOIN roles ON users_roles.roles_id = roles.id
						 LEFT JOIN user_entity_ip_addresses ON users.id = user_entity_ip_addresses.user_entity_id
						 LEFT JOIN blacklist ON user_entity_ip_addresses.ip_addresses = blacklist.ip_address
				GROUP BY avatar_url, cover_photo_url, email, username, users.id, role, blacklist.ip_address;`,
	)
	return
}
