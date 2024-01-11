package notifications

type Notification struct {
	Id             int    `json:"id" db:"id"`
	SenderAvatar   string `json:"senderAvatar" db:"sender_avatar"`
	SenderUsername string `json:"senderUsername" db:"sender_username"`
	Action         string `json:"action" db:"action"`
	LocationName   string `json:"locationName" db:"location_name"`
	CreatedAt      string `json:"createdAt" db:"created_at"`
}

type NotificationMarkAsReadData struct {
	Id int `json:"id" db:"id" valid:"required,int"`
}

type NotificationRequest struct {
	SenderAvatar   string `json:"senderAvatar" db:"sender_avatar"`
	SenderUsername string `json:"senderUsername" db:"sender_username" valid:"required"`
	SenderId       int    `json:"senderId" db:"sender_id" valid:"required"`
	Action         string `json:"action" db:"action" valid:"required"`
	LocationName   string `json:"locationName" db:"location_name" valid:"required"`
}
