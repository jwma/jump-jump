package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/i18n"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
)

// helper: check if the current user is admin of the given tenant (or is super admin).
func isTenantAdmin(user *models.User, tenantID string) (*models.TenantMember, error) {
	if user.IsSuper {
		return &models.TenantMember{TenantID: tenantID, UserID: user.ID, Role: models.RoleAdmin}, nil
	}
	memberRepo := repository.GetTenantMemberRepo(db.GetPostgresPool())
	return memberRepo.Get(tenantID, user.ID)
}

// helper: resolve member for a tenant; returns nil (not a member) rather than error for non-members.
func getMemberOrNil(user *models.User, tenantID string) *models.TenantMember {
	if user.IsSuper {
		return &models.TenantMember{TenantID: tenantID, UserID: user.ID, Role: models.RoleAdmin}
	}
	memberRepo := repository.GetTenantMemberRepo(db.GetPostgresPool())
	m, err := memberRepo.Get(tenantID, user.ID)
	if err != nil {
		return nil
	}
	return m
}

// ListTenantMembersAPI godoc
// @Summary 租户成员列表
// @Description 查看租户成员列表（租户内所有成员均可查看）
// @Tags 租户成员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "租户 ID"
// @Success 200 {object} models.Response{data=[]models.TenantMemberData}
// @Failure 401 {object} nil
// @Router /tenant/{id}/members [get]
func ListTenantMembersAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		tenantID := c.Param("id")
		member := getMemberOrNil(ctx.User, tenantID)
		if member == nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "tenant.notMember")))
			return
		}

		memberRepo := repository.GetTenantMemberRepo(db.GetPostgresPool())
		members, err := memberRepo.ListByTenant(tenantID)
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "tenant.memberListFailed")))
			return
		}

		userRepo := repository.GetUserRepo(db.GetPostgresPool())
		result := make([]*models.TenantMemberData, 0, len(members))
		for _, m := range members {
			u, err := userRepo.FindByID(m.UserID)
			if err != nil {
				continue
			}
			result = append(result, &models.TenantMemberData{
				UserID:   m.UserID,
				Username: u.Username,
				Role:     m.Role,
				JoinedAt: m.JoinedAt,
			})
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(result))
	})
}

// InviteUserAPI godoc
// @Summary 邀请用户加入租户
// @Description 邀请用户加入租户（仅 admin）
// @Tags 租户成员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "租户 ID"
// @Param body body models.InviteUserRequest true "邀请请求"
// @Success 200 {object} models.Response{data=models.TenantInvitation}
// @Failure 401 {object} nil
// @Router /tenant/{id}/invitations [post]
func InviteUserAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		tenantID := c.Param("id")
		member, err := isTenantAdmin(ctx.User, tenantID)
		if err != nil || !member.IsAdmin() {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.adminOnlyInvite")))
			return
		}

		req := &models.InviteUserRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "common.invalidParameters")))
			return
		}

		userRepo := repository.GetUserRepo(db.GetPostgresPool())
		invitee, err := userRepo.FindByUsername(req.Username)
		if err != nil {
			writeErrorResponse(c, err)
			return
		}

		memberRepo := repository.GetTenantMemberRepo(db.GetPostgresPool())
		isMember, _ := memberRepo.IsMember(tenantID, invitee.ID)
		if isMember {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.alreadyMember")))
			return
		}

		if invitee.ID == ctx.User.ID {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.cannotInviteSelf")))
			return
		}

		invRepo := repository.GetTenantInvitationRepo(db.GetPostgresPool())
		hasPending, _ := invRepo.HasPendingInvitation(tenantID, invitee.ID)
		if hasPending {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.pendingInvitationExists")))
			return
		}

		inv := &models.TenantInvitation{
			TenantID:  tenantID,
			InviterID: ctx.User.ID,
			InviteeID: invitee.ID,
		}
		if err := invRepo.Create(inv); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.invitationCreationFailed")))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(inv))
	})
}

// ListMyInvitationsAPI godoc
// @Summary 查看我的邀请
// @Description 查看当前用户收到的所有待处理邀请
// @Tags 邀请
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=[]models.InvitationData}
// @Failure 401 {object} nil
// @Router /invitations [get]
func ListMyInvitationsAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		invRepo := repository.GetTenantInvitationRepo(db.GetPostgresPool())
		invs, err := invRepo.FindPendingByInvitee(ctx.User.ID)
		if err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.invitationListFailed")))
			return
		}

		userRepo := repository.GetUserRepo(db.GetPostgresPool())
		tenantRepo := repository.GetTenantRepo(db.GetPostgresPool())
		result := make([]*models.InvitationData, 0, len(invs))
		for _, inv := range invs {
			t, err := tenantRepo.GetByID(inv.TenantID)
			if err != nil {
				continue
			}
			inviter, err := userRepo.FindByID(inv.InviterID)
			if err != nil {
				continue
			}
			result = append(result, &models.InvitationData{
				ID:              inv.ID,
				TenantID:        inv.TenantID,
				TenantName:      t.Name,
				InviterUsername: inviter.Username,
				Status:          inv.Status,
				CreatedAt:       inv.CreatedAt,
			})
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(result))
	})
}

