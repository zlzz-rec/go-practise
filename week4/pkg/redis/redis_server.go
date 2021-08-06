package redis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	RedisMaxIdleConnection int
	RedisIdleTimeout       time.Duration
	RedisConnectTimeout    time.Duration
	RedisReadTimeout       time.Duration
	RedisWriteTimeout      time.Duration
	RedisDatabase          int
	Password               string
}

// NewConfig 初始化redis配置
func NewConfig(maxIdleConnection int, idleTimeout time.Duration, connectTimeout time.Duration, readTimeout time.Duration, writeTimeout time.Duration, redisDB int, password string) *Config {

	return &Config{
		RedisMaxIdleConnection: maxIdleConnection,
		RedisIdleTimeout:       idleTimeout,
		RedisConnectTimeout:    connectTimeout,
		RedisReadTimeout:       readTimeout,
		RedisWriteTimeout:      writeTimeout,
		RedisDatabase:          redisDB,
		Password:               password,
	}
}

// RedisTemplate 模板方法模式 server只要是实现RedisTemplate接口的type都可以传入
//
// 使用说明
//
//  server := redis.NewRedisServer()
//  singletonTemp := redis.GetInstance(redisIps []string, config *Config)
//  server.SetRedisTemplate(singletonTemp)
type RedisTemplate interface {
	getPool() *redis.Pool
}

// RedisServer 组合RedisTemplate
type RedisServer struct {
	server RedisTemplate
}

// NewRedisServer redis服务的封装，同时支持单机版服务和分布式的集群
// 在初始化服务时，需要调用server.SetRedisTemplate，进行设置模板
// 单机版：初始化模板 redis.NewNormalTemplate
//
//  server := redis.NewRedisServer()
//  normalTemp := NewNormalTemplate(url, cfg)
//  server.SetRedisTemplate(normalTemp)
//
// 集群版：初始化模板 redis.GetInstance
//
//  server := redis.NewRedisServer()
//  singletonTemp := redis.GetInstance(ips, cfg)
//  server.SetRedisTemplate(singletonTemp)
func NewRedisServer() *RedisServer {
	return &RedisServer{}
}

// SetRedisTemplate SetRedisTemplate
func (s *RedisServer) SetRedisTemplate(t RedisTemplate) {
	s.server = t
}

// Get Get
func (s *RedisServer) Get(key string) (string, error) {
	pool := s.server.getPool()
	data, err := do(pool, "GET", key)
	if err != nil {
		return "", err
	}
	if data != nil {
		ret, err := redis.String(data, nil)
		if err != nil {
			if err == redis.ErrNil {
				return "", nil
			}
			return "", err
		}
		return ret, nil
	}
	return "", nil
}

// Ping ping
func (s *RedisServer) Ping() (string, error) {
	pool := s.server.getPool()
	return redis.String(do(pool, "PING"))
}

// Set Set
func (s *RedisServer) Set(key string, value string) error {
	pool := s.server.getPool()
	_, err := redis.String(do(pool, "SET", key, value))
	return err
}

// SetEx SetEx
func (s *RedisServer) SetEx(key string, value string, seconds int64) error {
	if seconds <= 0 {
		return s.Set(key, value)
	}
	pool := s.server.getPool()
	sec := strconv.FormatInt(seconds, 10)
	_, err := redis.String(do(pool, "SETEX", key, sec, value))
	return err
}

// SetNx SetNx
func (s *RedisServer) SetNx(key string, value string) (ok bool, err error) {
	pool := s.server.getPool()
	iRet, err := redis.Int(do(pool, "SETNX", key, value))
	if err != nil {
		return false, err
	}
	if iRet == 0 {
		return false, nil
	}
	return true, nil
}

// TryLock 分布式锁
func (s *RedisServer) TryLock(key string, value string, seconds int64) (ok bool, err error) {
	pool := s.server.getPool()
	result, err := redis.String(do(pool, "SET", key, value, "NX", "EX", seconds))
	if err != nil {
		if err == redis.ErrNil { // redigo: nil returned 锁被占用
			return false, nil
		} else { // 运行错误
			return false, err
		}
	}
	if result != "" { // 获取锁成功
		return true, nil
	}
	return false, nil
}

