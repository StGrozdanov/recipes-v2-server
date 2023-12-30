package notifications

type Notification struct {
	SenderAvatar   string `json:"senderAvatar" db:"sender_avatar"`
	SenderUsername string `json:"senderUsername" db:"sender_username"`
	Action         string `json:"action" db:"action"`
	LocationName   string `json:"locationName" db:"location_name"`
	CreatedAt      string `json:"createdAt" db:"created_at"`
}
