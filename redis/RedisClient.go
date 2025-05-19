package redis

import (
	"context"
	redisRaw "github.com/go-redis/redis/v8"
	"os"
	"strings"
	"time"
)

type Limiter = redisRaw.Limiter
type Options = redisRaw.Options
type Pipeliner = redisRaw.Pipeliner
type Cmder = redisRaw.Cmder
type Cmd = redisRaw.Cmd
type PubSub = redisRaw.PubSub
type PoolStats = redisRaw.PoolStats
type Conn = redisRaw.Conn
type XAddArgs = redisRaw.XAddArgs
type StringCmd = redisRaw.StringCmd
type CommandsInfoCmd = redisRaw.CommandsInfoCmd
type StringSliceCmd = redisRaw.StringSliceCmd
type FloatCmd = redisRaw.FloatCmd
type IntCmd = redisRaw.IntCmd
type GeoLocationCmd = redisRaw.GeoLocationCmd
type StatusCmd = redisRaw.StatusCmd
type XStreamSliceCmd = redisRaw.XStreamSliceCmd
type XReadGroupArgs = redisRaw.XReadGroupArgs
type BoolCmd = redisRaw.BoolCmd
type XReadArgs = redisRaw.XReadArgs
type StringStringMapCmd = redisRaw.StringStringMapCmd

const Nil = redisRaw.Nil

type RedisClient interface {
	XAdd(background context.Context, record *XAddArgs) *StringCmd
	XGroupCreateMkStream(ctx context.Context, stream, group, start string) *StatusCmd
	Get(ctx context.Context, key string) *StringCmd
	XAck(ctx context.Context, stream, group string, ids ...string) *IntCmd
	SIsMember(ctx context.Context, key string, member interface{}) *BoolCmd
	XRead(ctx context.Context, a *XReadArgs) *XStreamSliceCmd
	FlushAll(ctx context.Context) *StatusCmd
	SAdd(ctx context.Context, key string, members ...interface{}) *IntCmd
	SRem(ctx context.Context, key string, members ...interface{}) *IntCmd
	Del(ctx context.Context, keys ...string) *IntCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd
	Ping(ctx context.Context) *StatusCmd
	SMembers(ctx context.Context, key string) *StringSliceCmd
	HSet(ctx context.Context, key string, values ...interface{}) *IntCmd
	HDel(ctx context.Context, key string, fields ...string) *IntCmd
	HGetAll(ctx context.Context, key string) *StringStringMapCmd
	HGet(ctx context.Context, key, field string) *StringCmd
	HExists(ctx context.Context, key, field string) *BoolCmd
	XReadGroup(ctx context.Context, a *XReadGroupArgs) *XStreamSliceCmd
	Subscribe(ctx context.Context, channels ...string) *PubSub
	PSubscribe(ctx context.Context, channels ...string) *PubSub
	Context() context.Context
	Do(ctx context.Context, args ...interface{}) *Cmd
	Process(ctx context.Context, cmd Cmder) error
	PoolStats() *PoolStats
	Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)
	Pipeline() Pipeliner
	TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)
	TxPipeline() Pipeliner
	Close() error
}

var _ RedisClient = (*Client)(nil)

type Client struct {
	redisClient RedisClient
}

func (this *Client) HSet(ctx context.Context, key string, values ...interface{}) *IntCmd {
	return this.redisClient.HSet(ctx, key, values...)
}

func (this *Client) HDel(ctx context.Context, key string, fields ...string) *IntCmd {
	return this.redisClient.HDel(ctx, key, fields...)
}

func (this *Client) HGetAll(ctx context.Context, key string) *StringStringMapCmd {
	return this.redisClient.HGetAll(ctx, key)
}

func (this *Client) HGet(ctx context.Context, key, field string) *StringCmd {
	return this.redisClient.HGet(ctx, key, field)
}

func (this *Client) HExists(ctx context.Context, key, field string) *BoolCmd {
	return this.redisClient.HExists(ctx, key, field)
}

func (this *Client) SMembers(ctx context.Context, key string) *StringSliceCmd {
	return this.redisClient.SMembers(ctx, key)
}

