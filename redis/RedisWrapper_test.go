package redis

import "testing"

func TestRedisCLientWrapper4SingleHost(test *testing.T)  {
	var redisClient *Client
	redisClient = NewClient(&Options{})
	redisClient.PoolStats()
}
