package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/utils"
	"log"
	"time"
)

type RequestHistoryListResult struct {
	Histories []*models.RequestHistory `json:"histories"`
	Total     int                      `json:"total"`
}

func NewEmptyRequestHistoryResult() *RequestHistoryListResult {
	return &RequestHistoryListResult{
		Histories: make([]*models.RequestHistory, 0),
		Total:     0,
	}
}

func (r *RequestHistoryListResult) AddHistory(h ...*models.RequestHistory) {
	r.Histories = append(r.Histories, h...)
	r.Total = len(r.Histories)
}

type requestHistoryRepository struct {
	db *redis.Client
}

var requestHistoryRepo = &requestHistoryRepository{db.GetRedisClient()}

func GetRequestHistoryRepo() *requestHistoryRepository {
	return requestHistoryRepo
}

func (r *requestHistoryRepository) Save(rh *models.RequestHistory) {
	rh.Time = time.Now()
	key := utils.GetRequestHistoryKey(rh.Link.Id, rh.Time)
	j, err := json.Marshal(rh)
	if err != nil {
		log.Printf("fail to save request history with key: %s, error: %v\n", key, err)
		return
	}

	r.db.LPush(key, j)
}

func (r *requestHistoryRepository) FindByDate(linkId string, d ...time.Time) (*RequestHistoryListResult, error) {
	var start time.Time
	var end time.Time
	dayDuration := time.Hour * 24
	result := NewEmptyRequestHistoryResult()

	if len(d) <= 1 {
		start = time.Now()
		end = start.Add(dayDuration)
	} else {
		start = d[0]
		end = d[len(d)-1]
	}
	if end.Before(start) {
		return result, fmt.Errorf("结束日期不能早于开始日期")
	}

	rawRs := make([]*redis.StringSliceCmd, 0)
	p := r.db.Pipeline()
	for ; start.Before(end); start = start.Add(dayDuration) {
		rawRs = append(rawRs, p.LRange(utils.GetRequestHistoryKey(linkId, start), 0, -1))
	}
	_, _ = p.Exec()

	for _, rs := range rawRs {
		for _, one := range rs.Val() {
			rh := &models.RequestHistory{}
			_ = json.Unmarshal([]byte(one), rh)
			result.AddHistory(rh)
		}
	}
	return result, nil
}

func (r *requestHistoryRepository) FindLatest(linkId string, size int64) (*RequestHistoryListResult, error) {
	key := utils.GetRequestHistoryKey(linkId, time.Now())
	rawRs, err := r.db.LRange(key, 0, size).Result()
	if err != nil {
		log.Printf("failed to find request history latest records with key: %s, err: %v\n", key, err)
	}

	result := NewEmptyRequestHistoryResult()
	for _, one := range rawRs {
		rh := &models.RequestHistory{}
		_ = json.Unmarshal([]byte(one), rh)
		result.AddHistory(rh)
	}
	return result, nil
}

type userRepository struct {
	db *redis.Client
}

var userRepo = &userRepository{db.GetRedisClient()}

func GetUserRepo() *userRepository {
	return userRepo
}

func (r *userRepository) IsExists(username string) bool {
	if username == "" {
		return false
	}

	exists, err := r.db.HExists(utils.GetUserKey(), username).Result()
	if err != nil {
		log.Printf("fail to check user exists with username: %s, error: %v\n", username, err)
		return false
	}
	return exists
}

func (r *userRepository) Save(u *models.User) error {
	if u.Username == "" || u.RawPassword == "" {
		return fmt.Errorf("username or password can not be empty string")
	}
	if _, exists := models.Roles[u.Role]; !exists {
		return fmt.Errorf("user role can not be %b", u.Role)
	}
	if r.IsExists(u.Username) {
		return fmt.Errorf("%s already exitis", u.Username)
	}

	salt, _ := utils.RandomSalt(32)
	dk, _ := utils.EncodePassword([]byte(u.RawPassword), salt)
	u.Password = dk
	u.Salt = salt
	u.CreateTime = time.Now()

	j, _ := json.Marshal(u)
	r.db.HSet(utils.GetUserKey(), u.Username, j)
	return nil
}

func (r *userRepository) UpdatePassword(u *models.User) error {
	if u.RawPassword == "" {
		return fmt.Errorf("password can not be empty string")
	}

	salt, _ := utils.RandomSalt(32)
	dk, _ := utils.EncodePassword([]byte(u.RawPassword), salt)
	u.Password = dk
	u.Salt = salt

	j, _ := json.Marshal(u)
	r.db.HSet(utils.GetUserKey(), u.Username, j)
	return nil
}

func (r *userRepository) FindOneByUsername(username string) (*models.User, error) {
	if username == "" {
		return nil, fmt.Errorf("username can not be empty string")
	}

	j, err := r.db.HGet(utils.GetUserKey(), username).Result()
	if err != nil {
		log.Printf("fail to get user with username: %s, error: %v\n", username, err)
		return nil, fmt.Errorf("用户不存在")
	}

	u := &models.User{}
	err = json.Unmarshal([]byte(j), u)
	if err != nil {
		log.Printf("fail to Unmarshal user with username: %s, error: %v\n", username, err)
		return nil, fmt.Errorf("用户不存在")
	}
	return u, nil
}