func (this *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd {
	return this.redisClient.Set(ctx, key, value, expiration)
}

func (this *Client) Ping(ctx context.Context) *StatusCmd {
	return this.redisClient.Ping(ctx)
}

func (this *Client) Del(ctx context.Context, keys ...string) *IntCmd {
	return this.redisClient.Del(ctx, keys...)
}
func (this *Client) SRem(ctx context.Context, key string, members ...interface{}) *IntCmd {
	return this.redisClient.SRem(ctx, key, members...)
}
func (this *Client) SAdd(ctx context.Context, key string, members ...interface{}) *IntCmd {
	return this.redisClient.SAdd(ctx, key, members...)
}
func (this *Client) FlushAll(ctx context.Context) *StatusCmd {
	return this.redisClient.FlushAll(ctx)
}
func (this *Client) XRead(ctx context.Context, a *XReadArgs) *XStreamSliceCmd {
	return this.redisClient.XRead(ctx, a)
}
func (this *Client) XAck(ctx context.Context, stream, group string, ids ...string) *IntCmd {
	return this.redisClient.XAck(ctx, stream, group, ids...)
}
func (this *Client) SIsMember(ctx context.Context, key string, member interface{}) *BoolCmd {
	return this.redisClient.SIsMember(ctx, key, member)
}
func (this *Client) Get(ctx context.Context, key string) *StringCmd {
	return this.redisClient.Get(ctx, key)
}
func (this *Client) XReadGroup(ctx context.Context, a *XReadGroupArgs) *XStreamSliceCmd {
	return this.redisClient.XReadGroup(ctx, a)
}

// Env reads specified environment variable. If no value has been found,
// fallback is returned.
func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func NewClient(opt *Options) *Client {
	clusterOrHost := Env("REDIS_CLUSTER", "true")
	if strings.Compare("false", strings.ToLower(clusterOrHost)) == 0 {
		return &Client{redisRaw.NewClient(opt)}
	} else {
		return &Client{redisRaw.NewClusterClient(convert2ClusterOptions(opt))}

	}
}

func convert2ClusterOptions(opt *Options) *redisRaw.ClusterOptions {
	return &redisRaw.ClusterOptions{
		Addrs:              []string{opt.Addr},
		Dialer:             opt.Dialer,
		OnConnect:          opt.OnConnect,
		Username:           opt.Username,
		Password:           opt.Password,
		MaxRetries:         opt.MaxRetries,
		MinRetryBackoff:    opt.MinRetryBackoff,
		MaxRetryBackoff:    opt.MaxRetryBackoff,
		DialTimeout:        opt.DialTimeout,
		ReadTimeout:        opt.ReadTimeout,
		WriteTimeout:       opt.WriteTimeout,
		PoolFIFO:           opt.PoolFIFO,
		PoolSize:           opt.PoolSize,
		MinIdleConns:       opt.MinIdleConns,
		MaxConnAge:         opt.MaxConnAge,
		PoolTimeout:        opt.PoolTimeout,
		IdleTimeout:        opt.IdleTimeout,
		IdleCheckFrequency: opt.IdleCheckFrequency,
		TLSConfig:          opt.TLSConfig,
	}
}
func (this *Client) Context() context.Context {
	return this.redisClient.Context()
}
func (this *Client) Do(ctx context.Context, args ...interface{}) *Cmd {
	return this.redisClient.Do(ctx, args)
}
func (this *Client) Process(ctx context.Context, cmd Cmder) error {
	return this.redisClient.Process(ctx, cmd)
}
func (this *Client) PoolStats() *PoolStats {
	return this.redisClient.PoolStats()
}
func (this *Client) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return this.redisClient.Pipelined(ctx, fn)
}
func (this *Client) Pipeline() Pipeliner {
	return this.redisClient.Pipeline()
}
func (this *Client) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return this.redisClient.TxPipelined(ctx, fn)
}
func (this *Client) TxPipeline() Pipeliner {
	return this.redisClient.TxPipeline()
}
func (this *Client) Subscribe(ctx context.Context, channels ...string) *PubSub {
	return this.redisClient.Subscribe(ctx, channels...)
}
func (this *Client) PSubscribe(ctx context.Context, channels ...string) *PubSub {
	return this.redisClient.PSubscribe(ctx, channels...)
}
func (this *Client) XAdd(background context.Context, record *XAddArgs) *StringCmd {
	return this.redisClient.XAdd(background, record)
}
func (this *Client) Close() error {
	return this.redisClient.Close()
}
func (this *Client) XGroupCreateMkStream(ctx context.Context, stream, group, start string) *StatusCmd {
	return this.redisClient.XGroupCreateMkStream(ctx, stream, group, start)
}