// ReleaseRedisLock 删除redis锁
func (s *RedisServer) ReleaseRedisLock(strKey string) (int, error) {
	pool := s.server.getPool()
	iRet, err := redis.Int(do(pool, "del", strKey))
	return iRet, err
}

// Exists Exists
func (s *RedisServer) Exists(key string) (int, error) {
	pool := s.server.getPool()
	return redis.Int(do(pool, "EXISTS", key))
}

// Expire Expire
func (s *RedisServer) Expire(strKey string, iSecond int) (int, error) {
	pool := s.server.getPool()
	iRet, err := redis.Int(do(pool, "expire", strKey, iSecond))
	return iRet, err
}

// Del Del
func (s *RedisServer) Del(strKey string) (int, error) {
	pool := s.server.getPool()
	iRet, err := redis.Int(do(pool, "del", strKey))
	return iRet, err
}

// BatchDel 批量删除
func (s *RedisServer) BatchDel(strKeys []string) ([]int64, error) {
	conn := s.server.getPool().Get()
	if len(strKeys) == 0 {
		return nil, nil
	}
	defer conn.Close()
	conn.Send("MULTI")
	for _, key := range strKeys {
		conn.Send("DEL", key)
	}
	reply, err := conn.Do("EXEC")
	if err != nil {
		return nil, err
	}
	flags := []int64{}
	for _, v := range reply.([]interface{}) {
		flags = append(flags, v.(int64))
	}
	return flags, err
}

// Mget Mget
func (s *RedisServer) Mget(keys ...string) ([]string, error) {
	pool := s.server.getPool()

	args := []interface{}{}
	for _, key := range keys {
		args = append(args, key)
	}

	rv, err := do(pool, "MGET", args...)
	iRet, err := redis.Values(rv, err)
	if err != nil {
		return nil, err
	}
	strBuf := make([]string, 0, len(iRet))
	for _, value := range iRet {
		if value != nil {
			strBuf = append(strBuf, string(uint8s2Bytes(value.([]uint8))))
		} else {
			strBuf = append(strBuf, "")
		}

	}
	return strBuf, nil
}

func uint8s2Bytes(data []uint8) []byte {
	buf := make([]byte, 0, len(data))
	for _, v := range data {
		buf = append(buf, byte(v))
	}
	return buf
}

// Sadd Sadd
func (s *RedisServer) Sadd(key string, values ...string) (int64, error) {
	pool := s.server.getPool()
	args := redis.Args{}.Add(key).AddFlat(values)
	numberAdded, err := redis.Int64(do(pool, "SADD", args...))
	return numberAdded, err
}

// Srem Srem
func (s *RedisServer) Srem(key string, values ...string) (int64, error) {
	pool := s.server.getPool()
	args := redis.Args{}.Add(key).AddFlat(values)
	numberRemoved, err := redis.Int64(do(pool, "Srem", args...))
	return numberRemoved, err
}

// Scard Scard
func (s *RedisServer) Scard(key string) (int64, error) {
	pool := s.server.getPool()
	length, err := redis.Int64(do(pool, "SCARD", key))
	return length, err
}

// Sdiff Sdiff
func (s *RedisServer) Sdiff(key string, exceptKey []string) ([]string, error) {
	pool := s.server.getPool()
	args := redis.Args{}.Add(key).AddFlat(exceptKey)
	values, err := redis.Strings(do(pool, "SDIFF", args...))
	if err != nil {
		return []string{}, err
	}

	return values, nil
}

// SrandMember SrandMember
func (s *RedisServer) SrandMember(key string, count int64) ([]string, error) {

	pool := s.server.getPool()
	values, err := redis.Strings(do(pool, "SRANDMEMBER", key, count))
	if err != nil {
		return []string{}, err
	}

	return values, nil
}

