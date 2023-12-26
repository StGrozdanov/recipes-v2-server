package recipes

import "recipes-v2-server/database"

// GetAll gets the recipes in a pageable way
func GetAll(limit, cursor int) (recipes RecipePaginationInfo, err error) {
	err = database.GetMultipleRecordsNamedQuery(
		&recipes.BaseRecipeInfoArray,
		`SELECT recipe_name,
					   image_url
				FROM recipes
				WHERE id > :cursor
				ORDER BY created_at
				LIMIT :limit;`,
		map[string]interface{}{"limit": limit, "cursor": cursor},
	)

	if cursor == 0 || cursor-limit < 0 {
		recipes.PageData.PrevPage = 0
	} else {
		recipes.PageData.PrevPage = cursor - limit
	}

	recipes.PageData.NextPage = limit + cursor
	recipes.PageData.FirstPage = cursor == 0
	recipes.PageData.LastPage = len(recipes.BaseRecipeInfoArray) < limit

	return
}

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
