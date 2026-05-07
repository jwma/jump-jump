package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

// --- User Repository (PG, globally unique username) ---

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

func (r *userRepository) IsExists(username string) bool {
	var exists bool
	r.db.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`,
		username).Scan(&exists)
	return exists
}

func (r *userRepository) Save(u *models.User) error {
	if u.Username == "" || u.RawPassword == "" {
		return fmt.Errorf("username and password are required")
	}
	if r.IsExists(u.Username) {
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
	u.IsActive = true
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt

	return r.db.QueryRow(context.Background(),
		`INSERT INTO users (username, password, salt, is_active, is_super, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $6)
		 RETURNING id`,
		u.Username, u.Password, u.Salt, u.IsActive, u.IsSuper, u.CreatedAt).Scan(&u.ID)
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
		`UPDATE users SET password = $1, salt = $2, updated_at = now() WHERE id = $3`,
		dk, salt, u.ID)
	return err
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	if username == "" {
		return nil, fmt.Errorf("username can not be empty string")
	}

	u := &models.User{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, username, password, salt, is_active, is_super, created_at, updated_at
		 FROM users WHERE username = $1`,
		username).Scan(&u.ID, &u.Username, &u.Password, &u.Salt, &u.IsActive, &u.IsSuper, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	return u, nil
}

func (r *userRepository) FindByID(userID string) (*models.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id can not be empty string")
	}

	u := &models.User{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, username, password, salt, is_active, is_super, created_at, updated_at
		 FROM users WHERE id = $1`,
		userID).Scan(&u.ID, &u.Username, &u.Password, &u.Salt, &u.IsActive, &u.IsSuper, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	return u, nil
}

type ListOptions struct {
	Query    string
	Page     int64
	PageSize int64
}

type UserListResult struct {
	Users []*models.User `json:"users"`
	Total int64          `json:"total"`
}

func (r *userRepository) List(opts ListOptions) (*UserListResult, error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 || opts.PageSize > 100 {
		opts.PageSize = 20
	}
	offset := (opts.Page - 1) * opts.PageSize

	result := &UserListResult{Users: make([]*models.User, 0), Total: 0}

	var countQuery, dataQuery string
	var args []interface{}

	if opts.Query != "" {
		countQuery = `SELECT COUNT(*) FROM users WHERE username ILIKE $1`
		dataQuery = `SELECT id, username, is_active, is_super, created_at, updated_at
					 FROM users WHERE username ILIKE $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = append(args, "%"+opts.Query+"%")
	} else {
		countQuery = `SELECT COUNT(*) FROM users`
		dataQuery = `SELECT id, username, is_active, is_super, created_at, updated_at
					 FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	}

	err := r.db.QueryRow(context.Background(), countQuery, args...).Scan(&result.Total)
	if err != nil {
		return nil, fmt.Errorf("查询用户总数失败: %w", err)
	}
	if result.Total == 0 {
		return result, nil
	}

	dataArgs := append(args, opts.PageSize, offset)
	rows, err := r.db.Query(context.Background(), dataQuery, dataArgs...)
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.IsActive, &u.IsSuper, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("查询用户列表失败: %w", err)
		}
		result.Users = append(result.Users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}
	return result, nil
}

func (r *userRepository) UpdatePasswordByID(userID, newPassword string) error {
	salt, err := utils.RandomSalt(32)
	if err != nil {
		return fmt.Errorf("生成盐失败: %w", err)
	}
	dk, err := utils.EncodePassword([]byte(newPassword), salt)
	if err != nil {
		return fmt.Errorf("编码密码失败: %w", err)
	}
	ct, err := r.db.Exec(context.Background(),
		`UPDATE users SET password = $1, salt = $2, updated_at = now() WHERE id = $3`,
		dk, salt, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("用户不存在")
	}
	return nil
}

func (r *userRepository) UpdateStatus(userID string, isActive bool) error {
	ct, err := r.db.Exec(context.Background(),
		`UPDATE users SET is_active = $1, updated_at = now() WHERE id = $2`,
		isActive, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("用户不存在")
	}
	return nil
}

func (r *userRepository) GetUserTenants(u *models.User) ([]*models.UserTenantEntry, error) {
	if u.IsSuper {
		rows, err := r.db.Query(context.Background(),
			`SELECT id, name FROM tenants ORDER BY created_at`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		result := make([]*models.UserTenantEntry, 0)
		for rows.Next() {
			e := &models.UserTenantEntry{Role: models.RoleAdmin}
			rows.Scan(&e.TenantID, &e.TenantName)
			result = append(result, e)
		}
		return result, nil
	}

	memberRepo := GetTenantMemberRepo(r.db)
	members, err := memberRepo.ListByUser(u.ID)
	if err != nil {
		return nil, err
	}

	tenantRepo := GetTenantRepo(r.db)
	result := make([]*models.UserTenantEntry, 0, len(members))
	for _, m := range members {
		t, err := tenantRepo.GetByID(m.TenantID)
		if err != nil {
			continue
		}
		result = append(result, &models.UserTenantEntry{
			TenantID:   m.TenantID,
			TenantName: t.Name,
			Role:       m.Role,
		})
	}
	return result, nil
}

// --- Tenant Member Repository ---

type TenantMemberRepository struct {
	db *pgxpool.Pool
}

var tenantMemberRepo *TenantMemberRepository

func GetTenantMemberRepo(p *pgxpool.Pool) *TenantMemberRepository {
	if tenantMemberRepo == nil {
		tenantMemberRepo = &TenantMemberRepository{p}
	}
	return tenantMemberRepo
}

func (r *TenantMemberRepository) Save(m *models.TenantMember) error {
	if m.TenantID == "" || m.UserID == "" {
		return fmt.Errorf("tenant_id and user_id are required")
	}
	if m.Role == "" {
		m.Role = models.RoleMember
	}
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO tenant_members (tenant_id, user_id, role, joined_at)
		 VALUES ($1, $2, $3, now())
		 ON CONFLICT (tenant_id, user_id) DO NOTHING`,
		m.TenantID, m.UserID, m.Role)
	return err
}

