package cache

import "github.com/redis/go-redis/v9"

type Cache struct {
	Addr     string `config:"cache_addr"`
	Password string `config:"cache_password"`
}

func (c *Cache) Open() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password, // no password set
		DB:       0,          // use default DB
	})

	return rdb
}
