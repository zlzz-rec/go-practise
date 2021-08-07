package redis

import (
	"github.com/gomodule/redigo/redis"
)

type NormalTemplate struct {
	pool redis.Pool
}

func NewNormalTemplate(url string, password string, database int, opts ...Option) *NormalTemplate {
	s := &NormalTemplate{}
	s.pool = newPool(url, password, database, opts...)
	return s
}

func (s *NormalTemplate) getPool() *redis.Pool {
	return &s.pool
}