func (r *TenantMemberRepository) Get(tenantID, userID string) (*models.TenantMember, error) {
	m := &models.TenantMember{}
	err := r.db.QueryRow(context.Background(),
		`SELECT tenant_id, user_id, role, joined_at
		 FROM tenant_members WHERE tenant_id = $1 AND user_id = $2`,
		tenantID, userID).Scan(&m.TenantID, &m.UserID, &m.Role, &m.JoinedAt)
	if err != nil {
		return nil, fmt.Errorf("成员关系不存在")
	}
	return m, nil
}

func (r *TenantMemberRepository) Delete(tenantID, userID string) error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM tenant_members WHERE tenant_id = $1 AND user_id = $2`,
		tenantID, userID)
	return err
}

func (r *TenantMemberRepository) ListByTenant(tenantID string) ([]*models.TenantMember, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT tenant_id, user_id, role, joined_at
		 FROM tenant_members WHERE tenant_id = $1 ORDER BY joined_at`,
		tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.TenantMember, 0)
	for rows.Next() {
		m := &models.TenantMember{}
		rows.Scan(&m.TenantID, &m.UserID, &m.Role, &m.JoinedAt)
		result = append(result, m)
	}
	return result, nil
}

func (r *TenantMemberRepository) UpdateRole(tenantID, userID, role string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE tenant_members SET role = $1 WHERE tenant_id = $2 AND user_id = $3`,
		role, tenantID, userID)
	return err
}

func (r *TenantMemberRepository) IsMember(tenantID, userID string) bool {
	var exists bool
	r.db.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM tenant_members WHERE tenant_id = $1 AND user_id = $2)`,
		tenantID, userID).Scan(&exists)
	return exists
}

