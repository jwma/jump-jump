package repository

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
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
	r.Total += len(h)
}

type requestHistoryRepository struct {
	db *redis.Client
}

func NewRequestHistoryRepository(db *redis.Client) *requestHistoryRepository {
	return &requestHistoryRepository{db}
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

func NewUserRepository(db *redis.Client) *userRepository {
	return &userRepository{db}
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
