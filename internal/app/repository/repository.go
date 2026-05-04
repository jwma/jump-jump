package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/utils"
	"github.com/redis/go-redis/v9"
)

// --- Tenant Repository ---

type TenantRepository struct {
	db *pgxpool.Pool
}

var tenantRepo *TenantRepository

func GetTenantRepo(p *pgxpool.Pool) *TenantRepository {
	if tenantRepo == nil {
		tenantRepo = &TenantRepository{p}
	}
	return tenantRepo
}

func (r *TenantRepository) Create(req *models.CreateTenantRequest) (*models.Tenant, error) {
	t := &models.Tenant{}
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO tenants (name, slug) VALUES ($1, $2) RETURNING id, name, slug, is_active, created_at, updated_at`,
		req.Name, req.Slug).Scan(&t.ID, &t.Name, &t.Slug, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("创建租户失败: %w", err)
	}
	// Create default config
	r.db.Exec(context.Background(), `INSERT INTO tenant_configs (tenant_id) VALUES ($1)`, t.ID)
	return t, nil
}

func (r *TenantRepository) GetByID(id string) (*models.Tenant, error) {
	t := &models.Tenant{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, name, slug, is_active, created_at, updated_at FROM tenants WHERE id = $1`, id).
		Scan(&t.ID, &t.Name, &t.Slug, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("租户不存在")
	}
	return t, nil
}

func (r *TenantRepository) List() ([]*models.Tenant, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, name, slug, is_active, created_at, updated_at FROM tenants ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.Tenant, 0)
	for rows.Next() {
		t := &models.Tenant{}
		rows.Scan(&t.ID, &t.Name, &t.Slug, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
		result = append(result, t)
	}
	return result, nil
}

func (r *TenantRepository) AddDomain(tenantID, domain string, isDefault bool) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO tenant_domains (tenant_id, domain, is_default) VALUES ($1, $2, $3)`,
		tenantID, domain, isDefault)
	if err != nil {
		return fmt.Errorf("添加域名失败: %w", err)
	}
	return nil
}

func (r *TenantRepository) RemoveDomain(tenantID, domain string) error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM tenant_domains WHERE tenant_id = $1 AND domain = $2`, tenantID, domain)
	return err
}

func (r *TenantRepository) ListDomains(tenantID string) ([]*models.TenantDomain, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, tenant_id, domain, is_default, created_at FROM tenant_domains WHERE tenant_id = $1 ORDER BY created_at`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.TenantDomain, 0)
	for rows.Next() {
		d := &models.TenantDomain{}
		rows.Scan(&d.ID, &d.TenantID, &d.Domain, &d.IsDefault, &d.CreatedAt)
		result = append(result, d)
	}
	return result, nil
}

// --- User Repository (PG, tenant-scoped) ---

type userRepository struct {
	db *pgxpool.Pool
}

var userRepo *userRepository

func GetUserRepo(p *pgxpool.Pool) *userRepository {
	if userRepo == nil {
		userRepo = &userRepository{p}
	}
	return userRepo
}

func (r *userRepository) IsExists(tenantID, username string) bool {
	var exists bool
	r.db.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM users WHERE tenant_id = $1 AND username = $2)`,
		tenantID, username).Scan(&exists)
	return exists
}

func (r *userRepository) Save(u *models.User) error {
	if u.TenantID == "" || u.Username == "" || u.RawPassword == "" {
		return fmt.Errorf("tenant_id, username and password are required")
	}
	if _, exists := models.Roles[u.Role]; !exists {
		return fmt.Errorf("invalid user role: %d", u.Role)
	}
	if r.IsExists(u.TenantID, u.Username) {
		return fmt.Errorf("%s already exists", u.Username)
	}

	salt, err := utils.RandomSalt(32)
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}
	dk, err := utils.EncodePassword([]byte(u.RawPassword), salt)
	if err != nil {
		return fmt.Errorf("failed to encode password: %w", err)
	}
	u.Password = dk
	u.Salt = salt
	u.CreateTime = time.Now()

	_, err = r.db.Exec(context.Background(),
		`INSERT INTO users (tenant_id, username, password, salt, role, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $6)`,
		u.TenantID, u.Username, u.Password, u.Salt, u.Role, u.CreateTime)
	return err
}