func (r *TenantMemberRepository) ListByUser(userID string) ([]*models.TenantMember, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT tenant_id, user_id, role, joined_at
		 FROM tenant_members WHERE user_id = $1 ORDER BY joined_at`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.TenantMember, 0)
	for rows.Next() {
		m := &models.TenantMember{}
		rows.Scan(&m.TenantID, &m.UserID, &m.Role, &m.JoinedAt)
		result = append(result, m)
	}
	return result, nil
}

// --- Tenant Invitation Repository ---

type TenantInvitationRepository struct {
	db *pgxpool.Pool
}

var tenantInvitationRepo *TenantInvitationRepository

func GetTenantInvitationRepo(p *pgxpool.Pool) *TenantInvitationRepository {
	if tenantInvitationRepo == nil {
		tenantInvitationRepo = &TenantInvitationRepository{p}
	}
	return tenantInvitationRepo
}

func (r *TenantInvitationRepository) Create(inv *models.TenantInvitation) error {
	if inv.TenantID == "" || inv.InviterID == "" || inv.InviteeID == "" {
		return fmt.Errorf("tenant_id, inviter_id and invitee_id are required")
	}
	inv.Status = models.InvitationStatusPending
	return r.db.QueryRow(context.Background(),
		`INSERT INTO tenant_invitations (tenant_id, inviter_id, invitee_id, status, created_at)
		 VALUES ($1, $2, $3, $4, now())
		 RETURNING id, created_at`,
		inv.TenantID, inv.InviterID, inv.InviteeID, inv.Status).Scan(&inv.ID, &inv.CreatedAt)
}

func (r *TenantInvitationRepository) Get(id string) (*models.TenantInvitation, error) {
	inv := &models.TenantInvitation{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, tenant_id, inviter_id, invitee_id, status, created_at
		 FROM tenant_invitations WHERE id = $1`, id).
		Scan(&inv.ID, &inv.TenantID, &inv.InviterID, &inv.InviteeID, &inv.Status, &inv.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("邀请不存在")
	}
	return inv, nil
}

func (r *TenantInvitationRepository) UpdateStatus(id, status string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE tenant_invitations SET status = $1 WHERE id = $2`, status, id)
	return err
}

func (r *TenantInvitationRepository) ListByTenant(tenantID string) ([]*models.TenantInvitation, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, tenant_id, inviter_id, invitee_id, status, created_at
		 FROM tenant_invitations WHERE tenant_id = $1 ORDER BY created_at DESC`,
		tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.TenantInvitation, 0)
	for rows.Next() {
		inv := &models.TenantInvitation{}
		rows.Scan(&inv.ID, &inv.TenantID, &inv.InviterID, &inv.InviteeID, &inv.Status, &inv.CreatedAt)
		result = append(result, inv)
	}
	return result, nil
}

func (r *TenantInvitationRepository) ListByInvitee(inviteeID string) ([]*models.TenantInvitation, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, tenant_id, inviter_id, invitee_id, status, created_at
		 FROM tenant_invitations WHERE invitee_id = $1 ORDER BY created_at DESC`,
		inviteeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.TenantInvitation, 0)
	for rows.Next() {
		inv := &models.TenantInvitation{}
		rows.Scan(&inv.ID, &inv.TenantID, &inv.InviterID, &inv.InviteeID, &inv.Status, &inv.CreatedAt)
		result = append(result, inv)
	}
	return result, nil
}


func (r *TenantInvitationRepository) FindPendingByInvitee(inviteeID string) ([]*models.TenantInvitation, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, tenant_id, inviter_id, invitee_id, status, created_at
		 FROM tenant_invitations WHERE invitee_id = $1 AND status = $2 ORDER BY created_at DESC`,
		inviteeID, models.InvitationStatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.TenantInvitation, 0)
	for rows.Next() {
		inv := &models.TenantInvitation{}
		rows.Scan(&inv.ID, &inv.TenantID, &inv.InviterID, &inv.InviteeID, &inv.Status, &inv.CreatedAt)
		result = append(result, inv)
	}
	return result, nil
}

