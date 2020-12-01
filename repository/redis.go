package repository

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	errs "github.com/pkg/errors"

	"github.com/ercross/errcross"
)

type redisRepo struct {
	client *redis.Client
}

func newRedisClient(redisUrl string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	return client, err
}

func NewRedisRepo(redisUrl string) (errcross.ErrcrossRepository, error) {
	redisRepo := &redisRepo{}
	client, err := newRedisClient(redisUrl)
	if err != nil {
		return nil, errs.Wrap(err, "repository.NewredisRepo")
	}
	redisRepo.client = client
	return redisRepo, nil
}

//generateRedisKey creates the key to retreive the data from Redis database
func (r *redisRepo) generateRedisKey(urlKey string) string {
	return fmt.Sprintf("redirect:%s", urlKey)
}

func (r *redisRepo) Find(key string) (*errcross.Errcross, error) {
	e := errcross.Errcross{}
	redisKey := r.generateRedisKey(key)
	data, err := r.client.HGetAll(redisKey).Result()
	if err != nil {
		return nil, errs.Wrap(err, "repository.redis.Find")
	}
	if len(data) == 0 {
		return nil, errs.Wrap(errcross.ErrKeyNotFound, "repository.redis.Find")
	}

	//extract timestamp, as int64 bit, from redis through key = timestamp
	createdAt, err := strconv.ParseInt(data["timestamp"], 10, 64)
	if err != nil {
		return nil, errs.Wrap(err, "repository.redis.Find")
	}
	e.Key = data["key"]
	e.URL = data["url"]
	e.Timestamp = createdAt
	return &e, nil
}

func (r *redisRepo) Store(e *errcross.Errcross) error {
	redisKey := r.generateRedisKey(e.Key)
	data := map[string]interface{}{
		"key":       e.Key,
		"url":       e.URL,
		"timestamp": e.Timestamp,
	}

	//save key and data into redis
	_, err := r.client.HMSet(redisKey, data).Result()
	if err != nil {
		return errs.Wrap(err, "repository.redis.Store")
	}
	return nil
}
