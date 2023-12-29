package recipes

import (
	"recipes-v2-server/database"
)

// GetAll gets the recipes in a pageable way
func GetAll(limit, cursor int) (recipes RecipePaginationInfo, err error) {
	err = database.GetMultipleRecordsNamedQuery(
		&recipes.BaseRecipeInfoArray,
		`SELECT recipe_name,
					   image_url
				FROM recipes
				WHERE id > :cursor AND status = 'APPROVED'
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
				WHERE status = 'APPROVED'
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

// Search searches for recipes by name with the provided string
func Search(query string) (recipes []BaseRecipeInfo, err error) {
	filter := "%" + query + "%"

	err = database.GetMultipleRecordsNamedQuery(
		&recipes,
		`SELECT recipe_name,
					   image_url
				FROM recipes
				WHERE recipe_name LIKE :query AND status = 'APPROVED'
				ORDER BY visitations_count DESC;`,
		map[string]interface{}{"query": filter},
	)

	return
}

// SearchByCategory searches for recipes by category name with the provided string
func SearchByCategory(query string) (recipes []BaseRecipeInfo, err error) {
	err = database.GetMultipleRecordsNamedQuery(
		&recipes,
		`SELECT recipe_name,
					   image_url
				FROM recipes
				WHERE category = :query AND status = 'APPROVED'
				ORDER BY visitations_count DESC;`,
		map[string]interface{}{"query": query},
	)
	return
}

// GetASingleRecipe gets the recipe with provided name from the database
func GetASingleRecipe(recipeName string) (recipe RecipeData, err error) {
	err = database.GetSingleRecordNamedQuery(
		&recipe,
		`WITH steps_results AS (SELECT ARRAY_AGG(steps) AS steps
                       FROM recipe_entity_steps
                       WHERE recipe_entity_id = (SELECT id FROM recipes WHERE recipe_name = :recipe_name)),

					 products_results AS (SELECT ARRAY_AGG(products) AS products
										  FROM recipe_entity_products
										  WHERE recipe_entity_id = (SELECT id FROM recipes WHERE recipe_name = :recipe_name))
				
				SELECT recipe_name,
					   image_url,
					   COALESCE(calories, 0)                   AS calories,
					   preparation_time,
					   COALESCE(protein, 0)                    AS protein,
					   difficulty,
					   (SELECT steps FROM steps_results)       AS steps,
					   (SELECT products FROM products_results) AS products,
					   category,
					   users.id                                AS ownerId,
					   users.username                          AS owner_name
				FROM recipes
						 LEFT JOIN recipe_entity_products ON recipes.id = recipe_entity_products.recipe_entity_id
						 LEFT JOIN recipe_entity_steps ON recipes.id = recipe_entity_steps.recipe_entity_id
						 LEFT JOIN users ON users.id = recipes.owner_id
				WHERE recipe_name = :recipe_name;`,
		map[string]interface{}{"recipe_name": recipeName},
	)
	return
}

// GetRecipesFromUser gets the recipes created from the given user
func GetRecipesFromUser(username string) (recipes []BaseRecipeInfo, err error) {
	err = database.GetMultipleRecordsNamedQuery(
		&recipes,
		`SELECT recipe_name,
					   image_url
				FROM recipes
						 JOIN users ON recipes.owner_id = users.id
				WHERE username = :username
				ORDER BY visitations_count DESC;`,
		map[string]interface{}{"username": username},
	)
	return
}
