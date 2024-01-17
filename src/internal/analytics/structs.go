package analytics

type Data struct {
	Data []int `json:"data"`
}

type ChartData struct {
	Labels   []string `json:"labels"`
	Datasets []Data   `json:"datasets"`
}

type VisitationsData struct {
	Date        string `json:"date" db:"date"`
	Visitations int    `json:"visitations" db:"visitations"`
}

type MostActiveUserData struct {
	Username               string `json:"username" db:"username"`
	AvatarURL              string `json:"avatarURL" db:"avatar_url"`
	RecipesCount           int    `json:"recipesCount" db:"recipes_count"`
	CommentsCount          int    `json:"commentsCount" db:"comments_count"`
	TotalPublicationsCount int    `json:"totalPublicationsCount" db:"total_publications_count"`
}
