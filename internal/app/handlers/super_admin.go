package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
)

// CreateUserAPI godoc
// @Summary 创建用户
// @Description 超级管理员创建普通用户
// @Tags 超级管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body models.CreateUserRequest true "创建用户请求"
// @Success 200 {object} models.Response{data=models.UserData}
// @Failure 401 {object} nil
// @Router /super/users [post]
func CreateUserAPI(c *gin.Context) {
	req := &models.CreateUserRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse("请填写用户名和密码"))
		return
	}

	userRepo := repository.GetUserRepo(db.GetPostgresPool())
	u := &models.User{Username: strings.TrimSpace(req.Username), RawPassword: req.Password}
	if err := userRepo.Save(u); err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(models.ToUserData(u)))
}

// ListUsersAPI godoc
// @Summary 用户列表
// @Description 超级管理员查看所有用户，支持分页和搜索
// @Tags 超级管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Param query query string false "搜索用户名"
// @Success 200 {object} models.Response{data=models.ListUsersResponseData}
// @Failure 401 {object} nil
// @Router /super/users [get]
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

// GetUserAPI godoc
// @Summary 获取用户详情
// @Description 超级管理员查看用户详情，包括其所属的所有租户和角色
// @Tags 超级管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户 ID"
// @Success 200 {object} models.Response{data=models.UserDetailData}
// @Failure 401 {object} nil
// @Router /super/users/{id} [get]
func GetUserAPI(c *gin.Context) {
	id := c.Param("id")
	userRepo := repository.GetUserRepo(db.GetPostgresPool())

	u, err := userRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse("用户不存在"))
		return
	}

	tenants, err := userRepo.GetUserTenants(u)
	if err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.NewSuccessResponse(models.UserDetailData{
		ID: u.ID, Username: u.Username, IsActive: u.IsActive,
		IsSuper: u.IsSuper, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
		Tenants: tenants,
	}))
}

// ResetPasswordAPI godoc
// @Summary 重置用户密码
// @Description 超级管理员重置指定用户的密码
// @Tags 超级管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户 ID"
// @Param body body models.ResetPasswordRequest true "重置密码请求"
// @Success 200 {object} models.Response
// @Failure 401 {object} nil
// @Router /super/users/{id}/reset-password [post]
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

// UpdateUserStatusAPI godoc
// @Summary 禁用/启用用户
// @Description 超级管理员启用或禁用指定用户
// @Tags 超级管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户 ID"
// @Param body body models.UpdateUserStatusRequest true "更新用户状态请求"
// @Success 200 {object} models.Response
// @Failure 401 {object} nil
// @Router /super/users/{id}/status [patch]
func UpdateUserStatusAPI(c *gin.Context) {
	id := c.Param("id")
	req := &models.UpdateUserStatusRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse("请指定用户状态"))
		return
	}

	userRepo := repository.GetUserRepo(db.GetPostgresPool())
	if err := userRepo.UpdateStatus(id, *req.IsActive); err != nil {
		c.JSON(http.StatusOK, models.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
}
