package notifications

import "recipes-v2-server/database"

// GetNotificationsForUser get all unread notifications of the give user
func GetNotificationsForUser(username string) (notifications []Notification, err error) {
	err = database.GetMultipleRecordsNamedQuery(
		&notifications,
		`SELECT notifications.id,
    				   action,
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

// MarkAsRead marks the given notification as read
func MarkAsRead(id int) (err error) {
	_, err = database.ExecuteNamedQuery(
		`UPDATE notifications SET is_marked_as_read = true WHERE notifications.id = :id;`,
		map[string]interface{}{"id": id},
	)
	return
}
