package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Client struct {
	prefix string
	client *cache.Cache
}

func NewClient(prefix string, defCacheTime int) *Client {
	var cacheTime = time.Duration(defCacheTime) * time.Minute
	c := cache.New(cacheTime, 10*time.Minute)

	var cc = &Client{client: c, prefix: prefix}
	return cc
}

func (cc *Client) getKeyName(key string) string {
	return cc.prefix + "_" + key
}

func (cc *Client) Set(key string, val interface{}) {
	cc.client.SetDefault(cc.getKeyName(key), val)
}

func (cc *Client) Get(key string) (interface{}, bool) {
	return cc.client.Get(cc.getKeyName(key))
}

func (cc *Client) Delete(key string) {
	cc.client.Delete(key)
}
