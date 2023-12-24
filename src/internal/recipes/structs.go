package recipes

type ExtendedRecipeInfo struct {
	ImageURL   string `json:"imageURL" db:"image_url"`
	RecipeName string `json:"recipeName" db:"recipe_name"`
	Category   string `json:"category" db:"category"`
}

type BaseRecipeInfo struct {
	ImageURL   string `json:"imageURL" db:"image_url"`
	RecipeName string `json:"recipeName" db:"recipe_name"`
}
