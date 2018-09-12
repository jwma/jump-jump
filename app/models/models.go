package models

// 短链接
type Link struct {
	Slug        string
	Url         string
	IsEnabled   bool
	Description string
	CreatedBy	string
	CreatedAt   int64
	UpdatedAt   int64
}

// 短链接请求记录结构
type RequestRecord struct {
	RemoteAddr string
	UserAgent  string
	RequestAt  int64
}

type User struct {
	Username  string
	Password  string
	Salt      string
	CreatedAt int64
}
