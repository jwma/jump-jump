package models

import (
	"encoding/json"
	"time"
)

const RoleUser = 1
const RoleAdmin = 2

var Roles = map[int]string{
	RoleUser:  "user",
	RoleAdmin: "admin",
}

type ChangePasswordParameter struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

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
	Id          string    `json:"id"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	IsEnable    bool      `json:"is_enable"`
	CreatedBy   string    `json:"created_by"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

type CreateShortLinkParameter struct {
	*ShortLink
	IdLength int `json:"id_length"`
}

type UpdateShortLinkParameter struct {
	Url         string `json:"url" binding:"required"`
	Description string `json:"description"`
	IsEnable    bool   `json:"is_enable"`
}

type RequestHistory struct {
	Id   string     `json:"id"`
	Link *ShortLink `json:"-"`
	Url  string     `json:"url"` // 由于短链接的目标连接可能会被修改，可以在访问历史记录中记录一下当前的目标连接
	IP   string     `json:"ip"`
	UA   string     `json:"ua"`
	Time time.Time  `json:"time"`
}

func (r *RequestHistory) MarshalBinary() (data []byte, err error) {
	return json.Marshal(r)
}

func NewRequestHistory(link *ShortLink, IP string, UA string) *RequestHistory {
	return &RequestHistory{Link: link, IP: IP, UA: UA, Url: link.Url}
}
