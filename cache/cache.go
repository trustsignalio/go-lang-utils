package cache

import (
	"encoding/json"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/patrickmn/go-cache"
)

type Client struct {
	prefix string
	client *cache.Cache
}

type MultiClient struct {
	prefix     string
	expiration int
	client     *cache.Cache
	mc         *memcache.Client
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
	cc.client.Delete(cc.getKeyName(key))
}

func NewMultiClient(prefix, mcServer string, defCacheTime int) *MultiClient {
	var cacheTime = time.Duration(defCacheTime) * time.Minute
	c := cache.New(cacheTime, 10*time.Minute)
	mc := memcache.New(mcServer)

	var cc = &MultiClient{client: c, mc: mc, prefix: prefix, expiration: defCacheTime * 60}
	return cc
}

func (cc *MultiClient) getKeyName(key string) string {
	return cc.prefix + "_" + key
}

func (cc *MultiClient) Set(key string, val interface{}) {
	k := cc.getKeyName(key)
	cc.client.SetDefault(k, val)

	result, err := json.Marshal(val)
	if err == nil {
		cc.mc.Set(&memcache.Item{
			Key:        k,
			Value:      result,
			Expiration: int32(cc.expiration),
		})
	}
}

func (cc *MultiClient) Get(key string) (interface{}, bool) {
	k := cc.getKeyName(key)
	result, found := cc.client.Get(k)
	if found {
		return result, found
	}
	item, err := cc.mc.Get(k)
	if err == nil {
		var cacheObj interface{}
		err = json.Unmarshal(item.Value, cacheObj)
		if err == nil {
			return cacheObj, true
		}
	}

	return nil, false
}

func (cc *MultiClient) Delete(key string) {
	k := cc.getKeyName(key)
	cc.client.Delete(k)
	cc.mc.Delete(k)
}