// Smembers Smembers
func (s *RedisServer) Smembers(key string) ([]string, error) {

	pool := s.server.getPool()
	values, err := redis.Strings(do(pool, "SMEMBERS", key))
	if err != nil {
		return []string{}, err
	}

	return values, nil
}

// Sismember 成员是否属于set
func (s *RedisServer) Sismember(key string, member string) (bool, error) {

	pool := s.server.getPool()
	intVal, err := redis.Int64(do(pool, "SISMEMBER", key, member))
	if err != nil {
		return false, err
	}
	if intVal == 1 {
		return true, nil
	}
	return false, nil
}

// Lrange Lrange
func (s *RedisServer) Lrange(key string, start, stop int) ([]string, error) {
	pool := s.server.getPool()
	strValue, err := redis.Strings(do(pool, "LRANGE", key, start, stop))
	if err != nil {
		return nil, err
	}

	return strValue, nil
}

// Lpush Lpush
func (s *RedisServer) Lpush(key string, values ...string) (int64, error) {
	pool := s.server.getPool()
	args := redis.Args{}.Add(key).AddFlat(values)
	length, err := redis.Int64(do(pool, "LPUSH", args...))
	return length, err
}

// Rpush Rpush
func (s *RedisServer) Rpush(key string, values ...string) (int64, error) {
	pool := s.server.getPool()
	args := redis.Args{}.Add(key).AddFlat(values)
	length, err := redis.Int64(do(pool, "RPUSH", args...))
	return length, err
}

// Hset Hset
func (s *RedisServer) Hset(key string, field string, value string) error {
	pool := s.server.getPool()
	_, err := do(pool, "HSET", key, field, value)
	return err
}

// Hexists field是否存在
func (s *RedisServer) Hexists(key string, field string) (bool, error) {

	pool := s.server.getPool()
	intVal, err := redis.Int64(do(pool, "HEXISTS", key, field))
	if err != nil {
		return false, err
	}
	if intVal == 1 {
		return true, nil
	}
	return false, nil
}

// Hdel Hdel
func (s *RedisServer) Hdel(key string, field string) error {
	pool := s.server.getPool()
	_, err := redis.Int(do(pool, "HDEL", key, field))
	return err
}

// HdelMulti 移除多个hset的field
func (s *RedisServer) HdelMulti(keyName string, members []string) (int64, error) {
	params := make([]interface{}, len(members)+1)
	params[0] = keyName
	for i, v := range members {
		params[i+1] = v
	}

	pool := s.server.getPool()
	value, err := redis.Int64(do(pool, "HDEL", params...))
	if err != nil {
		return 0, err
	}

	return value, nil
}

// Hget Hget
func (s *RedisServer) Hget(key string, field string) (string, error) {
	pool := s.server.getPool()
	strValue, err := redis.String(do(pool, "HGET", key, field))
	if err != nil {
		if err == redis.ErrNil {
			return "", nil
		}
		return "", err
	}

	return strValue, nil
}

// Hgetall HGETALL
func (s *RedisServer) Hgetall(key string) (map[string]string, error) {
	pool := s.server.getPool()

	mapValue, err := redis.StringMap(do(pool, "HGETALL", key))
	if err != nil {
		return nil, err
	}

	return mapValue, nil
}

// MultiGet 获取多个key的值
func (s *RedisServer) MultiGet(keys []string) ([]string, error) {
	// 拼接脚本
	keyCount := len(keys)
	callArr := []string{}
	for i := 1; i <= keyCount; i++ {
		callArr = append(callArr, fmt.Sprintf(`redis.call('GET', KEYS[%d])`, i))
	}
	scriptStr := fmt.Sprintf(`return{%s}`, strings.Join(callArr, ","))

	// 构建参数
	args := redis.Args{}.AddFlat(keys)

	return runScriptForStrings(s, keyCount, scriptStr, args)
}

