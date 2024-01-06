package recipes

import (
	"bytes"
	"errors"
	"mime/multipart"
	"recipes-v2-server/database"
	"recipes-v2-server/utils"
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
	if err != nil {
		return
	}

	var totalRecipesCount int
	err = database.GetSingleRecord(&totalRecipesCount, `SELECT COUNT(id) FROM recipes;`)
	if err != nil {
		return
	}

	if cursor == 0 || cursor-limit < 0 {
		recipes.PageData.PrevPage = 0
	} else {
		recipes.PageData.PrevPage = cursor - limit
	}

	recipes.PageData.NextPage = limit + cursor
	recipes.PageData.FirstPage = cursor == 0
	recipes.PageData.LastPage = limit+cursor >= totalRecipesCount

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
		`SELECT recipe_name,
					   image_url,
					   COALESCE(calories, 0)                   AS calories,
					   preparation_time,
					   COALESCE(protein, 0)                    AS protein,
					   difficulty,
					   steps,
					   products,
					   category,
					   users.id                                AS owner_id,
					   users.username                          AS owner_name
				FROM recipes
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

// GetFavourites gets the recipes favourite of the given user
func GetFavourites(username string) (recipes []BaseRecipeInfo, err error) {
	err = database.GetMultipleRecordsNamedQuery(
		&recipes,
		`SELECT recipe_name,
					   image_url
				FROM users
						 JOIN users_favourites ON users_favourites.user_entity_id = users.id
						 JOIN recipes ON recipes.id = users_favourites.favourites_id
				WHERE username = :username
				ORDER BY visitations_count DESC;`,
		map[string]interface{}{"username": username},
	)
	return
}

// IsInFavourites checks if recipe is in user favourites and returns a boolean value
func IsInFavourites(data FavouritesRequest) (isInFavourites bool, err error) {
	err = database.GetSingleRecordNamedQuery(
		&isInFavourites,
		`SELECT EXISTS(SELECT recipe_name
              FROM recipes
                       JOIN users_favourites
                            ON recipes.id = users_favourites.favourites_id
                                AND users_favourites.user_entity_id = :user_id
              WHERE recipe_name = :recipe_name);`,
		data,
	)
	return
}

// AddToFavourites adds and recipe to the user favourites collection
func AddToFavourites(data FavouritesRequest) (err error) {
	_, err = database.ExecuteNamedQuery(
		`WITH recipes AS (SELECT id FROM recipes WHERE recipe_name = :recipe_name)
				
				INSERT
				INTO users_favourites(user_entity_id, favourites_id)
				VALUES (:user_id, (SELECT recipes.id FROM recipes));`,
		data,
	)
	return
}

// RemoveFromFavourites removes a recipe from user favourites collection
func RemoveFromFavourites(data FavouritesRequest) (err error) {
	_, err = database.ExecuteNamedQuery(
		`WITH recipes AS (SELECT id FROM recipes WHERE recipe_name = :recipe_name)

				DELETE
				FROM users_favourites
				WHERE user_entity_id = :user_id
				  AND favourites_id = (SELECT id FROM recipes);`,
		data,
	)
	return
}

// Create creates a new recipe
func Create(recipe RecipeData, authToken string) (response RecipeData, err error) {
	recipe, err = adjustRecipeStatus(recipe, authToken)
	if err != nil {
		return
	}

	err = database.GetSingleRecordNamedQuery(
		&response,
		`INSERT INTO recipes (category,
                     created_at,
                     image_url,
                     owner_id,
                     recipe_name,
                     status,
                     visitations_count,
                     calories,
                     protein,
                     preparation_time,
                     difficulty,
                     steps,
                     products)
				VALUES (:category,
						NOW(),
						:image_url,
						:owner_id,
						:recipe_name,
						:status,
						0,
						:calories,
						:protein,
						:preparation_time,
						:difficulty,
						:steps,
						:products)
				RETURNING recipe_name,
					image_url,
					COALESCE(calories, 0) AS calories,
					preparation_time,
					COALESCE(protein, 0) AS protein,
					difficulty,
					steps,
					products,
					category,
					owner_id;`,
		recipe,
	)

	response.OwnerData.Username = recipe.OwnerData.Username
	return
}

func adjustRecipeStatus(recipe RecipeData, authToken string) (RecipeData, error) {
	claims, isValid, err := utils.ParseJWT(authToken)
	if err != nil {
		return recipe, err
	}
	if !isValid {
		return recipe, errors.New("invalid token")
	}

	if claims.Role == "ADMINISTRATOR" || claims.Role == "MODERATOR" {
		recipe.Status = "APPROVED"
	} else {
		recipe.Status = "PENDING"
	}
	return recipe, nil
}

// RecipeNameExists checks for existing recipe with this name and returns boolean value
func RecipeNameExists(recipeName string) (exists bool, err error) {
	err = database.GetSingleRecordNamedQuery(
		&exists,
		`SELECT EXISTS(SELECT id FROM recipes WHERE recipe_name = :recipe_name);`,
		map[string]interface{}{"recipe_name": recipeName},
	)
	return
}

// UploadRecipeImage Uploads recipe image to s3 bucket and returns the URL
func UploadRecipeImage(file *multipart.FileHeader, fileKey string) (imageURL string, err error) {
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
	return
}

// Edit edits a recipe
func Edit(recipeName string, data RecipeData) (result RecipeData, err error) {
	extendedData := ExtendedRecipeData{data, recipeName}

	err = database.GetSingleRecordNamedQuery(
		&result,
		`UPDATE recipes
				SET recipe_name      = :recipe_name,
					preparation_time = :preparation_time,
					category         = :category,
					image_url        = :image_url,
					owner_id         = :owner_id,
					calories         = :calories,
					protein          = :protein,
					difficulty       = :difficulty,
					steps            = :steps,
					products         = :products
				WHERE recipe_name = :old_recipe_name
				RETURNING *`,
		extendedData,
	)
	return
}
