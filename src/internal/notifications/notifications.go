package notifications

import (
	"fmt"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"recipes-v2-server/database"
	"recipes-v2-server/utils"
)

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

// Create creates a notification to the related receivers
func Create(request NotificationRequest) (err error) {
	const query = `INSERT INTO notifications (action,
                           created_at,
                           location_id,
                           location_name,
                           is_marked_as_read,
                           receiver_id,
                           sender_avatar,
                           sender_username)
					SELECT :action,
						   Now(),
						   recipes.id,
						   :recipe_name,
						   false,
						   :receiver_id,
						   :sender_avatar,
						   :sender_username
					FROM recipes
					WHERE recipe_name = :recipe_name;`

	queryParams := map[string]interface{}{
		"action":          request.Action,
		"recipe_name":     request.LocationName,
		"sender_avatar":   request.SenderAvatar,
		"sender_username": request.SenderUsername,
	}

	receiverIds, err := findNotificationReceivers(request)
	if err != nil {
		return
	}

	statement, err := database.PrepareNamedStatement(query)
	if err != nil {
		return
	}

	for _, receiverId := range receiverIds {
		queryParams["receiver_id"] = receiverId

		if _, err := statement.Exec(queryParams); err != nil {
			utils.
				GetLogger().
				WithFields(log.Fields{"error": err.Error(), "request": request}).
				Error("Error on executing statement for creating notification")
		}
	}

	return
}

func findNotificationReceivers(request NotificationRequest) (results pq.Int32Array, err error) {
	var notificationActionsReceiversMap = map[string]func(request NotificationRequest) (pq.Int32Array, error){
		"EDITED_RECIPE":   findEditRecipeActionReceivers,
		"CREATED_RECIPE":  findCreateRecipeActionReceivers,
		"DELETED_RECIPE":  findDeleteRecipeActionReceivers,
		"CREATED_COMMENT": findCreateCommentActionReceivers,
		"EDITED_COMMENT":  findCreateCommentActionReceivers,
		"DELETED_COMMENT": findEditRecipeActionReceivers,
	}

	handlerFunc, found := notificationActionsReceiversMap[request.Action]
	if !found {
		err = fmt.Errorf("such action is not registered - %s", request.Action)
		return
	}

	return handlerFunc(request)
}

func findEditRecipeActionReceivers(request NotificationRequest) (results pq.Int32Array, err error) {
	err = database.GetSingleRecordNamedQuery(
		&results,
		`WITH admin_and_moderator_groups AS (SELECT id
                                    FROM users
                                             JOIN users_roles ON user_entity_id = users.id
                                    WHERE roles_id IN (1, 2)
                                      AND id != :sender_id),

				 resource_owner_that_is_not_the_sender AS (SELECT users.id
														   FROM users
																	JOIN recipes ON owner_id = users.id
														   WHERE recipes.recipe_name = :location_name
															 AND users.id != :sender_id)
			
			SELECT ARRAY(SELECT id
						 FROM admin_and_moderator_groups
			
						 UNION
			
						 SELECT id
						 FROM resource_owner_that_is_not_the_sender
				   ) AS results;`,
		request,
	)
	return
}

func findCreateRecipeActionReceivers(request NotificationRequest) (results pq.Int32Array, err error) {
	err = database.GetSingleRecordNamedQuery(
		&results,
		`SELECT ARRAY(SELECT id
                            FROM users
                                JOIN users_roles ON user_entity_id = users.id
                            WHERE roles_id = 1
                                AND id != :sender_id) AS results;`,
		request,
	)
	return
}

func findDeleteRecipeActionReceivers(request NotificationRequest) (results pq.Int32Array, err error) {
	err = database.GetSingleRecordNamedQuery(
		&results,
		`WITH admin_groups AS (SELECT id
                                    FROM users
                                             JOIN users_roles ON user_entity_id = users.id
                                    WHERE roles_id = 1
                                      AND id != :sender_id),

				 resource_owner_that_is_not_the_sender AS (SELECT users.id
														   FROM users
																	JOIN recipes ON owner_id = users.id
														   WHERE recipes.recipe_name = :location_name
															 AND users.id != :sender_id)
			
			SELECT ARRAY(SELECT id
						 FROM admin_and_moderator_groups
			
						 UNION
			
						 SELECT id
						 FROM resource_owner_that_is_not_the_sender
				   ) AS results;`,
		request,
	)
	return
}

func findCreateCommentActionReceivers(request NotificationRequest) (results pq.Int32Array, err error) {
	err = database.GetSingleRecordNamedQuery(
		&results,
		`WITH admin_and_moderator_groups AS (SELECT id
                                    FROM users
                                             JOIN users_roles ON user_entity_id = users.id
                                    WHERE roles_id IN (1, 2)
                                      AND id != :sender_id),

					 resource_owner_that_is_not_the_sender AS (SELECT users.id
															   FROM users
																		JOIN recipes ON owner_id = users.id
															   WHERE recipes.recipe_name = :location_name
																 AND users.id != :sender_id),
				
					 users_involved_into_the_conversation AS (SELECT comments.owner_id
															  FROM comments
																	   JOIN recipes ON comments.target_recipe_id = recipes.id
															  WHERE recipe_name = :location_name
																AND comments.owner_id != :sender_id)
				
				SELECT ARRAY(SELECT id
							 FROM admin_and_moderator_groups
				
							 UNION
				
							 SELECT id
							 FROM resource_owner_that_is_not_the_sender
				
							 UNION
				
							 SELECT owner_id
							 FROM users_involved_into_the_conversation
					   ) AS results;`,
		request,
	)
	return
}