type shortLinkRepository struct {
	db *redis.Client
}

var shortLinkRepo = &shortLinkRepository{db.GetRedisClient()}

func GetShortLinkRepo() *shortLinkRepository {
	return shortLinkRepo
}

func (r *shortLinkRepository) generateId(l int) (string, error) {
	var id string
	for true {
		id = utils.RandStringRunes(l)
		rs, err := r.db.Exists(utils.GetShortLinkKey(id)).Result()
		if rs == 0 {
			break
		}
		if err != nil {
			log.Println(err)
			return "", err
		}
	}
	return id, nil
}

func (r *shortLinkRepository) save(s *models.ShortLink, isUpdate bool) error {
	if isUpdate && s.Id == "" {
		return fmt.Errorf("id错误")
	}
	if s.Url == "" {
		return fmt.Errorf("请填写url")
	}
	if s.CreatedBy == "" {
		return fmt.Errorf("未设置创建者，请通过接口创建短链接")
	}

	if !isUpdate {
		id, err := r.generateId(6)
		if err != nil {
			log.Println(err)
			return errors.New("服务器繁忙，请稍后再试")
		}
		if s.Id == "" {
			s.Id = id
		}
		s.Id = utils.TrimShortLinkId(s.Id)
		if s.Id == "" {
			return fmt.Errorf("id错误")
		}
		s.CreateTime = time.Now()
	}
	s.UpdateTime = time.Now()
	j, _ := json.Marshal(s)

	pipeline := r.db.Pipeline()
	// 保存短链接
	pipeline.Set(utils.GetShortLinkKey(s.Id), j, 0)
	// 保存用户的短链接记录，保存到创建者及全局
	record := redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: s.Id,
	}
	pipeline.ZAdd(utils.GetUserShortLinksKey(s.CreatedBy), record)
	pipeline.ZAdd(utils.GetShortLinksKey(), record)
	_, err := pipeline.Exec()
	if err != nil {
		log.Println(err)
		return errors.New("服务器繁忙，请稍后再试")
	}
	return nil
}

func (r *shortLinkRepository) Save(s *models.ShortLink) error {
	return r.save(s, false)
}

func (r *shortLinkRepository) Update(s *models.ShortLink, params *models.UpdateShortLinkParameter) error {
	s.Url = params.Url
	s.Description = params.Description
	s.IsEnable = params.IsEnable

	return r.save(s, true)
}

func (r *shortLinkRepository) Delete(s *models.ShortLink) {
	pipeline := r.db.Pipeline()

	// 删除短链接本身
	pipeline.Del(utils.GetShortLinkKey(s.Id))
	// 删除用户的短链接记录及全局短链接记录
	pipeline.ZRem(utils.GetUserShortLinksKey(s.CreatedBy), s.Id)
	pipeline.ZRem(utils.GetShortLinksKey(), s.Id)
	_, _ = pipeline.Exec()

	// 删除访问历史
	keys, _ := r.db.Keys(fmt.Sprintf("history:%s:*", s.Id)).Result()
	if len(keys) > 0 {
		r.db.Del(keys...)
	}
}

func (r *shortLinkRepository) Get(id string) (*models.ShortLink, error) {
	if id == "" {
		return nil, fmt.Errorf("短链接不存在")
	}

	key := utils.GetShortLinkKey(id)
	s := &models.ShortLink{}
	rs, err := r.db.Get(key).Result()
	if err != nil {
		log.Printf("fail to get short Link with Key: %s, error: %v\n", key, err)
		return s, fmt.Errorf("短链接不存在")
	}

	err = json.Unmarshal([]byte(rs), s)
	if err != nil {
		log.Printf("fail to unmarshal short Link, Key: %s, error: %v\n", key, err)
		return s, fmt.Errorf("短链接不存在")
	}

	return s, nil
}

type shortLinkListResult struct {
	ShortLinks []*models.ShortLink `json:"shortLinks"`
	Total      int64               `json:"total"`
}

func makeEmptyShortLinkListResult() *shortLinkListResult {
	return &shortLinkListResult{
		ShortLinks: make([]*models.ShortLink, 0),
		Total:      0,
	}
}

func (r *shortLinkListResult) AddLink(links ...*models.ShortLink) {
	r.ShortLinks = append(r.ShortLinks, links...)
}

func (r *shortLinkRepository) List(key string, start int64, stop int64) (*shortLinkListResult, error) {
	result := makeEmptyShortLinkListResult()

	total, _ := r.db.ZCard(key).Result()
	result.Total = total
	if total == 0 {
		return result, nil
	}

	ids, err := r.db.ZRevRange(key, start, stop).Result()
	if err != nil {
		return result, errors.New("系统繁忙请稍后再试")
	}

	if len(ids) == 0 {
		return result, nil
	}

	linkRs := make([]*redis.StringCmd, 0)
	pipeline := r.db.Pipeline()
	for _, id := range ids {
		r := pipeline.Get(utils.GetShortLinkKey(id))
		linkRs = append(linkRs, r)
	}
	_, _ = pipeline.Exec()

	for _, cmd := range linkRs {
		s := &models.ShortLink{}
		err = json.Unmarshal([]byte(cmd.Val()), s)
		result.AddLink(s)
	}
	return result, nil
}
