package utils

import (
	"bytes"
	"encoding/gob"
	"errors"
	"reflect"
	"time"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"github.com/FZambia/go-sentinel"
)

var BillMasterRedisPool *redis.Pool

var (
	ErrCacheMiss = errors.New("CacheMiss")
)

func InitCache() bool {
	BillMasterRedisPool = newSentinelPool()
	return true
}

func newSentinelPool() *redis.Pool {
	sntnl := &sentinel.Sentinel{
		Addrs:      []string{"10.16.25.121:26379", "10.16.25.121:26379", "10.16.25.121:26379"},
		MasterName: "master-dev",
		Dial: func(addr string) (redis.Conn, error) {
			timeout := 500 * time.Millisecond
			c, err := redis.DialTimeout("tcp", addr, timeout, timeout, timeout)

			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	add, _ := sntnl.MasterAddr()
	return redisPool(add, 0, "mei@1q2w3d", 1000, 1000, 1000, 1000, 50)
	/*return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   64,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				return nil, err
			}
			c, err := redis.Dial("tcp", masterAddr)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if !sentinel.TestRole(c, "master") {
				return errors.New("Role check failed")
			} else {
				return nil
			}
		},
	}*/
}
func redisPool(addr string, db int, pass string, connectTimeoutMs, readTimeoutMs, writeTimeoutMs, idleTimeoutSec int64, maxIdle int) *redis.Pool {

	beego.Info("addr:%s,db:%d,pass:%s,connectTimeoutMs:%d,readTimeoutMs:%d,writeTimeoutMs:%d,idleTimeoutSec:%d,maxid:%d",
		addr, db, pass, time.Millisecond*time.Duration(connectTimeoutMs),
		time.Millisecond*time.Duration(readTimeoutMs), time.Millisecond*time.Duration(writeTimeoutMs),
		time.Duration(idleTimeoutSec)*time.Second, maxIdle)
	return &redis.Pool{
		MaxActive:   maxIdle,
		MaxIdle:     maxIdle,
		IdleTimeout: time.Duration(idleTimeoutSec) * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", addr, time.Millisecond*time.Duration(connectTimeoutMs),
				time.Millisecond*time.Duration(readTimeoutMs), time.Millisecond*time.Duration(writeTimeoutMs))
			if err != nil {
				beego.Error("redis connect %s fail", addr)
				return nil, err
			}
			if pass != "" {
				_, err = c.Do("AUTH", pass)
				if err != nil {
					beego.Error("redis password error: %s.", err.Error())
					return nil, err
				}
			}

			err = c.Send("SELECT", db)
			if err != nil {
				beego.Error("redis select database %d fail, error: %s", db, err.Error())
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			/*_, err := c.Do("PING")
			return err*/
			if !sentinel.TestRole(c, "master") {
				return errors.New("Role check failed")
			} else {
				return nil
			}
		},
	}
}

func testRedisPool(pool *redis.Pool, name string) bool {
	c := pool.Get()

	if c != nil {
		defer c.Close()
	}

	err := pool.TestOnBorrow(c, time.Now())
	if err == nil {
		//l4g.Trace("redis[%s] connect success.", name)
		return true
	} else {
		//l4g.Error("redis[%s] connect fail.", name)
		return false
	}
	return false
}

func SetCache(key string, value interface{}, timeout int) error {
	data, err := Encode(value)
	if err != nil {
		return err
	}

	ttl := timeout
	conn := BillMasterRedisPool.Get()
	if conn != nil {
		defer conn.Close()
	}

	_, err = conn.Do("PING")
	if err != nil {
		beego.Info("BillMasterRedisPool PING failed,err:%s", err)
		err = errors.New("PING failed")
		return err
	}

	if _, err = conn.Do("SET", key, data); err != nil {
		beego.Info("redis err = %s,key = %s,data:%s", err, key, data)
		err = errors.New("Set failed")
		return err
	}

	// ttl>0的才设置过期时间
	if ttl > 0 {
		if _, err = conn.Do("expire", key, ttl); err != nil {
			beego.Info("expire key=%s,err=%s", key, err)
			err = errors.New("Expire failed")
			return err
		}
	}

	return err
}

func ExpireCache(key string, timeout int) error {

	ttl := timeout
	conn := BillMasterRedisPool.Get()
	if conn != nil {
		defer conn.Close()
	}

	_, err := conn.Do("PING")
	if err != nil {
		beego.Info("BillMasterRedisPool PING failed,err:%s", err)
		err = errors.New("PING failed")
		return err
	}

	if _, err = conn.Do("expire", key, ttl); err != nil {
		beego.Info("expire key=%s,err=%s", key, err)
		err = errors.New("Expire failed")
		return err
	}

	return err
}

func GetCache(key string, to interface{}) error {
	conn := BillMasterRedisPool.Get()
	if conn != nil {
		defer conn.Close()
	}

	_, err := conn.Do("PING")
	if err != nil {
		beego.Info("BillMasterRedisPool PING failed,err:%s", err)
		err = errors.New("PING failed")
		return err
	}

	//key := fmt.Sprintf("%s:%d", keyPrefix, mobile)
	if data, err := conn.Do("get", key); err != nil {
		beego.Info("get err:%s,key = %s", err, key)
		err = errors.New("Get failed")
		return err
	} else if reflect.TypeOf(data) == nil {
		//key not found
		beego.Info("Key not found get err:%s,key = %s", err, key)
		err = errors.New("Key not found failed")
		return err
	} else if err = Decode(data.([]byte), to); err != nil {
		return err
	}

	return err
}

func DelCache(key string) error {
	conn := BillMasterRedisPool.Get()
	if conn != nil {
		defer conn.Close()
	}

	if _, err := conn.Do("PING"); err != nil {
		beego.Info("BillMasterRedisPool PING failed,err:%s", err)
		return errors.New("PING failed")
	}

	if _, err := conn.Do("DEL", key); err != nil {
		beego.Info("del key=%s err=%s", key, err)
		return errors.New("Del key failed")
	}

	return nil
}

func IsExistCache(key string) (bool, error) {
	conn := BillMasterRedisPool.Get()
	if conn != nil {
		defer conn.Close()
	}

	if _, err := conn.Do("PING"); err != nil {
		beego.Info("BillMasterRedisPool PING failed,err:%s", err)
		return false, errors.New("PING failed")
	}

	if b, err := redis.Bool(conn.Do("EXISTS", key)); err != nil {
		beego.Info("IsExist key=%s err=%s", key, err)
		return false, errors.New("IsExist key failed")
	} else {
		return b, nil
	}
}

func IncrCache(key string) error {
	conn := BillMasterRedisPool.Get()
	if conn != nil {
		defer conn.Close()
	}

	if _, err := conn.Do("PING"); err != nil {
		beego.Info("BillMasterRedisPool PING failed,err:%s", err)
		return errors.New("PING failed")
	}

	if _, err := redis.Bool(conn.Do("INCRBY", key, 1)); err != nil {
		beego.Info("Incr key=%s err=%s", key, err)
		return errors.New("Incr key failed")
	}

	return nil
}

func DecrCache(key string) (bool, error) {
	conn := BillMasterRedisPool.Get()
	if conn != nil {
		defer conn.Close()
	}

	if _, err := conn.Do("PING"); err != nil {
		beego.Info("BillMasterRedisPool PING failed,err:%s", err)
		return false, errors.New("PING failed")
	}

	if b, err := redis.Bool(conn.Do("INCRBY", key, -1)); err != nil {
		beego.Info("Decr key=%s err=%s", key, err)
		return false, err
	} else {
		return b, nil
	}
}

func IncrValue(key string) (int, error) {
	conn := BillMasterRedisPool.Get()
	if conn != nil {
		defer conn.Close()
	}

	if _, err := conn.Do("PING"); err != nil {
		beego.Info("BillMasterRedisPool PING failed,err:%s", err)
		return 0, errors.New("PING failed")
	}

	if v, err := redis.Int(conn.Do("get", key)); err != nil {
		beego.Info("get err:%s,key = %s", err, key)
		return 0, errors.New("Get incr value failed")
	} else {
		return v, nil
	}
}

func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