// MultiHgetall 获取多个key的哈希结构
func (s *RedisServer) MultiHgetall(keys []string) ([]map[string]string, error) {
	// 拼接脚本
	keyCount := len(keys)
	callArr := []string{}
	for i := 1; i <= keyCount; i++ {
		callArr = append(callArr, fmt.Sprintf(`redis.call('HGETALL', KEYS[%d])`, i))
	}
	scriptStr := fmt.Sprintf(`return{%s}`, strings.Join(callArr, ","))

	// 构建参数
	args := redis.Args{}.AddFlat(keys)

	return runScriptForMapArr(s, keyCount, scriptStr, args)

}

// MultiHget 获取哈希结构的多个feild的value
func (s *RedisServer) MultiHget(key string, fields []string) ([]string, error) {
	memberNum := len(fields)
	callArr := []string{}
	for i := 1; i <= memberNum; i++ {
		callArr = append(callArr, fmt.Sprintf(`redis.call('HGET', KEYS[1], ARGV[%d])`, i))
	}
	scriptStr := fmt.Sprintf(`return{%s}`, strings.Join(callArr, ","))

	// 构建参数
	args := redis.Args{}.Add(key).AddFlat(fields)

	return runScriptForStrings(s, 1, scriptStr, args)
}

// Hmset Hmset
func (s *RedisServer) Hmset(strKey string, mapAttrs map[string]string) (string, error) {
	pool := s.server.getPool()
	param := redis.Args{}.Add(strKey).AddFlat(mapAttrs)
	strRet, err := redis.String(do(pool, "HMSET", param...))
	if err != nil {
		return "", err
	}

	return strRet, nil
}

// Hincrby HINCRBY
func (s *RedisServer) Hincrby(strKey string, field string, increment int64) (int64, error) {
	pool := s.server.getPool()
	value, err := redis.Int64(do(pool, "HINCRBY", strKey, field, increment))
	if err != nil {
		return 0, err
	}

	return value, nil
}

// ZrangeByLex 返回指定区间内的成员，支持offset，limit
// 在全体成员内，返回数据
//
//  redis.ZrangeByLex(key, "-", "+", true, offset, limit)
//
// 在指定区间 a（包含）-c（不包含） 内返回数据
//
//  redis.ZrangeByLex(key, "[a", "(c", true, offset, limit)
func (s *RedisServer) ZrangeByLex(strKey, min, max string, limitFlag bool, offset, count uint64) ([]string, error) {
	pool := s.server.getPool()
	var strValue []string
	var err error
	if limitFlag {
		strValue, err = redis.Strings(do(pool, "ZRANGEBYLEX", strKey, min, max, "LIMIT", offset, count))
	} else {
		strValue, err = redis.Strings(do(pool, "ZRANGEBYLEX", strKey, min, max))
	}
	if err != nil {
		return strValue, err
	}
	return strValue, nil
}

// ZrevrangeByLex 倒序返回指定区间内的成员，支持offset，limit
// 在全体成员内，返回数据
//
//  redis.ZrevrangeByLex(key, "-", "+", true, offset, limit)
//
// 在指定区间 a（包含）-c（不包含） 内返回数据
//
//  redis.ZrevrangeByLex(key, "[a", "(c", true, offset, limit)
func (s *RedisServer) ZrevrangeByLex(strKey, max, min string, limitFlag bool, offset, count uint64) ([]string, error) {
	pool := s.server.getPool()
	var strValue []string
	var err error
	if limitFlag {
		strValue, err = redis.Strings(do(pool, "ZREVRANGEBYLEX", strKey, max, min, "LIMIT", offset, count))
	} else {
		strValue, err = redis.Strings(do(pool, "ZREVRANGEBYLEX", strKey, max, min))
	}
	if err != nil {
		return strValue, err
	}
	return strValue, nil
}

// ZremRangeByScore 移除zset中，指定分数（score）区间内的所有成员
func (s *RedisServer) ZremRangeByScore(keyName string, min, max string) (int64, error) {
	pool := s.server.getPool()
	i, e := do(pool, "zremrangebyscore", keyName, min, max)
	mapValues, err := redis.Int64(i, e)
	if err != nil {
		return 0, err
	}

	return mapValues, nil
}

