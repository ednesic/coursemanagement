package cache

import (
	"sync"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

var (
	instance Cache
	once     sync.Once
)

//Cache is an interface to handle cache
type Cache interface {
	Get(string, interface{}) error
	Set(string, interface{}, time.Duration) error
	Delete(string) error
	Initialize(map[string]string)
	Disconnect()
}

type rImpl struct {
	codec *cache.Codec
	ring  *redis.Ring
}

//GetInstance to return a redis client
func GetInstance() Cache {
	once.Do(func() {
		if instance == nil {
			instance = &rImpl{}
		}
	})
	return instance
}

func (rc *rImpl) Initialize(hosts map[string]string) {
	rc.ring = redis.NewRing(&redis.RingOptions{
		Addrs: hosts,
	})
	rc.codec = &cache.Codec{
		Redis: rc.ring,

		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

func (rc rImpl) Get(key string, object interface{}) error {
	if err := rc.codec.Get(key, object); err != nil {
		return &RedisErr{Msg: err.Error()}
	}
	return nil
}

func (rc rImpl) Set(k string, obj interface{}, d time.Duration) error {
	if err := rc.codec.Set(&cache.Item{Key: k, Object: obj, Expiration: d}); err != nil {
		return &RedisErr{Msg: err.Error()}
	}
	return nil
}

func (rc rImpl) Delete(key string) error {
	if err := rc.codec.Delete(key); err != nil {
		return &RedisErr{Msg: err.Error()}
	}
	return nil
}

func (rc rImpl) Disconnect() {
	_ = rc.ring.Close()
}
