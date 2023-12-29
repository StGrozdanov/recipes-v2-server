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
					   COALESCE(users.avatar_url, '') AS avatar_url
				FROM comments
						 JOIN recipes ON comments.target_recipe_id = recipes.id
						 JOIN users ON comments.owner_id = users.id
				ORDER BY created_at DESC
				LIMIT 6;`,
	)

	return
}

// GetCommentsForRecipe gets the comments for the requested recipe name
func GetCommentsForRecipe(recipeName string) (comments []Comment, err error) {
	err = database.GetMultipleRecordsNamedQuery(
		&comments,
		`SELECT comments.content,
					   comments.created_at,
					   recipes.recipe_name,
					   users.username,
					   COALESCE(users.avatar_url, '') AS avatar_url
				FROM comments
						 JOIN recipes ON comments.target_recipe_id = recipes.id
						 JOIN users ON comments.owner_id = users.id
				WHERE recipes.recipe_name = :recipe_name
				ORDER BY created_at DESC;`,
		map[string]interface{}{"recipe_name": recipeName},
	)
	return
}
