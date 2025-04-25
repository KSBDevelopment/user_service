package request

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Bio      string `json:"bio"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
}
