package redis

import (
	"github.com/gomodule/redigo/redis"
)

type NormalTemplate struct {
	pool redis.Pool
}

func NewNormalTemplate(url string, config *Config) *NormalTemplate {
	s := &NormalTemplate{}
	s.pool = newPool(url, config.Password, config.RedisMaxIdleConnection, config.RedisDatabase,
		config.RedisIdleTimeout, config.RedisConnectTimeout, config.RedisReadTimeout, config.RedisWriteTimeout)
	return s
}

func (s *NormalTemplate) getPool() *redis.Pool {
	return &s.pool
}