func (r *userRepository) UpdatePassword(u *models.User) error {
	if u.RawPassword == "" {
		return fmt.Errorf("password can not be empty string")
	}

	salt, _ := utils.RandomSalt(32)
	dk, _ := utils.EncodePassword([]byte(u.RawPassword), salt)
	u.Password = dk
	u.Salt = salt

	_, err := r.db.Exec(context.Background(),
		`UPDATE users SET password = $1, salt = $2, updated_at = now() WHERE tenant_id = $3 AND username = $4`,
		dk, salt, u.TenantID, u.Username)
	return err
}

func (r *userRepository) FindOneByUsername(tenantID, username string) (*models.User, error) {
	if username == "" {
		return nil, fmt.Errorf("username can not be empty string")
	}

	u := &models.User{}
	err := r.db.QueryRow(context.Background(),
		`SELECT tenant_id, username, password, salt, role, created_at
		 FROM users WHERE tenant_id = $1 AND username = $2`,
		tenantID, username).Scan(&u.TenantID, &u.Username, &u.Password, &u.Salt, &u.Role, &u.CreateTime)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	return u, nil
}

// --- Short Link Repository (PG + Redis cache, tenant-scoped) ---

type shortLinkRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

var shortLinkRepo *shortLinkRepository

func GetShortLinkRepo(p *pgxpool.Pool, rdb *redis.Client) *shortLinkRepository {
	if shortLinkRepo == nil {
		shortLinkRepo = &shortLinkRepository{p, rdb}
	}
	return shortLinkRepo
}

func (r *shortLinkRepository) GenerateId(l int) (string, error) {
	for {
		id := utils.RandStringRunes(l)
		var exists bool
		r.db.QueryRow(context.Background(),
			`SELECT EXISTS(SELECT 1 FROM short_links WHERE id = $1)`, id).Scan(&exists)
		if !exists {
			return id, nil
		}
	}
}

func (r *shortLinkRepository) Save(s *models.ShortLink) error {
	if s.Id == "" {
		return fmt.Errorf("id不能为空")
	}
	if s.Url == "" {
		return fmt.Errorf("请填写url")
	}
	if s.CreatedBy == "" {
		return fmt.Errorf("未设置创建者")
	}

	s.CreateTime = time.Now()
	s.UpdateTime = time.Now()

	_, err := r.db.Exec(context.Background(),
		`INSERT INTO short_links (id, tenant_id, url, description, is_enabled, created_by, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $7)`,
		s.Id, s.TenantID, s.Url, s.Description, s.IsEnable, s.CreatedBy, s.CreateTime)
	if err != nil {
		log.Printf("fail to save short link: %v", err)
		return errors.New("服务器繁忙，请稍后再试")
	}
	return nil
}

func (r *shortLinkRepository) Update(s *models.ShortLink, params *models.UpdateShortLinkAPIRequest) error {
	s.Url = params.Url
	s.Description = params.Description
	s.IsEnable = params.IsEnable
	s.UpdateTime = time.Now()

	_, err := r.db.Exec(context.Background(),
		`UPDATE short_links SET url = $1, description = $2, is_enabled = $3, updated_at = $4 WHERE id = $5`,
		s.Url, s.Description, s.IsEnable, s.UpdateTime, s.Id)
	if err != nil {
		return errors.New("服务器繁忙，请稍后再试")
	}

	r.rdb.Del(context.Background(), utils.GetShortLinkCacheKey(s.Id))
	return nil
}

func (r *shortLinkRepository) Delete(s *models.ShortLink) {
	r.db.Exec(context.Background(), `DELETE FROM short_links WHERE id = $1`, s.Id)
	r.rdb.Del(context.Background(), utils.GetShortLinkCacheKey(s.Id))
	r.rdb.Del(context.Background(), utils.GetRequestHistoryKey(s.Id))
	r.rdb.Del(context.Background(), utils.GetDailyReportKey(s.Id))
}

