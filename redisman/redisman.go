// Package mongodb provides support for accessing and executing commands against
// a mongoDB database
package redisman

import (
	// "fmt"
	"log"
	// "strings"
	"github.com/garyburd/redigo/redis"
	"time"
)

// default configuration
const (
	maxIdle     = 100
	maxActive   = 1000
	idleTimeout = 180
)

const (
	host     = "localhost"
	port     = "6379"
	password = ""
	database = 0
	timeout  = 60
)

var (
	redisPool *redis.Pool
)

type (
	RedisConfig struct {
		Host     string
		Port     string
		Password string
		Database int
		Timeout  time.Duration
	}

	PoolConfig struct {
		MaxIdle     int
		MaxActive   int
		IdleTimeout time.Duration
	}
)

func Startup(config *RedisConfig, pool *PoolConfig) error {
	log.Println("Redis : Startup : Started")

	if redisPool != nil {
		return nil
	}

	if config == nil {
		config = &RedisConfig{
			Host:     host,
			Port:     port,
			Password: password,
			Database: database,
			Timeout:  timeout * time.Second,
		}
	}

	if pool == nil {
		pool = &PoolConfig{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: idleTimeout * time.Second,
		}
	}

	// options[0] = redis.DialReadTimeout(config.Timeout)
	// options[0] = redis.DialWriteTimeout(config.Timeout)
	optTimout := redis.DialConnectTimeout(config.Timeout)
	optDatabase := redis.DialDatabase(config.Database)
	optPassword := redis.DialPassword(config.Password)
	network := "tcp"
	address := config.Host + ":" + config.Port

	redisPool = &redis.Pool{
		// 从配置文件获取 maxidle以及 maxactive，取不到则用后面的默认值
		MaxIdle:     pool.MaxIdle,
		MaxActive:   pool.MaxActive,
		IdleTimeout: pool.IdleTimeout,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(network, address, optTimout, optDatabase, optPassword)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	log.Println("Redis : Startup : Completed")
	return nil
}

func GetConn() redis.Conn {
	return redisPool.Get()
}

// Gracefully shutdown
// func Shutdown() {
// 	redis.Conn.Close()
// }
