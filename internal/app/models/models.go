package models

import (
	"encoding/json"
	"time"
)

// Role constants (string-based for tenant_members)
const RoleAdmin = "admin"
const RoleMember = "member"

// Invitation status constants
const InvitationStatusPending = "pending"
const InvitationStatusAccepted = "accepted"
const InvitationStatusRejected = "rejected"

type Response struct {
	Msg  string      `json:"msg" example:"ok" default:"ok"`
	Code int         `json:"code" example:"0" format:"int" default:"0"`
	Data interface{} `json:"data"`
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{Msg: "ok", Code: 0, Data: data}
}

func NewErrorResponse(msg string) *Response {
	return &Response{Msg: msg, Code: 4999, Data: nil}
}

// --- Tenant ---

type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TenantDomain struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	Domain    string    `json:"domain"`
	IsDefault bool      `json:"isDefault"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateTenantRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type AddDomainRequest struct {
	Domain    string `json:"domain" binding:"required"`
	IsDefault bool   `json:"isDefault"`
}

// --- Auth ---

type LoginAPIRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginAPIResponseData struct {
	Token string `json:"token,omitempty"`
}

type GetUserInfoAPIResponseData struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type ChangePasswordAPIRequest struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

// --- Config ---

type GetConfigAPIResponseData struct {
	Config interface{} `json:"config"`
}

type UpdateIdLengthRequest struct {
	IdLength        int `json:"idLength" binding:"required"`
	IdMinimumLength int `json:"idMinimumLength" binding:"required"`
	IdMaximumLength int `json:"idMaximumLength" binding:"required"`
}

type UpdateNotFoundConfigRequest struct {
	Mode  string `json:"mode" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// --- User ---

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	RawPassword string    `json:"-"`
	Password    []byte    `json:"password"`
	Salt        []byte    `json:"salt"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// --- Tenant Member ---

type TenantMember struct {
	TenantID string    `json:"tenantId"`
	UserID   string    `json:"userId"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joinedAt"`
}

func (m *TenantMember) IsAdmin() bool {
	return m.Role == RoleAdmin
}

// --- Tenant Invitation ---

type TenantInvitation struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	InviterID string    `json:"inviterId"`
	InviteeID string    `json:"inviteeId"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// --- Short Link ---

type ShortLink struct {
	Id          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	IsEnable    bool      `json:"is_enable"`
	CreatedBy   string    `json:"created_by"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

func NewShortLink(tenantID string, createdBy string, r *CreateShortLinkAPIRequest) *ShortLink {
	return &ShortLink{
		TenantID:    tenantID,
		Id:          r.Id,
		Url:         r.Url,
		Description: r.Description,
		IsEnable:    r.IsEnable,
		CreatedBy:   createdBy,
	}
}

type ShortLinkData struct {
	Id          string    `json:"id"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	IsEnable    bool      `json:"isEnable"`
	CreatedBy   string    `json:"createdBy"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

func ToShortLinkData(s *ShortLink) *ShortLinkData {
	return &ShortLinkData{
		Id: s.Id, Url: s.Url, Description: s.Description,
		IsEnable: s.IsEnable, CreatedBy: s.CreatedBy,
		CreateTime: s.CreateTime, UpdateTime: s.UpdateTime,
	}
}

func ToShortLinkDataSlice(s []*ShortLink) []*ShortLinkData {
	r := make([]*ShortLinkData, 0, len(s))
	for _, ss := range s {
		r = append(r, ToShortLinkData(ss))
	}
	return r
}

type CreateShortLinkAPIRequest struct {
	Id          string `json:"id"`
	Url         string `json:"url"`
	Description string `json:"description"`
	IsEnable    bool   `json:"isEnable"`
	IdLength    int    `json:"idLength"`
}

type GetShortLinkAPIResponseData struct {
	ShortLinkData *ShortLinkData `json:"shortLink"`
}

type CreateShortLinkAPIResponseData struct {
	ShortLinkData *ShortLinkData `json:"shortLink"`
}

type UpdateShortLinkAPIResponseData struct {
	ShortLinkData *ShortLinkData `json:"shortLink"`
}

type UpdateShortLinkAPIRequest struct {
	Url         string `json:"url" binding:"required"`
	Description string `json:"description"`
	IsEnable    bool   `json:"isEnable"`
}

type ListShortLinksAPIResponseData struct {
	ShortLinks []*ShortLinkData `json:"shortLinks"`
	Total      int64            `json:"total"`
}

// --- Request History ---

type RequestHistory struct {
	Id          string    `json:"id,omitempty"`
	ShortLinkID string    `json:"shortLinkId"`
	TenantID    string    `json:"tenantId,omitempty"`
	Url         string    `json:"url"`
	IP          string    `json:"ip"`
	UA          string    `json:"ua"`
	OS          string    `json:"os,omitempty"`
	Time        time.Time `json:"time"`
}

func (r *RequestHistory) MarshalBinary() (data []byte, err error) {
	return json.Marshal(r)
}

func NewRequestHistory(link *ShortLink, IP string, UA string, OS string) *RequestHistory {
	return &RequestHistory{
		ShortLinkID: link.Id,
		TenantID:    link.TenantID,
		Url:         link.Url,
		IP:          IP,
		UA:          UA,
		OS:          OS,
	}
}

type ShortLinkDataAPIResponseData struct {
	Histories []*RequestHistory `json:"histories"`
	Daily     []*DailyStats     `json:"daily"`
	OSDist    map[string]int    `json:"osDist"`
}

type DailyStats struct {
	Date string `json:"date"`
	PV   int    `json:"pv"`
	UV   int    `json:"uv"`
}

type UserPreference struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