// Zadd zset添加元素
func (s *RedisServer) Zadd(keyName string, score string, member string) (int64, error) {
	pool := s.server.getPool()
	value, err := redis.Int64(do(pool, "ZADD", keyName, score, member))
	if err != nil {
		return value, err
	}

	return value, nil
}

// Zadds 批量添加成员到有序集合
func (s *RedisServer) Zadds(keyName string, input map[string]string) (interface{}, error) {
	output := []interface{}{}
	output = append(output, keyName)
	for key, value := range input {
		output = append(output, value, key)
	}
	pool := s.server.getPool()
	//println(args)
	mapValues, err := do(pool, "ZADD", output...)
	if err != nil {
		return nil, err
	}

	return mapValues, nil
}

// zrevrange 从有序集合中查询数据，递减排列,获取索引区间的成员
func (s *RedisServer) Zrevrange(keyName string, start int64, stop int64) ([]string, error) {
	pool := s.server.getPool()
	result, err := redis.Strings(do(pool, "zrevrange", keyName, start, stop))
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ZrangeByScore 递增排列, 通过分数返回有序集合指定区间内的成员
func (s *RedisServer) ZrangeByScore(keyName string, min int64, max int64) ([]string, error) {
	pool := s.server.getPool()
	values, err := redis.Strings(do(pool, "ZRANGEBYSCORE", keyName, min, max))
	if err != nil {
		return []string{}, err
	}

	return values, nil
}

// ZrevrangeByPreMemberId  倒序获取preMember下标后面pagesize个 Member成员
func (s *RedisServer) ZrevrangeByPreMember(key string, preMember uint64, pagesize int64) ([]uint64, error) {
	var err error
	if preMember != 0 {
		index, err := s.Zrevrank(key, strconv.FormatUint(preMember, 10)) //获取成员的下标
		if err == nil {
			if index >= 0 {
				return searchMembers(s, key, index+1, index+pagesize, true)
			}
		}
	} else if preMember <= 0 {
		return searchMembers(s, key, 0, pagesize-1, true)
	}
	return nil, err
}

// ZrangeByPreMember 正序获取preMember(成员)后面 pagesize个 Member成员
func (s *RedisServer) ZrangeByPreMember(key string, preMember uint64, pagesize int64) ([]uint64, error) {
	var err error
	if preMember == 0 {
		return searchMembers(s, key, 0, pagesize-1, false)
	}
	var preIndex int64
	preIndex, err = s.Zrevrank(key, strconv.FormatUint(preMember, 10)) //获取成员的下标
	if err != nil {
		return nil, err
	}
	if preIndex >= 0 {
		return searchMembers(s, key, preIndex+1, preIndex+pagesize, false)
	}
	return nil, nil
}

// searchMembers 正序(或逆序)获取preMember(成员)后面pagesize个 Member成员
func searchMembers(s *RedisServer, key string, start int64, end int64, reverseOrder bool) ([]uint64, error) {
	var result []string
	var err error
	if reverseOrder {
		result, err = s.Zrevrange(key, start, end)
	} else {
		result, err = s.Zrange(key, start, end)
	}
	if err != nil {
		return nil, err
	}
	if len(result) <= 0 {
		return nil, nil
	}
	ids := []uint64{}
	for _, v := range result {
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, i)
	}
	return ids, nil
}

// Zrange 递增排列, 获取索引区间的成员
func (s *RedisServer) Zrange(keyName string, start int64, stop int64) ([]string, error) {
	pool := s.server.getPool()
	values, err := redis.Strings(do(pool, "zrange", keyName, start, stop))
	if err != nil {
		return []string{}, err
	}

	return values, nil
}

// ZrevrangeWithscores 递减排列, 获取索引区间的成员和权重
func (s *RedisServer) ZrevrangeWithscores(keyName string, start int64, stop int64) ([]string, error) {
	pool := s.server.getPool()
	values, err := redis.Strings(do(pool, "zrevrange", keyName, start, stop, "WITHSCORES"))
	if err != nil {
		return []string{}, err
	}

	return values, nil
}

