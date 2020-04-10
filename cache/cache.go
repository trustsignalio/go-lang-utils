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
	expiration int32
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

// NewMultiClient method will return a pointer to MultiClient object
func NewMultiClient(prefix, mcServer string, defCacheTime int) *MultiClient {
	var cacheTime = time.Duration(defCacheTime) * time.Minute
	c := cache.New(cacheTime, 10*time.Minute)
	mc := memcache.New(mcServer)

	var cc = &MultiClient{client: c, mc: mc, prefix: prefix, expiration: int32(defCacheTime * 36)}
	return cc
}

func (cc *MultiClient) getKeyName(key string) string {
	return cc.prefix + "_" + key
}

// Set method will set the object in both memory cache and memcache
func (cc *MultiClient) Set(key string, val interface{}) {
	k := cc.getKeyName(key)
	cc.client.SetDefault(k, val)

	result, err := json.Marshal(val)
	if err == nil {
		cc.mc.Set(&memcache.Item{
			Key:        k,
			Value:      result,
			Expiration: cc.expiration,
		})
	}
}

// SetInMemory method will set the object in memory cache
func (cc *MultiClient) SetInMemory(key string, val interface{}) {
	k := cc.getKeyName(key)
	cc.client.Set(k, val, time.Duration(cc.expiration)*time.Second)
}

// SetWithExpire method will set the object in both memory cache and memcache
func (cc *MultiClient) SetWithExpire(key string, val interface{}, secs int) {
	k := cc.getKeyName(key)
	cc.client.Set(k, val, time.Duration(secs)*time.Second)

	result, err := json.Marshal(val)
	if err == nil {
		cc.mc.Set(&memcache.Item{
			Key:        k,
			Value:      result,
			Expiration: int32(secs),
		})
	}
}

// Get method tires to find the key from memory cache then check memcache
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

// GetWithSet method tries to get the key from program memory cache and if
// it fails then tries memcache and if the item is found in memcache then it
// is set in program memory for faster lookup
func (cc *MultiClient) GetWithSet(key string, resultObj interface{}) (interface{}, bool) {
	k := cc.getKeyName(key)
	result, found := cc.client.Get(k)
	if found {
		return result, found
	}
	item, err := cc.mc.Get(k)
	if err == nil {
		err = json.Unmarshal(item.Value, resultObj)
		if err == nil {
			cc.client.Set(k, resultObj, time.Duration(cc.expiration)*time.Second)
			return resultObj, true
		}
	}

	return nil, false
}

// GetSliceWithSet method tries to get the key from program memory cache and if
// it fails then tries memcache and if the item is found in memcache then it
// is set in program memory for faster lookup
func (cc *MultiClient) GetSliceWithSet(key string, resultObj interface{}) (interface{}, bool) {
	k := cc.getKeyName(key)
	result, found := cc.client.Get(k)
	if found {
		return result, found
	}
	item, err := cc.mc.Get(k)
	if err == nil {
		return item.Value, true
	}

	return nil, false
}

// Delete method will remove the key from both memory cache and memcache
func (cc *MultiClient) Delete(key string) {
	k := cc.getKeyName(key)
	cc.client.Delete(k)
	cc.mc.Delete(k)
}
