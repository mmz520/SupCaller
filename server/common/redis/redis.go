package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"SupCaller/common/config"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() {
	cfg := config.Config.Redis

	RDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		IdleTimeout:  parseDuration(cfg.IdleTimeout),
		ReadTimeout:  parseDuration(cfg.ReadTimeout),
		WriteTimeout: parseDuration(cfg.WriteTimeout),
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
}

// parseDuration 解析时间字符串
func parseDuration(durationStr string) time.Duration {
	d, err := time.ParseDuration(durationStr)
	if err != nil {
		return 300 * time.Second
	}
	return d
}

// Set 设置键值对，带过期时间
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RDB.Set(ctx, key, value, expiration).Err()
}

// Get 获取键值
func Get(ctx context.Context, key string) (string, error) {
	return RDB.Get(ctx, key).Result()
}

// GetBytes 获取二进制数据
func GetBytes(ctx context.Context, key string) ([]byte, error) {
	return RDB.Get(ctx, key).Bytes()
}

// Del 删除键
func Del(ctx context.Context, keys ...string) error {
	return RDB.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := RDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Incr 自增
func Incr(ctx context.Context, key string) (int64, error) {
	return RDB.Incr(ctx, key).Result()
}

// IncrBy 自增指定值
func IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return RDB.IncrBy(ctx, key, value).Result()
}

// Decr 自减
func Decr(ctx context.Context, key string) (int64, error) {
	return RDB.Decr(ctx, key).Result()
}

// DecrBy 自减指定值
func DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return RDB.DecrBy(ctx, key, value).Result()
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return RDB.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func TTL(ctx context.Context, key string) (time.Duration, error) {
	return RDB.TTL(ctx, key).Result()
}

// SetNX 当键不存在时设置
func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return RDB.SetNX(ctx, key, value, expiration).Result()
}

// Hash相关操作

// HSet 设置哈希字段
func HSet(ctx context.Context, key string, values ...interface{}) error {
	return RDB.HSet(ctx, key, values...).Err()
}

// HGet 获取哈希字段
func HGet(ctx context.Context, key string, field string) (string, error) {
	return RDB.HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希所有字段
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return RDB.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func HDel(ctx context.Context, key string, fields ...string) error {
	return RDB.HDel(ctx, key, fields...).Err()
}

// HExists 检查哈希字段是否存在
func HExists(ctx context.Context, key string, field string) (bool, error) {
	return RDB.HExists(ctx, key, field).Result()
}

// HIncrBy 哈希字段自增
func HIncrBy(ctx context.Context, key string, field string, value int64) (int64, error) {
	return RDB.HIncrBy(ctx, key, field, value).Result()
}

// List相关操作

// LPush 从列表左侧插入
func LPush(ctx context.Context, key string, values ...interface{}) error {
	return RDB.LPush(ctx, key, values...).Err()
}

// RPush 从列表右侧插入
func RPush(ctx context.Context, key string, values ...interface{}) error {
	return RDB.RPush(ctx, key, values...).Err()
}

// LPop 从列表左侧弹出
func LPop(ctx context.Context, key string) (string, error) {
	return RDB.LPop(ctx, key).Result()
}

// RPop 从列表右侧弹出
func RPop(ctx context.Context, key string) (string, error) {
	return RDB.RPop(ctx, key).Result()
}

// LLen 获取列表长度
func LLen(ctx context.Context, key string) (int64, error) {
	return RDB.LLen(ctx, key).Result()
}

// LRange 获取列表指定范围元素
func LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return RDB.LRange(ctx, key, start, stop).Result()
}

// Set相关操作

// SAdd 向集合添加元素
func SAdd(ctx context.Context, key string, members ...interface{}) error {
	return RDB.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有元素
func SMembers(ctx context.Context, key string) ([]string, error) {
	return RDB.SMembers(ctx, key).Result()
}

// SRem 从集合移除元素
func SRem(ctx context.Context, key string, members ...interface{}) error {
	return RDB.SRem(ctx, key, members...).Err()
}

// SIsMember 检查元素是否在集合中
func SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return RDB.SIsMember(ctx, key, member).Result()
}

// SCard 获取集合大小
func SCard(ctx context.Context, key string) (int64, error) {
	return RDB.SCard(ctx, key).Result()
}

// ZSet相关操作

// ZAdd 向有序集合添加元素
func ZAdd(ctx context.Context, key string, members ...*redis.Z) error {
	return RDB.ZAdd(ctx, key, members...).Err()
}

// ZRange 获取有序集合指定范围元素
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return RDB.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores 获取有序集合指定范围元素及分数
func ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return RDB.ZRangeWithScores(ctx, key, start, stop).Result()
}

// ZRem 从有序集合移除元素
func ZRem(ctx context.Context, key string, members ...interface{}) error {
	return RDB.ZRem(ctx, key, members...).Err()
}

// ZScore 获取元素分数
func ZScore(ctx context.Context, key string, member string) (float64, error) {
	return RDB.ZScore(ctx, key, member).Result()
}

// ZRank 获取元素排名
func ZRank(ctx context.Context, key string, member string) (int64, error) {
	return RDB.ZRank(ctx, key, member).Result()
}

// Pub/Sub相关操作

// Publish 发布消息
func Publish(ctx context.Context, channel string, message interface{}) error {
	return RDB.Publish(ctx, channel, message).Err()
}

// Subscribe 订阅频道
func Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return RDB.Subscribe(ctx, channels...)
}

// 事务相关

// Multi 开始事务
func Multi(ctx context.Context) redis.Pipeliner {
	return RDB.Pipeline()
}

// 批量操作

// MSet 批量设置键值对
func MSet(ctx context.Context, values ...interface{}) error {
	return RDB.MSet(ctx, values...).Err()
}

// MGet 批量获取键值
func MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return RDB.MGet(ctx, keys...).Result()
}

// Scan 扫描键
func Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return RDB.Scan(ctx, cursor, match, count).Result()
}

// FlushDB 清空当前数据库
func FlushDB(ctx context.Context) error {
	return RDB.FlushDB(ctx).Err()
}

// FlushAll 清空所有数据库
func FlushAll(ctx context.Context) error {
	return RDB.FlushAll(ctx).Err()
}

// Close 关闭Redis连接
func Close() error {
	return RDB.Close()
}