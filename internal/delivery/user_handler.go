package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user_service/internal/service"
	"user_service/internal/transport/request"
)

type UserHandler struct {
	s *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{s: s}
}

func (h *UserHandler) GetUserByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid id")
		return
	}

	res, err := h.s.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Could not get user")
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var req request.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid data")
		return
	}

	res, err := h.s.CreateUser(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Could not create user")
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	var req request.UpdateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	requestUserIdStr := ctx.Param("id")
	requestUserIdUint64, err := strconv.ParseUint(requestUserIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid params"})
		return
	}
	requestUserId := uint(requestUserIdUint64)

	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "You're unauthorized"})
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user id in context"})
		return
	}

	res, err := h.s.UpdateUser(userID, req, requestUserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	requestUserIdStr := ctx.Param("id")
	requestUserIdUint64, err := strconv.ParseUint(requestUserIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID parameter"})
		return
	}
	requestUserID := uint(requestUserIdUint64)

	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDUint64, ok := userIDVal.(uint64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}
	userID := uint(userIDUint64)

	if err := h.s.DeleteUser(userID, requestUserID); err != nil {
		if err.Error() == "unauthorized: users can only delete their own account" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own account"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) GetUsersPaginated(ctx *gin.Context) {

	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")

	page := 1
	pageSize := 10

	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
			return
		}
	}

	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page_size parameter"})
			return
		}
	}

	resp, err := h.s.GetUsersPaginated(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