// ZrangeWithscores 递增排列, 获取索引区间的成员和权重
func (s *RedisServer) ZrangeWithscores(keyName string, start int64, stop int64) ([]string, error) {
	pool := s.server.getPool()
	values, err := redis.Strings(do(pool, "zrange", keyName, start, stop, "WITHSCORES"))
	if err != nil {
		return []string{}, err
	}

	return values, nil
}

// Zscore 返回zset成员的score值
func (s *RedisServer) Zscore(keyName string, member string) (string, error) {
	pool := s.server.getPool()
	mapValues, err := redis.String(do(pool, "zscore", keyName, member))
	if err != nil {
		return "", err
	}

	return mapValues, nil
}

// ZscoreMulti 返回zset多个成员的score值
func (s *RedisServer) ZscoreMulti(keyName string, members []string) ([]string, error) {

	// 拼接脚本
	memberNum := len(members)
	callArr := []string{}
	for i := 1; i <= memberNum; i++ {
		callArr = append(callArr, fmt.Sprintf(`redis.call('ZSCORE', KEYS[1], ARGV[%d])`, i))
	}
	scriptStr := fmt.Sprintf(`return{%s}`, strings.Join(callArr, ","))

	// 构建参数
	args := redis.Args{}.Add(keyName).AddFlat(members)

	return runScriptForStrings(s, 1, scriptStr, args)
}

// Zrem 移除zset成员
func (s *RedisServer) Zrem(keyName string, member string) (int64, error) {
	pool := s.server.getPool()
	i, e := do(pool, "zrem", keyName, member)
	mapValues, err := redis.Int64(i, e)
	if err != nil {
		return 0, err
	}

	return mapValues, nil
}

// ZremMulti 移除多个zset成员
func (s *RedisServer) ZremMulti(keyName string, members []string) (int64, error) {
	params := make([]interface{}, len(members)+1)
	params[0] = keyName
	for i, v := range members {
		params[i+1] = v
	}

	pool := s.server.getPool()
	value, err := redis.Int64(do(pool, "zrem", params...))
	if err != nil {
		return 0, err
	}

	return value, nil
}

