package models

// 短链接
type Link struct {
	Slug        string
	Url         string
	IsEnabled   bool
	Description string
	CreatedAt   int64
	UpdatedAt   int64
}

// 短链接请求记录结构
type RequestRecord struct {
	RemoteAddr string
	UserAgent  string
	RequestAt  int64
}
