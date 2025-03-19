package zwx

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/valkey-io/valkey-go"
	"log/slog"
	"os"
	"time"
)

type Logger interface {
	Debugf(format string, v ...any)
	Infof(format string, v ...any)
	Errorf(format string, v ...any)
	Fatalf(format string, v ...any)
}
type Storage interface {
	Get(key string) string
	Del(key string)
	SetEX(key string, val string, expire time.Duration)
	SetNX(key string, val string, expire time.Duration) bool
	SAdd(key string, val ...string)
	SRem(key string, val ...string)
	SMembers(key string) []string
	HSet(key string, val map[string]string)
	HGetAll(key string) map[string]string
	HIncrBy(key string, field string, incr int64)
}
type Options struct {
	// 开启调试模式，会打印更多日志
	Debug bool
	// 自定义日志器
	Logger Logger
	// 存储器自定义实现
	Storage Storage
	// 支持valkey客户端，github.com/valkey-io/valkey-go
	ValkeyClient valkey.Client
	// 支持redis客户端，github.com/go-redis/redis/v8
	RedisClient redis.UniversalClient
	// 存储器前缀，默认为空
	StoragePrefix string
	// 刷新token的间隔时间，默认55分钟
	AccessTokenRefresh time.Duration
	// 每次启动前清理缓存，默认false，如果开启，每次启动之前都会遗忘之前托管的app
	AlwaysCleanBeforeStart bool
}

func (o *Options) Validate() {
	if o.Logger == nil {
		o.Logger = &defaultLogger{}
	}
	if o.ValkeyClient != nil {
		o.Storage = &valkeyStorage{o.ValkeyClient}
	} else if o.RedisClient != nil {
		o.Storage = &redisStorage{o.RedisClient}
	} else if o.Storage == nil {
		o.Logger.Fatalf("storage/valkeyClient/redisClient must have one!")
	}
	if o.AccessTokenRefresh == 0 {
		o.AccessTokenRefresh = 55 * time.Minute
	}
}

// 默认日志实现
type defaultLogger struct{}

func (l *defaultLogger) Debugf(format string, v ...any) {
	slog.Debug(fmt.Sprintf(format, v...))
}
func (l *defaultLogger) Infof(format string, v ...any) {
	slog.Info(fmt.Sprintf(format, v...))
}
func (l *defaultLogger) Errorf(format string, v ...any) {
	slog.Error(fmt.Sprintf(format, v...))
}
func (l *defaultLogger) Fatalf(format string, v ...any) {
	l.Errorf(format, v...)
	os.Exit(1)
}

// valkeyStorage
// @Description: valkey 实现
type valkeyStorage struct {
	valkey.Client
}

func (s *valkeyStorage) Get(key string) string {
	str, _ := s.Do(context.TODO(), s.B().Get().Key(key).Build()).ToString()
	return str
}
func (s *valkeyStorage) Del(key string) {
	s.Do(context.TODO(), s.B().Del().Key(key).Build())
}
func (s *valkeyStorage) SetEX(key string, val string, expire time.Duration) {
	s.Do(context.TODO(), s.B().Set().Key(key).Value(val).Ex(expire).Build())
}
func (s *valkeyStorage) SetNX(key string, val string, expire time.Duration) bool {
	res, _ := s.Do(context.TODO(), s.B().Set().Key(key).Value(val).Nx().Ex(expire).Build()).AsBool()
	return res
}
func (s *valkeyStorage) SAdd(key string, val ...string) {
	s.Do(context.TODO(), s.B().Sadd().Key(key).Member(val...).Build())
}
func (s *valkeyStorage) SRem(key string, val ...string) {
	s.Do(context.TODO(), s.B().Srem().Key(key).Member(val...).Build())
}
func (s *valkeyStorage) SMembers(key string) []string {
	res, _ := s.Do(context.TODO(), s.B().Smembers().Key(key).Build()).AsStrSlice()
	return res
}
func (s *valkeyStorage) HSet(key string, val map[string]string) {
	cmd := s.B().Hset().Key(key).FieldValue()
	for k, v := range val {
		cmd.FieldValue(k, v)
	}
	s.Do(context.TODO(), cmd.Build())
}
func (s *valkeyStorage) HGetAll(key string) map[string]string {
	res, _ := s.Do(context.TODO(), s.B().Hgetall().Key(key).Build()).AsStrMap()
	return res
}
func (s *valkeyStorage) HIncrBy(key string, field string, incr int64) {
	s.Do(context.TODO(), s.B().Hincrby().Key(key).Field(field).Increment(incr).Build())
}

// redisStorage
// @Description: redis 实现
type redisStorage struct {
	c redis.UniversalClient
}

func (s *redisStorage) Get(key string) string {
	return s.c.Get(context.TODO(), key).Val()
}
func (s *redisStorage) Del(key string) {
	s.c.Del(context.TODO(), key)
}
func (s *redisStorage) SetEX(key string, val string, expire time.Duration) {
	s.c.Set(context.TODO(), key, val, expire)
}
func (s *redisStorage) SetNX(key string, val string, expire time.Duration) bool {
	return s.c.SetNX(context.TODO(), key, val, expire).Val()
}
func (s *redisStorage) SAdd(key string, val ...string) {
	s.c.SAdd(context.TODO(), key, val)
}
func (s *redisStorage) SRem(key string, val ...string) {
	s.c.SRem(context.TODO(), key, val)
}
func (s *redisStorage) SMembers(key string) []string {
	return s.c.SMembers(context.TODO(), key).Val()
}
func (s *redisStorage) HSet(key string, val map[string]string) {
	s.c.HSet(context.TODO(), key, val)
}
func (s *redisStorage) HGetAll(key string) map[string]string {
	return s.c.HGetAll(context.TODO(), key).Val()
}
func (s *redisStorage) HIncrBy(key string, field string, incr int64) {
	s.c.HIncrBy(context.TODO(), key, field, incr)
}
