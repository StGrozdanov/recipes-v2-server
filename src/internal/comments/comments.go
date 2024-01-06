package comments

import (
	"recipes-v2-server/database"
)

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
		`SELECT comments.id,
    				   comments.content,
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

// Edit edits a comment
func Edit(data CommentData) (result Comment, err error) {
	err = database.GetSingleRecordNamedQuery(
		&result,
		`UPDATE comments SET content = :content WHERE id = :id
				RETURNING *`,
		data,
	)
	return
}

// Delete deletes a comment
func Delete(id string) (err error) {
	_, err = database.ExecuteNamedQuery(
		`DELETE
				FROM comments
				WHERE id = :id;`,
		map[string]interface{}{"id": id},
	)
	return
}

// Create creates a new comment
func Create(data CommentData) (result Comment, err error) {
	err = database.GetSingleRecordNamedQuery(
		&result,
		`INSERT INTO comments(content, created_at, owner_id, target_recipe_id)
				VALUES (:content, Now(), :owner_id, (SELECT id FROM recipes WHERE recipe_name = :recipe_name))
				RETURNING *`,
		data,
	)
	return
}
