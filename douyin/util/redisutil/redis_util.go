// Package redisutil
// @Author shaofan
// @Date 2022/5/13
// @DESC redis连接初始化与工具
package redisutil

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

// Set 							插入string类型
// key 							键
// value 						需要插入的值，内部需要进行序列化
func Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = RedisDB.Set(key, data, 0).Err()
	return err
}

// SetWithExpireTime 			插入string类型并设置过期时间
// key 							键
// value 						需要插入的值，内部需要进行序列化
// duration						过期时间
func SetWithExpireTime(key string, value interface{}, duration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = RedisDB.Set(key, data, duration).Err()
	return err
}

// Get 							获取string类型
// key 							键
// value 						获取的值存储的指针
func Get[T any](key string, value *T) error {
	v, err := RedisDB.Get(key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(v), value)
	if err != nil {
		return err
	}
	return nil
}

// GetAndDelete					获取string并删除string类型
// key							键
// value						值
func GetAndDelete[T any](key string, value *T, isTx bool, p redis.Pipeliner) error {
	// Get
	v, err := RedisDB.Get(key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(v), value)
	if err != nil {
		return err
	}
	if isTx {
		// Delete
		_, err = p.Del(key).Result()
		if err == redis.Nil {
			return nil
		} else if err != nil {
			return err
		}
	} else {
		// Delete
		_, err = RedisDB.Del(key).Result()
		if err == redis.Nil {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ZAdd 						插入set类型
// key						 	键
// value 						需要插入的值
func ZAdd(key string, value []redis.Z, isTx bool, pipeline redis.Pipeliner) error {
	if isTx {
		_, err := pipeline.ZAdd(key, value...).Result()
		if err != nil {
			return err
		}
	} else {
		_, err := RedisDB.ZAdd(key, value...).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

// ZRevRange 					逆序获取set类型
// key 							键
// value 						获取的值存储的指针
func ZRevRange[T any](key string, value *[]T) error {
	val, err := RedisDB.ZRevRange(key, 0, -1).Result()
	if *value == nil || len(*value) < len(val) {
		*value = make([]T, len(val))
	}
	for i, v := range val {
		err := json.Unmarshal([]byte(v), &(*value)[i])
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

// ZAddWithExpireTime 			插入set类型并设置过期时间
// key							键
// value						需要插入的值
// duration						时间
func ZAddWithExpireTime(key string, value []redis.Z, duration time.Duration, isTx bool, pipeline redis.Pipeliner) error {
	// 如果处于外部事务内，使用外部事务进行处理
	if isTx {
		_, err := pipeline.ZAdd(key, value...).Result()
		if err != nil {
			return nil
		}
		_, err = pipeline.Expire(key, duration).Result()
		if err != nil {
			return err
		}
	} else {
		//如果不处于外部事务，内部使用事务对插入和过期时间两个指令进行事务处理
		pipeline = RedisDB.TxPipeline()
		_, err := pipeline.ZAdd(key, value...).Result()
		if err != nil {
			return err
		}
		_, err = pipeline.Expire(key, duration).Result()
		if err != nil {
			err1 := pipeline.Discard()
			if err1 != nil {
				return err1
			}
			return err
		}
		_, err = pipeline.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

// ZRem 						删除sort_set中的数据
// key							需要删除数据的键
// value						需要删除的数据
func ZRem[T any](key string, value *[]T, isTx bool, pipeline redis.Pipeliner) error {
	for _, v := range *value {
		val, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if isTx {
			_, err := pipeline.ZRem(key, val).Result()
			if err != nil {
				return err
			}
		} else {
			_, err := RedisDB.ZRem(key, val).Result()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Keys 获取匹配的键集
func Keys(prefix string, keys *[]string) error {
	res := RedisDB.Keys(prefix + "*")
	*keys = res.Val()
	return res.Err()
}

// TTL 获取的键值的过期时间
func TTL(key string) (time.Duration, error) {
	res := RedisDB.TTL(key)
	return res.Val(), res.Err()
}

// Begin 开启事务
func Begin() redis.Pipeliner {
	return RedisDB.TxPipeline()
}

// Expire 	设置键的过期时间
// key		键
// duration	过期时间
func Expire(key string, duration time.Duration) error {
	_, err := RedisDB.Expire(key, duration).Result()
	if err != nil {
		return err
	}
	return nil
}