// Zcard 查看zset的成员个数
func (s *RedisServer) Zcard(keyName string) (int64, error) {
	pool := s.server.getPool()
	count, err := redis.Int64(do(pool, "zcard", keyName))
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ZcardMulti 查看多个zset的成员个数
func (s *RedisServer) ZcardMulti(keys []string) ([]int64, error) {
	// 拼接脚本
	keyCount := len(keys)
	callArr := []string{}
	for i := 1; i <= keyCount; i++ {
		callArr = append(callArr, fmt.Sprintf(`redis.call('ZCARD', KEYS[%d])`, i))
	}
	scriptStr := fmt.Sprintf(`return{%s}`, strings.Join(callArr, ","))

	// 构建参数
	args := redis.Args{}.AddFlat(keys)

	return runScriptForIntArr(s, keyCount, scriptStr, args)
}

// SismemberMulti 批量判断成员是否属于set
func (s *RedisServer) SismemberMulti(key string, members []string) (map[string]bool, error) {
	result := make(map[string]bool)

	callArr := []string{}
	for i := 1; i <= len(members); i++ {
		callArr = append(callArr, fmt.Sprintf(`redis.call('SISMEMBER', KEYS[1], ARGV[%d])`, i))
	}
	scriptStr := fmt.Sprintf(`return{%s}`, strings.Join(callArr, ","))

	// 构建参数
	args := redis.Args{}.Add(key).AddFlat(members)

	strs, err := runScriptForIntArr(s, 1, scriptStr, args)
	if err != nil {
		return result, err
	}

	for i := range strs {
		result[members[i]] = (strs[i] > 0)
	}

	return result, nil
}

// Zinterstore zset的交集操作
func (s *RedisServer) Zinterstore(desKey string, orgKeys []string) (int64, error) {

	keys := make([]interface{}, len(orgKeys)+2)
	keys[0] = desKey
	keys[1] = len(orgKeys)
	for i, v := range orgKeys {
		keys[i+2] = v
	}

	pool := s.server.getPool()
	num, err := redis.Int64(do(pool, "ZINTERSTORE", keys...))
	if err != nil {
		return 0, err
	}

	return num, nil
}

// Zunionstore zset的并集操作
func (s *RedisServer) Zunionstore(desKey string, orgKeys []string) (int64, error) {

	keys := make([]interface{}, len(orgKeys)+2)
	keys[0] = desKey
	keys[1] = len(orgKeys)
	for i, v := range orgKeys {
		keys[i+2] = v
	}

	pool := s.server.getPool()
	num, err := redis.Int64(do(pool, "ZUNIONSTORE", keys...))
	if err != nil {
		return 0, err
	}

	return num, nil
}

// Zincrby 改变zset成员的score (正数是增加, 负数是减少)
func (s *RedisServer) Zincrby(key string, score string, member string) (string, error) {

	keys := make([]interface{}, 3)
	keys[0] = key
	keys[1] = score
	keys[2] = member

	pool := s.server.getPool()
	value, err := redis.String(do(pool, "ZINCRBY", keys...))
	if err != nil {
		return "0", err
	}

	return value, nil
}

// Incrby  INCRBY(正数是增加, 负数是减少)
func (s *RedisServer) Incrby(key string, increment int64, expiredSecend int64) (int64, error) {

	pool := s.server.getPool()
	keys := make([]interface{}, 2)
	keys[0] = key
	keys[1] = increment
	value, err := redis.Int64(do(pool, "INCRBY", keys...))
	if err != nil {
		return 0, err
	}
	if expiredSecend > 0 {
		_, err = redis.Int64(do(pool, "EXPIRE", []interface{}{key, expiredSecend}...))
	}
	if err != nil {
		return 0, err
	}

	return value, nil
}

// Incrbyfloat Incrbyfloat
func (s *RedisServer) Incrbyfloat(key string, increment float64, expiredSecend int64) (string, error) {

	keys := make([]interface{}, 2)
	keys[0] = key
	keys[1] = increment

	pool := s.server.getPool()
	value, err := redis.String(do(pool, "INCRBYFLOAT", keys...))
	if err != nil {
		return "0", err
	}
	if expiredSecend > 0 {
		_, err = redis.Int64(do(pool, "EXPIRE", []interface{}{key, expiredSecend}...))
	}
	if err != nil {
		return "0", err
	}

	return value, nil
}

// Zcount 查看指定分数区间内的成员个数
func (s *RedisServer) Zcount(keyName string, start int64, stop int64) (int64, error) {
	pool := s.server.getPool()
	num, err := redis.Int64(do(pool, "zcount", keyName, start, stop))
	if err != nil {
		return 0, err
	}

	return num, nil
}

// Zrank 获取zset成员的下标位置，如果值不存在返回-1
func (s *RedisServer) Zrank(keyName string, member string) (int64, error) {
	pool := s.server.getPool()
	index, err := do(pool, "zrank", keyName, member)
	if err != nil {
		return -1, err
	}
	if index != nil {
		num, err := redis.Int64(index, err)
		return num, err
	}
	return -1, nil
}

// ZrankMulti 获取zset成员的下标位置，如果值不存在返回-1
func (s *RedisServer) ZrankMulti(keyName string, members []string) ([]int64, error) {
	// 拼接脚本
	memberNum := len(members)
	callArr := []string{}
	for i := 1; i <= memberNum; i++ {
		callArr = append(callArr, fmt.Sprintf(`redis.call('ZRANK', KEYS[1], ARGV[%d])`, i))
	}
	scriptStr := fmt.Sprintf(`return{%s}`, strings.Join(callArr, ","))

	// 构建参数
	args := redis.Args{}.Add(keyName).AddFlat(members)

	return runScriptForInt64(s, 1, scriptStr, args)
}

// Zrevrank 获取zset成员的倒序下标位置，如果值不存在返回-1
func (s *RedisServer) Zrevrank(keyName string, member string) (int64, error) {
	pool := s.server.getPool()
	index, err := do(pool, "ZREVRANK", keyName, member)
	if err != nil {
		return -1, err
	}
	if index != nil {
		num, err := redis.Int64(index, err)
		return num, err
	}
	return -1, nil
}

// HmgetArgs HmgetArgs
func (s *RedisServer) HmgetArgs(args ...interface{}) (interface{}, error) {
	pool := s.server.getPool()
	mapValues, err := do(pool, "hmget", args...)
	if err != nil {
		return nil, err
	}

	return mapValues, nil
}

// WriteCache 写redis，value为string类型时，直接写redis，其他类型先转成json结构，再写redis
func (s *RedisServer) WriteCache(key string, value interface{}, seconds int64) error {
	str, ok := value.(string)
	if ok { // 类型为string，直接写redis
		err := s.SetEx(key, str, seconds) // 写缓存
		if err != nil {
			return err
		}
		return nil
	}

	// 非string，转json再写redis
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	str = string(bytes)
	err = s.SetEx(key, str, seconds) // 写缓存
	if err != nil {
		return err
	}
	return nil
}

// runScriptForStrings 运行脚本得到string数组
func runScriptForStrings(s *RedisServer, keyCount int, scriptStr string, args redis.Args) ([]string, error) {
	return redis.Strings(runScript(s, keyCount, scriptStr, args))
}

// runScriptForIntArr 运行脚本得到int64数组
func runScriptForIntArr(s *RedisServer, keyCount int, scriptStr string, args redis.Args) ([]int64, error) {
	return redis.Int64s(runScript(s, keyCount, scriptStr, args))
}

// runScriptForMapArr 运行脚本得到map[string]string数组
func runScriptForMapArr(s *RedisServer, keyCount int, scriptStr string, args redis.Args) ([]map[string]string, error) {
	values, err := redis.Values(runScript(s, keyCount, scriptStr, args))
	result := []map[string]string{}
	for _, mb := range values {
		m, err := redis.StringMap(mb, err)
		if err != nil {
			return result, err
		}
		result = append(result, m)
	}
	return result, nil
}

// runScriptForInt64 运行脚本得到int64数组
func runScriptForInt64(s *RedisServer, keyCount int, scriptStr string, args redis.Args) ([]int64, error) {
	indexes, err := runScript(s, keyCount, scriptStr, args)
	if err != nil {
		return nil, err
	}

	arr, ok := indexes.([]interface{})
	if !ok {
		return nil, fmt.Errorf("runScriptForInt64 err:%s", arr)
	}

	result := []int64{}
	for _, inter := range arr {
		switch t := inter.(type) {
		case int64:
			result = append(result, t)
		case nil:
			result = append(result, -1)
		}
	}

	return result, nil
}

// runScript 运行脚本
func runScript(s *RedisServer, keyCount int, scriptStr string, args redis.Args) (interface{}, error) {
	// 创建脚本
	script := redis.NewScript(keyCount, scriptStr)

	// 创建连接
	conn := s.server.getPool().Get()
	defer conn.Close()

	// 运行脚本
	return script.Do(conn, args...)
}

func newPool(address, password string, maxIdleConnection int, db int, idleTimeout, connectTimeout, readTimeout, writeTimeout time.Duration) redis.Pool {
	return redis.Pool{
		MaxIdle:     maxIdleConnection,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			conn, err := newConn("tcp", address, password, db, connectTimeout, readTimeout, writeTimeout)
			if err != nil {
				return nil, err
			}
			return conn, err
		},
	}
}

func newConn(network, address, password string, db int, connectTimeout, readTimeout, writeTimeout time.Duration) (redis.Conn, error) {
	return redis.Dial(network, address,
		redis.DialPassword(password),
		redis.DialConnectTimeout(connectTimeout),
		redis.DialReadTimeout(readTimeout),
		redis.DialWriteTimeout(writeTimeout),
		redis.DialDatabase(db))
}

func do(redisPool *redis.Pool, cmd string, args ...interface{}) (interface{}, error) {
	conn := redisPool.Get()
	defer conn.Close()

	retVal, err := conn.Do(cmd, args...)

	return retVal, err
}
