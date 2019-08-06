package redis

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/orzzzli/goutils/convert"
	"time"
)

var GlobalRedisPool *redis.Pool

// NewRedisPool初始化连接池
func NewRedisPool(url string, password string, idle int, idleTime int) {
	GlobalRedisPool = &redis.Pool{
		MaxIdle:     idle,
		IdleTimeout: time.Duration(idleTime) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(url)
			if err != nil {
				panic(err)
			}
			if password != "" {
				//验证redis密码
				if _, authErr := c.Do("AUTH", password); authErr != nil {
					panic(err)
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				panic(err)
			}
			return nil
		},
	}
}

//expire
func Expire(k string, ex int) error {
	if GlobalRedisPool == nil {
		return errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	var err error
	_, err = conn.Do("EXPIRE", k, ex)
	return err
}

//string
func Set(k string, v string, ex int) error {
	if GlobalRedisPool == nil {
		return errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	var err error
	if ex <= 0 {
		_, err = conn.Do("SET", k, v)
	} else {
		_, err = conn.Do("SET", k, v, "EX", ex)
	}
	return err
}
func Get(k string) (string, bool, error) {
	if GlobalRedisPool == nil {
		return "", false, errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	res, err := conn.Do("GET", k)
	resOp := ""
	find := false
	if res != nil {
		resOp = string(res.([]uint8))
		find = true
	}
	return resOp, find, err
}

//SortedSet
func ZAdd(key string, k string, v float32) error {
	if GlobalRedisPool == nil {
		return errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	var err error
	_, err = conn.Do("ZADD", key, v, k)
	return err
}
func ZRevRank(key string, k string) (int, bool, error) {
	if GlobalRedisPool == nil {
		return 0, false, errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	resultOp := 0
	find := false
	result, err := conn.Do("ZREVRANK", key, k)
	if result != nil {
		resultOp, _ = convert.Int64to32(result.(int64))
		find = true
	}
	return resultOp, find, err
}
