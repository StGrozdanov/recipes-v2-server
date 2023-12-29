package users

import "recipes-v2-server/database"

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