func (r *TenantInvitationRepository) HasPendingInvitation(tenantID, inviteeID string) bool {
	var exists bool
	r.db.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM tenant_invitations WHERE tenant_id = $1 AND invitee_id = $2 AND status = $3)`,
		tenantID, inviteeID, models.InvitationStatusPending).Scan(&exists)
	return exists
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

// --- Request History Repository (Write-behind: Redis buffer → PG flush) ---

type requestHistoryRepository struct {
	rdb *redis.Client
	db  *pgxpool.Pool
}

var requestHistoryRepo *requestHistoryRepository

func GetRequestHistoryRepo(rdb *redis.Client, db *pgxpool.Pool) *requestHistoryRepository {
	if requestHistoryRepo == nil {
		requestHistoryRepo = &requestHistoryRepository{rdb, db}
	}
	return requestHistoryRepo
}

func (r *requestHistoryRepository) Save(rh *models.RequestHistory) {
	rh.Time = time.Now()

	// Try Redis buffer first
	if r.rdb != nil {
		_, err := r.rdb.ZAdd(context.Background(), utils.RequestHistoryBufferKey, redis.Z{
			Score:  float64(rh.Time.Unix()),
			Member: rh,
		}).Result()
		if err == nil {
			return
		}
		log.Printf("request history: redis buffer write failed, falling back to PG: %v", err)
	}

	// Fallback: direct PG write
	r.saveDirect(rh)
}

func (r *requestHistoryRepository) saveDirect(rh *models.RequestHistory) {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO request_histories (short_link_id, tenant_id, url, ip, ua, os, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		rh.ShortLinkID, rh.TenantID, rh.Url, rh.IP, rh.UA, rh.OS, rh.Time)
	if err != nil {
		log.Printf("request history: direct PG write failed: %v", err)
	}
}

func (r *requestHistoryRepository) FindByDateRange(linkId string, startTime, endTime time.Time) []*models.RequestHistory {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, short_link_id, url, ip, ua, os, created_at
		 FROM request_histories
		 WHERE short_link_id = $1 AND created_at BETWEEN $2 AND $3
		 ORDER BY created_at DESC`,
		linkId, startTime, endTime)
	if err != nil {
		log.Printf("request history query failed: %v", err)
		return make([]*models.RequestHistory, 0)
	}
	defer rows.Close()

	rhs := make([]*models.RequestHistory, 0)
	for rows.Next() {
		rh := &models.RequestHistory{}
		rows.Scan(&rh.Id, &rh.ShortLinkID, &rh.Url, &rh.IP, &rh.UA, &rh.OS, &rh.Time)
		rhs = append(rhs, rh)
	}
	return rhs
}

func (r *requestHistoryRepository) GetAggregatedStats(linkId string, startTime, endTime time.Time) ([]*models.DailyStats, map[string]int) {
	// Daily PV/UV
	rows, err := r.db.Query(context.Background(),
		`SELECT DATE(created_at) AS day, COUNT(*) AS pv, COUNT(DISTINCT ip) AS uv
		 FROM request_histories
		 WHERE short_link_id = $1 AND created_at BETWEEN $2 AND $3
		 GROUP BY day ORDER BY day`,
		linkId, startTime, endTime)
	if err != nil {
		log.Printf("request history aggregation failed: %v", err)
		return nil, nil
	}
	defer rows.Close()

	daily := make([]*models.DailyStats, 0)
	for rows.Next() {
		ds := &models.DailyStats{}
		rows.Scan(&ds.Date, &ds.PV, &ds.UV)
		daily = append(daily, ds)
	}

	// OS distribution
	osRows, err := r.db.Query(context.Background(),
		`SELECT os, COUNT(*) FROM request_histories
		 WHERE short_link_id = $1 AND created_at BETWEEN $2 AND $3 AND os != ''
		 GROUP BY os ORDER BY count DESC`, linkId, startTime, endTime)
	if err != nil {
		return daily, nil
	}
	defer osRows.Close()

	osDist := make(map[string]int)
	for osRows.Next() {
		var osName string
		var count int
		osRows.Scan(&osName, &count)
		osDist[osName] = count
	}

	return daily, osDist
}

// --- User Preference Repository (PG) ---

type userPreferenceRepository struct {
	db *pgxpool.Pool
}

var userPrefRepo *userPreferenceRepository

func GetUserPreferenceRepo(p *pgxpool.Pool) *userPreferenceRepository {
	if userPrefRepo == nil {
		userPrefRepo = &userPreferenceRepository{p}
	}
	return userPrefRepo
}

func (r *userPreferenceRepository) GetAll(userID string) ([]*models.UserPreference, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT key, value FROM user_preferences WHERE user_id = $1 ORDER BY key`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prefs := make([]*models.UserPreference, 0)
	for rows.Next() {
		p := &models.UserPreference{}
		rows.Scan(&p.Key, &p.Value)
		prefs = append(prefs, p)
	}
	return prefs, nil
}

func (r *userPreferenceRepository) Upsert(userID string, prefs []*models.UserPreference) error {
	for _, p := range prefs {
		_, err := r.db.Exec(context.Background(),
			`INSERT INTO user_preferences (user_id, key, value) VALUES ($1, $2, $3)
			 ON CONFLICT (user_id, key) DO UPDATE SET value = $3`,
			userID, p.Key, p.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
