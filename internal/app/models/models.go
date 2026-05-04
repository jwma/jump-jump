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

type LoginAPIRequest struct {
	Username string `json:"username" binding:"required" example:"your_username"`
	Password string `json:"password" binding:"required" example:"your_password"`
}

type LoginAPIResponseData struct {
	Token string `json:"token,omitempty" example:"xxx.xxx.xxx"`
}

type GetUserInfoAPIResponseData struct {
	Username string `json:"username" example:"admin"`
	Role     int    `json:"role" example:"1" enums:"1,2"`
}

type ChangePasswordAPIRequest struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

type GetConfigAPIResponseData struct {
	Config interface{} `json:"config"`
}

type UpdateLandingHostsAPIRequest struct {
	Hosts []string `json:"hosts" format:"array" example:"https://a.com/,https://b.com/"`
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
	Id          string    `json:"id" example:"RANDOM_ID" format:"string"`
	Url         string    `json:"url" example:"https://github.com/jwma/jump-jump" format:"string"`
	Description string    `json:"description" example:"Jump Jump project" format:"string"`
	IsEnable    bool      `json:"isEnable" example:"true" format:"boolean"`
	CreatedBy   string    `json:"createdBy" example:"admin" format:"string"`
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
	Id          string `json:"id" format:"string" example:"RANDOM_ID"`
	Url         string `json:"url" example:"https://github.com/jwma/jump-jump"`
	Description string `json:"description" example:"Jump Jump project"`
	IsEnable    bool   `json:"isEnable" example:"true" format:"boolean"`
	IdLength    int    `json:"idLength" example:"4" format:"int"`
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
	Url         string `json:"url" binding:"required" example:"https://github.com/jwma/jump-jump"`
	Description string `json:"description" example:"Jump Jump project"`
	IsEnable    bool   `json:"isEnable" example:"true" format:"boolean"`
}

type ListShortLinksAPIResponseData struct {
	ShortLinks []*ShortLinkData `json:"shortLinks"`
	Total      int64            `json:"total" example:"10" format:"10"`
}

type RequestHistory struct {
	Id   string     `json:"id"`
	Link *ShortLink `json:"-"`
	Url  string     `json:"url"`
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

type ShortLinkDataAPIResponseData struct {
	Histories []*RequestHistory `json:"histories"`
}

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