func (r *shortLinkRepository) Get(id string) (*models.ShortLink, error) {
	if id == "" {
		return nil, fmt.Errorf("短链接不存在")
	}

	// Try cache
	cacheKey := utils.GetShortLinkCacheKey(id)
	val, err := r.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		s := &models.ShortLink{}
		if json.Unmarshal([]byte(val), s) == nil {
			return s, nil
		}
	}

	// Cache miss → PG
	s := &models.ShortLink{}
	err = r.db.QueryRow(context.Background(),
		`SELECT id, tenant_id, url, description, is_enabled, created_by, created_at, updated_at
		 FROM short_links WHERE id = $1`, id).Scan(
		&s.Id, &s.TenantID, &s.Url, &s.Description, &s.IsEnable, &s.CreatedBy, &s.CreateTime, &s.UpdateTime)
	if err != nil {
		return nil, fmt.Errorf("短链接不存在")
	}

	j, _ := json.Marshal(s)
	r.rdb.Set(context.Background(), cacheKey, j, 5*time.Minute)
	return s, nil
}

type shortLinkListResult struct {
	ShortLinks []*models.ShortLink `json:"shortLinks"`
	Total      int64               `json:"total"`
}

func makeEmptyShortLinkListResult() *shortLinkListResult {
	return &shortLinkListResult{ShortLinks: make([]*models.ShortLink, 0), Total: 0}
}

