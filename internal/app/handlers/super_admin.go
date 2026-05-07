package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
)

// CreateUserAPI creates a new user (super admin only).
func CreateUserAPI(c *gin.Context) {
	req := &models.CreateUserRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse("请填写用户名和密码"))
		return
	}

	userRepo := repository.GetUserRepo(db.GetPostgresPool())
	u := &models.User{Username: req.Username, RawPassword: req.Password}
	if err := userRepo.Save(u); err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(models.ToUserData(u)))
}

// ListUsersAPI lists all users with pagination and search (super admin only).
func ListUsersAPI(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("pageSize", "20"), 10, 64)
	query := c.Query("query")

	userRepo := repository.GetUserRepo(db.GetPostgresPool())
	result, err := userRepo.List(repository.ListOptions{
		Query: query, Page: page, PageSize: pageSize,
	})
	if err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
		return
	}

	users := make([]*models.UserData, 0, len(result.Users))
	for _, u := range result.Users {
		users = append(users, models.ToUserData(u))
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(models.ListUsersResponseData{
		Users: users, Total: result.Total,
	}))
}

// GetUserAPI gets a user's detail including tenants and roles (super admin only).
func GetUserAPI(c *gin.Context) {
	id := c.Param("id")
	userRepo := repository.GetUserRepo(db.GetPostgresPool())

	u, err := userRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse("用户不存在"))
		return
	}

	tenants, _ := userRepo.GetUserTenants(u)
	c.JSON(http.StatusOK, models.NewSuccessResponse(models.UserDetailData{
		ID: u.ID, Username: u.Username, IsActive: u.IsActive,
		IsSuper: u.IsSuper, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
		Tenants: tenants,
	}))
}

// ResetPasswordAPI resets a user's password (super admin only).
func ResetPasswordAPI(c *gin.Context) {
	id := c.Param("id")
	req := &models.ResetPasswordRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse("请填写新密码"))
		return
	}

	userRepo := repository.GetUserRepo(db.GetPostgresPool())
	if err := userRepo.UpdatePasswordByID(id, req.NewPassword); err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
}

// UpdateUserStatusAPI enables or disables a user (super admin only).
func UpdateUserStatusAPI(c *gin.Context) {
	id := c.Param("id")
	req := &models.UpdateUserStatusRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse("请指定用户状态"))
		return
	}

	userRepo := repository.GetUserRepo(db.GetPostgresPool())
	if err := userRepo.UpdateStatus(id, req.IsActive); err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
}
