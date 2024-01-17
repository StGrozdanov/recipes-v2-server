package analytics

import (
	"recipes-v2-server/database"
)

// VisitationsForTheLastSixMonths returns the visitations for the last six months, in ready to use format for
// react chart
func VisitationsForTheLastSixMonths() (result ChartData, err error) {
	var data []VisitationsData

	err = database.GetMultipleRecords(
		&data,
		`
		SELECT TO_CHAR(DATE_TRUNC('month', visited_at), 'Month') 	  AS date,
			   COUNT(id)                                              AS visitations
		FROM visitations
		WHERE DATE_TRUNC('month', visited_at) > DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '6 months'
		GROUP BY DATE_TRUNC('month', visited_at)
		ORDER BY EXTRACT(YEAR FROM DATE_TRUNC('month', visited_at)),
				 EXTRACT(MONTH FROM DATE_TRUNC('month', visited_at));`,
	)

	result.Datasets = append(result.Datasets, Data{})

	for _, visitationData := range data {
		result.Labels = append(result.Labels, visitationData.Date)
		result.Datasets[0].Data = append(result.Datasets[0].Data, visitationData.Visitations)
	}
	return
}

// VisitationsForToday returns the visitations count for the current day
func VisitationsForToday() (visitations int, err error) {
	err = database.GetSingleRecord(
		&visitations,
		`
		SELECT COUNT(id)
		FROM visitations
		WHERE visited_at >= Now() - INTERVAL '1 day';`,
	)
	return
}

// MostActiveUser determines the most active user by his total publications count
func MostActiveUser() (user MostActiveUserData, err error) {
	err = database.GetSingleRecord(
		&user,
		`
		WITH users_recipes AS (SELECT username,
                              avatar_url,
                              COUNT(recipes.id) AS recipes_count
                       FROM users
                                LEFT JOIN recipes ON recipes.owner_id = users.id
                       GROUP BY users.username, avatar_url),

			 users_comments AS (SELECT username,
									   COUNT(comments.id) AS comments_count
								FROM users
										 LEFT JOIN comments ON comments.owner_id = users.id
								GROUP BY users.username)
		
		SELECT users_recipes.username,
			   users_recipes.avatar_url,
			   recipes_count,
			   comments_count,
			   SUM(recipes_count + comments_count) AS total_publications_count
		FROM users_recipes
				 LEFT JOIN users_comments ON users_recipes.username = users_comments.username
		GROUP BY users_recipes.username, users_recipes.avatar_url, users_recipes.recipes_count, users_comments.comments_count
		ORDER BY total_publications_count DESC
		LIMIT 1;`,
	)
	return
}