func (r *shortLinkRepository) List(tenantID, username string, isAdmin bool, start, pageSize int64) (*shortLinkListResult, error) {
	result := makeEmptyShortLinkListResult()

	var countQuery, dataQuery string
	var args []interface{}

	args = append(args, tenantID)

	if isAdmin {
		countQuery = `SELECT COUNT(*) FROM short_links WHERE tenant_id = $1`
		dataQuery = `SELECT id, tenant_id, url, description, is_enabled, created_by, created_at, updated_at
					 FROM short_links WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	} else {
		countQuery = `SELECT COUNT(*) FROM short_links WHERE tenant_id = $1 AND created_by = $2`
		dataQuery = `SELECT id, tenant_id, url, description, is_enabled, created_by, created_at, updated_at
					 FROM short_links WHERE tenant_id = $1 AND created_by = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`
		args = append(args, username)
	}

	err := r.db.QueryRow(context.Background(), countQuery, args...).Scan(&result.Total)
	if err != nil || result.Total == 0 {
		return result, nil
	}

	dataArgs := append(args, pageSize, start)
	rows, err := r.db.Query(context.Background(), dataQuery, dataArgs...)
	if err != nil {
		return result, errors.New("系统繁忙请稍后再试")
	}
	defer rows.Close()

	for rows.Next() {
		s := &models.ShortLink{}
		rows.Scan(&s.Id, &s.TenantID, &s.Url, &s.Description, &s.IsEnable, &s.CreatedBy, &s.CreateTime, &s.UpdateTime)
		result.ShortLinks = append(result.ShortLinks, s)
	}
	return result, nil
}

// --- Request History Repository (Redis, unchanged) ---

type requestHistoryListResult struct {
	Histories []*models.RequestHistory `json:"histories"`
	Total     int                      `json:"total"`
}

func newEmptyRequestHistoryResult() *requestHistoryListResult {
	return &requestHistoryListResult{Histories: make([]*models.RequestHistory, 0), Total: 0}
}

func (r *requestHistoryListResult) addHistory(h ...*models.RequestHistory) {
	r.Histories = append(r.Histories, h...)
	r.Total = len(r.Histories)
}

type requestHistoryRepository struct {
	db *redis.Client
}

var requestHistoryRepo *requestHistoryRepository

func GetRequestHistoryRepo(rdb *redis.Client) *requestHistoryRepository {
	if requestHistoryRepo == nil {
		requestHistoryRepo = &requestHistoryRepository{rdb}
	}
	return requestHistoryRepo
}

func (r *requestHistoryRepository) Save(rh *models.RequestHistory) {
	rh.Id = utils.RandStringRunes(6)
	rh.Time = time.Now()
	key := utils.GetRequestHistoryKey(rh.Link.Id)

	_, err := r.db.ZAdd(context.Background(), key, redis.Z{
		Score:  float64(rh.Time.Unix()),
		Member: rh,
	}).Result()
	if err != nil {
		log.Printf("fail to save request history: %v", err)
	}
}

func (r *requestHistoryRepository) FindLatest(linkId string, size int64) (*requestHistoryListResult, error) {
	key := utils.GetRequestHistoryKey(linkId)
	rs, err := r.db.ZRangeWithScores(context.Background(), key, -size, -1).Result()
	if err != nil {
		log.Printf("failed to find request history: %v", err)
	}

	utils.ReverseAny(rs)
	result := newEmptyRequestHistoryResult()
	for _, one := range rs {
		rh := &models.RequestHistory{}
		_ = json.Unmarshal([]byte(one.Member.(string)), rh)
		result.addHistory(rh)
	}
	return result, nil
}

func (r *requestHistoryRepository) FindByDateRange(linkId string, startTime, endTime time.Time) []*models.RequestHistory {
	rs, _ := r.db.ZRangeByScoreWithScores(context.Background(), utils.GetRequestHistoryKey(linkId), &redis.ZRangeBy{
		Min: strconv.Itoa(int(startTime.Unix())),
		Max: strconv.Itoa(int(endTime.Unix())),
	}).Result()

	rhs := make([]*models.RequestHistory, 0)
	for _, one := range rs {
		rh := &models.RequestHistory{}
		_ = json.Unmarshal([]byte(one.Member.(string)), rh)
		rhs = append(rhs, rh)
	}
	return rhs
}

// --- Active Link Repository (Redis) ---

type activeLinkRepository struct {
	db *redis.Client
}

var activeLinkRepo *activeLinkRepository

func GetActiveLinkRepo(rdb *redis.Client) *activeLinkRepository {
	if activeLinkRepo == nil {
		activeLinkRepo = &activeLinkRepository{rdb}
	}
	return activeLinkRepo
}

func (r *activeLinkRepository) Save(linkId string) {
	r.db.ZAdd(context.Background(), utils.GetActiveLinkKey(), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: linkId,
	})
}

func (r *activeLinkRepository) FindByDateRange(startTime, endTime time.Time) []*models.ActiveLink {
	result := make([]*models.ActiveLink, 0)
	rs, _ := r.db.ZRangeByScoreWithScores(context.Background(), utils.GetActiveLinkKey(), &redis.ZRangeBy{
		Min: strconv.Itoa(int(startTime.Unix())),
		Max: strconv.Itoa(int(endTime.Unix())),
	}).Result()
	for _, one := range rs {
		result = append(result, &models.ActiveLink{Id: one.Member.(string), Time: time.Unix(int64(one.Score), 0)})
	}
	return result
}

// --- Daily Report Repository (Redis) ---

type dailyReportRepository struct {
	db *redis.Client
}

var dailyReportRepo *dailyReportRepository

func GetDailyReportRepo(rdb *redis.Client) *dailyReportRepository {
	if dailyReportRepo == nil {
		dailyReportRepo = &dailyReportRepository{rdb}
	}
	return dailyReportRepo
}

func (r *dailyReportRepository) Save(linkId string, reportKey string, report *models.DailyReport) {
	r.db.HSet(context.Background(), utils.GetDailyReportKey(linkId), reportKey, report)
}

func (r *dailyReportRepository) FindRecent(linkId string, days int) []*models.DailyReportItem {
	if days < 1 {
		days = 1
	}
	now := time.Now()
	d := now.AddDate(0, 0, -days+1)
	reportKeys := make([]string, 0)
	for d.Before(now) {
		reportKeys = append(reportKeys, d.Format("2006-01-02"))
		d = d.AddDate(0, 0, 1)
	}
	reportKeys = append(reportKeys, now.Format("2006-01-02"))
	reports := make([]*models.DailyReportItem, days)
	rs, _ := r.db.HMGet(context.Background(), utils.GetDailyReportKey(linkId), reportKeys...).Result()
	for i := 0; i < days; i++ {
		r := &models.DailyReport{}
		if rs[i] != nil {
			json.Unmarshal([]byte(rs[i].(string)), r)
		}
		reports[i] = &models.DailyReportItem{Date: reportKeys[i], Report: r}
	}
	return reports
}
