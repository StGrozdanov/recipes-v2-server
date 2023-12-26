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