// AcceptInvitationAPI godoc
// @Summary 接受邀请
// @Description 接受租户邀请
// @Tags 邀请
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "邀请 ID"
// @Success 200 {object} models.Response
// @Failure 401 {object} nil
// @Router /invitations/{id}/accept [post]
func AcceptInvitationAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		invID := c.Param("id")
		invRepo := repository.GetTenantInvitationRepo(db.GetPostgresPool())

		inv, err := invRepo.Get(invID)
		if err != nil {
			writeErrorResponse(c, err)
			return
		}

		if inv.InviteeID != ctx.User.ID {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "invitation.noPermission")))
			return
		}

		if inv.Status != models.InvitationStatusPending {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "invitation.alreadyProcessed")))
			return
		}

		memberRepo := repository.GetTenantMemberRepo(db.GetPostgresPool())
		if err := memberRepo.Save(&models.TenantMember{
			TenantID: inv.TenantID,
			UserID:   ctx.User.ID,
			Role:     models.RoleMember,
		}); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "invitation.joinFailed")))
			return
		}

		if err := invRepo.UpdateStatus(inv.ID, models.InvitationStatusAccepted); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "common.operationFailed")))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
	})
}

// RejectInvitationAPI godoc
// @Summary 拒绝邀请
// @Description 拒绝租户邀请
// @Tags 邀请
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "邀请 ID"
// @Success 200 {object} models.Response
// @Failure 401 {object} nil
// @Router /invitations/{id}/reject [post]
func RejectInvitationAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		invID := c.Param("id")
		invRepo := repository.GetTenantInvitationRepo(db.GetPostgresPool())

		inv, err := invRepo.Get(invID)
		if err != nil {
			writeErrorResponse(c, err)
			return
		}

		if inv.InviteeID != ctx.User.ID {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "invitation.noPermission")))
			return
		}

		if inv.Status != models.InvitationStatusPending {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "invitation.alreadyProcessed")))
			return
		}

		if err := invRepo.UpdateStatus(inv.ID, models.InvitationStatusRejected); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "common.operationFailed")))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
	})
}

// UpdateMemberRoleAPI godoc
// @Summary 修改成员角色
// @Description 修改租户成员角色（仅 admin，不能修改自己的角色）
// @Tags 租户成员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "租户 ID"
// @Param userId path string true "用户 ID"
// @Param body body models.UpdateRoleRequest true "角色修改请求"
// @Success 200 {object} models.Response
// @Failure 401 {object} nil
// @Router /tenant/{id}/members/{userId}/role [patch]
func UpdateMemberRoleAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		tenantID := c.Param("id")
		targetUserID := c.Param("userId")

		member, err := isTenantAdmin(ctx.User, tenantID)
		if err != nil || !member.IsAdmin() {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.adminOnlyChangeRole")))
			return
		}

		if targetUserID == ctx.User.ID {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.cannotChangeOwnRole")))
			return
		}

		req := &models.UpdateRoleRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "common.invalidParameters")))
			return
		}

		if req.Role != models.RoleAdmin && req.Role != models.RoleMember {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.invalidRole")))
			return
		}

		memberRepo := repository.GetTenantMemberRepo(db.GetPostgresPool())
		isMember, _ := memberRepo.IsMember(tenantID, targetUserID)
		if !isMember {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.userNotMember")))
			return
		}

		if err := memberRepo.UpdateRole(tenantID, targetUserID, req.Role); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.roleChangeFailed")))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
	})
}

// RemoveMemberAPI godoc
// @Summary 移除成员
// @Description 移除租户成员（仅 admin，不能移除自己）
// @Tags 租户成员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "租户 ID"
// @Param userId path string true "用户 ID"
// @Success 200 {object} models.Response
// @Failure 401 {object} nil
// @Router /tenant/{id}/members/{userId} [delete]
func RemoveMemberAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		tenantID := c.Param("id")
		targetUserID := c.Param("userId")

		member, err := isTenantAdmin(ctx.User, tenantID)
		if err != nil || !member.IsAdmin() {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.adminOnlyRemove")))
			return
		}

		if targetUserID == ctx.User.ID {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.cannotRemoveSelf")))
			return
		}

		memberRepo := repository.GetTenantMemberRepo(db.GetPostgresPool())
		isMember, _ := memberRepo.IsMember(tenantID, targetUserID)
		if !isMember {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.userNotMember")))
			return
		}

		if err := memberRepo.Delete(tenantID, targetUserID); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.removeFailed")))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
	})
}

// LeaveTenantAPI godoc
// @Summary 退出租户
// @Description 退出租户（任何成员均可操作）
// @Tags 租户成员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "租户 ID"
// @Success 200 {object} models.Response
// @Failure 401 {object} nil
// @Router /tenant/{id}/leave [post]
func LeaveTenantAPI() gin.HandlerFunc {
	return Authenticator(func(c *gin.Context, ctx *AuthContext) {
		tenantID := c.Param("id")
		member := getMemberOrNil(ctx.User, tenantID)
		if member == nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "tenant.notMember")))
			return
		}

		memberRepo := repository.GetTenantMemberRepo(db.GetPostgresPool())
		if err := memberRepo.Delete(tenantID, ctx.User.ID); err != nil {
			c.JSON(http.StatusOK, models.NewErrorResponse(i18n.T(c, "member.leaveFailed")))
			return
		}

		c.JSON(http.StatusOK, models.NewSuccessResponse(nil))
	})
}
