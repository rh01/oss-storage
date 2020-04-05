// Package cache : redis连接池
package cache

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	poll      *redis.Pool
	redisHost = "127.0.0.1:6379"
	redisPass = "yourpassword"
)

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			//1 打开连接
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Println("err: ", err.Error())
				return nil, err
			}
			//2 访问认证
			if _, err = c.Do("AUTH", redisPass); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		// 定期检查redis连接
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := c.Do("PING")
			return err
		},
	}
}

func init() {
	poll = newRedisPool()
}

// RedisPoll : 返回一个redis连接对象
func RedisPoll() *redis.Pool {
	return poll
}
