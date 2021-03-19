package models

import (
	"encoding/json"
	"github.com/jwma/jump-jump/internal/app/config"
	"time"
)

const RoleUser = 1
const RoleAdmin = 2

var Roles = map[int]string{
	RoleUser:  "user",
	RoleAdmin: "admin",
}

type Response struct {
	Msg  string      `json:"msg" example:"ok" default:"ok"`
	Code int         `json:"code" example:"0" format:"int" default:"0"`
	Data interface{} `json:"data"`
} // @name Response

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Msg:  "ok",
		Code: 0,
		Data: data,
	}
}

func NewErrorResponse(msg string) *Response {
	return &Response{
		Msg:  msg,
		Code: 4999,
		Data: nil,
	}
}

type LoginAPIRequest struct {
	// 用户名
	Username string `json:"username" binding:"required" example:"your_username"`

	// 密码
	Password string `json:"password" binding:"required" example:"your_password"`
} // @name LoginAPIRequest

type LoginAPIResponseData struct {
	// json web token
	Token string `json:"token,omitempty" example:"xxx.xxx.xxx"`
} // @name LoginAPIResponseData

type GetUserInfoAPIResponseData struct {
	// 用户名
	Username string `json:"username" example:"admin"`

	// 角色，1 普通用户 | 2 管理员
	Role int `json:"role" example:"1" enums:"1,2"`
} // @name GetUserInfoAPIResponseData

type ChangePasswordAPIRequest struct {
	// 原密码
	Password string `json:"password"`

	// 新密码
	NewPassword string `json:"newPassword"`
} // @name ChangePasswordAPIRequest

type GetConfigAPIResponseData struct {
	Config *config.SystemConfig `json:"config"`
} // @name GetConfigAPIResponseData

type UpdateLandingHostsAPIRequest struct {
	Hosts []string `json:"hosts" format:"array" example:"https://a.com/,https://b.com/"`
} // @name UpdateLandingHostsAPIRequest

type User struct {
	Username    string    `json:"username"`
	Role        int       `json:"role"`
	RawPassword string    `json:"-"`
	Password    []byte    `json:"password"`
	Salt        []byte    `json:"salt"`
	CreateTime  time.Time `json:"create_time"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

type ShortLink struct {
	Id          string    `json:"id" `
	Url         string    `json:"url"`
	Description string    `json:"description"`
	IsEnable    bool      `json:"is_enable"`
	CreatedBy   string    `json:"created_by"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

func NewShortLink(createdBy string, r *CreateShortLinkAPIRequest) *ShortLink {
	return &ShortLink{
		Id:          r.Id,
		Url:         r.Url,
		Description: r.Description,
		IsEnable:    r.IsEnable,
		CreatedBy:   createdBy,
	}
}

type ShortLinkData struct {
	Id string `json:"id" example:"RANDOM_ID" format:"string"`
	// 目标链接
	Url string `json:"url" example:"https://github.com/jwma/jump-jump" format:"string"`

	// 描述
	Description string `json:"description" example:"Jump Jump project" format:"string"`

	// 是否启用
	IsEnable bool `json:"isEnable" example:"true" format:"boolean"`

	// 创建者
	CreatedBy string `json:"createdBy" example:"admin" format:"string"`

	// 创建时间
	CreateTime time.Time `json:"createTime"`

	// 最后更新时间
	UpdateTime time.Time `json:"updateTime"`
} // @name ShortLinkData

func ToShortLinkData(s *ShortLink) *ShortLinkData {
	return &ShortLinkData{
		Id:          s.Id,
		Url:         s.Url,
		Description: s.Description,
		IsEnable:    s.IsEnable,
		CreatedBy:   s.CreatedBy,
		CreateTime:  s.CreateTime,
		UpdateTime:  s.UpdateTime,
	}
}

func ToShortLinkDataSlice(s []*ShortLink) []*ShortLinkData {
	r := make([]*ShortLinkData, 0)
	for _, ss := range s {
		r = append(r, ToShortLinkData(ss))
	}
	return r
}

type CreateShortLinkAPIRequest struct {
	// 只有管理员可以在创建的时候指定 ID
	Id string `json:"id" format:"string" example:"RANDOM_ID"`

	// 目标链接
	Url string `json:"url" example:"https://github.com/jwma/jump-jump"`

	// 描述
	Description string `json:"description" example:"Jump Jump project"`

	// 是否启用
	IsEnable bool `json:"isEnable" example:"true" format:"boolean"`

	// 短链接 ID 长度
	IdLength int `json:"idLength" example:"4" format:"int"`
} // @name CreateShortLinkAPIRequest

type GetShortLinkAPIResponseData struct {
	ShortLinkData *ShortLinkData `json:"shortLink"`
} // @name GetShortLinkAPIResponseData

type CreateShortLinkAPIResponseData struct {
	ShortLinkData *ShortLinkData `json:"shortLink"`
} // @name CreateShortLinkAPIResponseData

type UpdateShortLinkAPIResponseData struct {
	ShortLinkData *ShortLinkData `json:"shortLink"`
} // @name UpdateShortLinkAPIResponseData

type UpdateShortLinkAPIRequest struct {
	// 目标链接
	Url string `json:"url" binding:"required" example:"https://github.com/jwma/jump-jump"`

	// 描述
	Description string `json:"description" example:"Jump Jump project"`

	// 是否启用
	IsEnable bool `json:"isEnable" example:"true" format:"boolean"`
} // @name UpdateShortLinkAPIRequest

type ListShortLinksAPIResponseData struct {
	ShortLinks []*ShortLinkData `json:"shortLinks"`
	Total      int64            `json:"total" example:"10" format:"10"`
} // @name ListShortLinksAPIResponseData

type RequestHistory struct {
	Id   string     `json:"id"`
	Link *ShortLink `json:"-"`
	Url  string     `json:"url"` // 由于短链接的目标连接可能会被修改，可以在访问历史记录中记录一下当前的目标连接
	IP   string     `json:"ip"`
	UA   string     `json:"ua"`
	Time time.Time  `json:"time"`
} // @name RequestHistory

func (r *RequestHistory) MarshalBinary() (data []byte, err error) {
	return json.Marshal(r)
}

func NewRequestHistory(link *ShortLink, IP string, UA string) *RequestHistory {
	return &RequestHistory{Link: link, IP: IP, UA: UA, Url: link.Url}
}

type ShortLinkDataAPIResponseData struct {
	Histories []*RequestHistory `json:"histories"`
} // @name ShortLinkDataAPIResponseData

type ActiveLink struct {
	Id   string
	Time time.Time
}

type DailyReport struct {
	PV int            `json:"pv"`
	UV int            `json:"uv"`
	OS map[string]int `json:"os"`
}

func (d *DailyReport) MarshalBinary() (data []byte, err error) {
	return json.Marshal(d)
}

type DailyReportItem struct {
	Date   string       `json:"date"`
	Report *DailyReport `json:"report"`
}
