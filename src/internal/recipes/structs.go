package recipes

import (
	"github.com/lib/pq"
	"recipes-v2-server/internal/users"
)

type ExtendedRecipeInfo struct {
	ImageURL   string `json:"imageURL" db:"image_url"`
	RecipeName string `json:"recipeName" db:"recipe_name"`
	Category   string `json:"category" db:"category"`
}

type BaseRecipeInfo struct {
	ImageURL   string `json:"imageURL" db:"image_url"`
	RecipeName string `json:"recipeName" db:"recipe_name"`
}

type BaseRecipeInfoArray = []BaseRecipeInfo

type RecipePaginationInfo struct {
	BaseRecipeInfoArray `json:"pages"`
	PageData            `json:"pageData"`
}

type PageData struct {
	PrevPage  int  `json:"prevPage"`
	NextPage  int  `json:"nextPage"`
	FirstPage bool `json:"firstPage"`
	LastPage  bool `json:"lastPage"`
}

type RecipeData struct {
	RecipeName      string         `db:"recipe_name" json:"recipeName"`
	Products        pq.StringArray `db:"products" json:"products"`
	Steps           pq.StringArray `db:"steps" json:"steps"`
	ImageURL        string         `db:"image_url" json:"imageURL"`
	CategoryName    string         `db:"category" json:"category"`
	Difficulty      string         `db:"difficulty" json:"difficulty"`
	PreparationTime int            `db:"preparation_time" json:"preparationTime"`
	Calories        int            `db:"calories" json:"calories"`
	Protein         int            `db:"protein" json:"protein"`
	users.OwnerData `json:"owner"`
}
