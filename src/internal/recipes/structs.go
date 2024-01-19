package recipes

import (
	"encoding/json"
	"recipes-v2-server/internal/users"
)

type ExtendedRecipeInfo struct {
	ImageURL   string `json:"imageURL" db:"image_url"`
	RecipeName string `json:"recipeName" db:"recipe_name"`
	Category   string `json:"category" db:"category"`
}

type BaseRecipeInfo struct {
	ImageURL   string `json:"imageURL" db:"image_url"`
	RecipeName string `json:"recipeName" db:"recipe_name" valid:"required"`
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
	RecipeName      string          `db:"recipe_name" json:"recipeName" valid:"required,minstringlength(4)"`
	Products        json.RawMessage `db:"products" json:"products" valid:"required"`
	Steps           json.RawMessage `db:"steps" json:"steps" valid:"required"`
	ImageURL        string          `db:"image_url" json:"imageURL" valid:"required,url"`
	CategoryName    string          `db:"category" json:"category" valid:"required"`
	Difficulty      string          `db:"difficulty" json:"difficulty" valid:"required"`
	PreparationTime int             `db:"preparation_time" json:"preparationTime" valid:"required"`
	Calories        int             `db:"calories" json:"calories"`
	Protein         int             `db:"protein" json:"protein"`
	Status          string          `db:"status" json:"-"`
	users.OwnerData `json:"owner"`
}

type FavouritesRequest struct {
	RecipeName string `json:"recipeName" db:"recipe_name" valid:"required"`
	UserId     int    `json:"userId" db:"user_id" valid:"required"`
}

type ExtendedRecipeData struct {
	RecipeData
	OriginalRecipeName string `db:"old_recipe_name"`
}

type AdminRecipeData struct {
	RecipeName string `db:"recipe_name" json:"recipeName" valid:"required,minstringlength(4)"`
	ImageURL   string `db:"image_url" json:"imageURL" valid:"required,url"`
	Status     string `db:"status" json:"status"`
	OwnerName  string `db:"owner_name" json:"ownerName" valid:"required"`
}
