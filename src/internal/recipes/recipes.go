package recipes

import "recipes-v2-server/database"

// GetLatest gets the latest 3 recipes
func GetLatest() (recipes []ExtendedRecipeInfo, err error) {
	err = database.GetMultipleRecords(
		&recipes,
		`SELECT recipe_name,
					   image_url,
					   category
				FROM recipes
				ORDER BY created_at DESC
				LIMIT 3;`,
	)
	return
}

// GetMostPopular gets the most visited 3 recipes
func GetMostPopular() (recipes []BaseRecipeInfo, err error) {
	err = database.GetMultipleRecords(
		&recipes,
		`SELECT recipe_name,
					   image_url
				FROM recipes
				ORDER BY visitations_count DESC
				LIMIT 3;`,
	)
	return
}
