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

// Delete deletes a user and transfers his recipes to a preferred admin user
func Delete(id int) (err error) {
	var oldImageURLs UserImages

	err = database.GetSingleRecordNamedQuery(
		&oldImageURLs,
		`WITH transfer_recipes_to_admin AS (UPDATE recipes SET owner_id = 2 WHERE recipes.owner_id = :id),
					 delete_favourites AS (DELETE FROM users_favourites WHERE user_entity_id = :id),
					 delete_comments AS (DELETE FROM comments WHERE owner_id = :id),
					 delete_roles AS (DELETE FROM users_roles WHERE user_entity_id = :id),
					 delete_ip_address AS (DELETE FROM user_entity_ip_addresses WHERE user_entity_id = :id)
				
				DELETE
				FROM users
				WHERE id = :id
				RETURNING COALESCE(avatar_url, '') AS avatar_url, 
						  COALESCE(cover_photo_url, '') AS cover_photo_url;`,
		map[string]interface{}{"id": id},
	)
	if err != nil {
		return
	}

	if oldImageURLs.AvatarURL != "" {
		err = utils.DeleteFromS3(oldImageURLs.AvatarURL)
	}

	if oldImageURLs.CoverPhotoURL != "" {
		err = utils.DeleteFromS3(oldImageURLs.CoverPhotoURL)
	}

	return
}

// ChangeRole changes a user role
func ChangeRole(data UserChangeRoleData) (err error) {
	_, err = database.ExecuteNamedQuery(
		`UPDATE users_roles
				SET roles_id = (SELECT id FROM roles WHERE role = :role LIMIT 1)
				WHERE user_entity_id = :user_id;`,
		data,
	)
	return
}

// Block blocks a user
func Block(data BlockUserData) (err error) {
	_, err = database.ExecuteNamedQuery(
		`INSERT INTO blacklist (ip_address, reason)
				SELECT ip_addresses, :reason
				FROM user_entity_ip_addresses
				WHERE user_entity_id = :user_id;`,
		data,
	)
	return
}

// Unblock unblocks a user
func Unblock(userId int) (err error) {
	_, err = database.ExecuteNamedQuery(
		`DELETE
				FROM blacklist
				WHERE ip_address IN ((SELECT ip_addresses FROM user_entity_ip_addresses WHERE user_entity_id = :user_id));`,
		map[string]interface{}{"user_id": userId},
	)
	return
}
