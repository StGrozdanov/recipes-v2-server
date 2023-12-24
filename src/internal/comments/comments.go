package comments

import "recipes-v2-server/database"

// GetLatestComments gets the latest 6 comments from the database
func GetLatestComments() (comments []Comment, err error) {
	err = database.GetMultipleRecords(
		&comments,
		`SELECT comments.content,
					   comments.created_at,
					   recipes.recipe_name,
					   users.username,
					   users.avatar_url
				FROM comments
						 JOIN recipes ON comments.target_recipe_id = recipes.id
						 JOIN users ON comments.owner_id = users.id
				ORDER BY created_at DESC
				LIMIT 6;`,
	)

	return
}
