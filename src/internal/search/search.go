package search

import (
	"github.com/lib/pq"
	"recipes-v2-server/database"
)

// UserSearch search users by username and returns them
func UserSearch(query string) (results pq.StringArray, err error) {
	filter := "%" + query + "%"
	err = database.GetSingleRecordNamedQuery(
		&results,
		`SELECT ARRAY(SELECT username FROM users WHERE username LIKE :search);`,
		map[string]interface{}{"search": filter},
	)
	return
}

// RecipesSearch searches recipes by name and returns them
func RecipesSearch(query string) (results pq.StringArray, err error) {
	filter := "%" + query + "%"
	err = database.GetSingleRecordNamedQuery(
		&results,
		`SELECT ARRAY(SELECT recipe_name FROM recipes WHERE recipe_name LIKE :search);`,
		map[string]interface{}{"search": filter},
	)
	return
}

// CommentsSearch searches comments by content and returns them
func CommentsSearch(query string) (results pq.StringArray, err error) {
	filter := "%" + query + "%"
	err = database.GetSingleRecordNamedQuery(
		&results,
		`SELECT ARRAY(SELECT content FROM comments WHERE content LIKE :search);`,
		map[string]interface{}{"search": filter},
	)
	return
}

// Global searches users, recipes, comments and returns them if their name / content matches the filter
func Global(query string) (results []GlobalSearch, err error) {
	filter := "%" + query + "%"
	err = database.GetMultipleRecordsNamedQuery(
		&results,
		`WITH users_search AS (SELECT 'users' AS collection_name, username AS content FROM users WHERE username LIKE :search),
					 recipes_search AS (SELECT 'recipes' AS collection_name, recipe_name AS content
										FROM recipes
										WHERE recipe_name LIKE :search),
					 comments_search AS (SELECT 'comments' AS collection_name, content FROM comments WHERE content LIKE :search)
				
				SELECT collection_name, ARRAY_AGG(content) AS results
				FROM (SELECT *
					  FROM users_search
					  UNION ALL
					  SELECT *
					  FROM recipes_search
					  UNION ALL
					  SELECT *
					  FROM comments_search) AS all_results
				GROUP BY collection_name;`,
		map[string]interface{}{"search": filter},
	)
	return
}
