package notifications

import "recipes-v2-server/database"

// GetNotificationsForUser get all unread notifications of the give user
func GetNotificationsForUser(username string) (notifications []Notification, err error) {
	err = database.GetMultipleRecordsNamedQuery(
		&notifications,
		`SELECT action,
					   created_at,
					   location_name,
					   sender_avatar,
					   sender_username
				FROM notifications
						 JOIN users ON notifications.receiver_id = users.id
				WHERE is_marked_as_read = false AND username = :username;`,
		map[string]interface{}{"username": username},
	)
	return
}
