package response

type UserResponseFull struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	AvatarURL      string `json:"avatar_url"`
	Bio            string `json:"bio"`
	FollowersCount uint   `json:"followers_count"`
	FollowingCount uint   `json:"following_count"`
}

type UserResponseShort struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type PaginatedUsersResponse struct {
	Users []UserResponseShort `json:"users"`
	Total int                 `json:"total"`
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
}
