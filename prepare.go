package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hopwesley/rta-mapping/common"
	"strconv"
)

type RedisCfg struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
}

func InitRtaMap(cfg *RedisCfg) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,     //"localhost:6379", // Redis地址
		Password: cfg.Password, // 密码（如果有的话）
		DB:       0,            // 使用默认DB
	})
	ctx := context.Background()

	var cursor uint64
	var keys []string
	for {
		var err error
		keys, cursor, err = rdb.Scan(ctx, cursor, "", 0).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			jsonNumbers, err := rdb.Get(ctx, key).Result()
			if err != nil {
				return err
			}
			var numbers []int
			err = json.Unmarshal([]byte(jsonNumbers), &numbers)
			if err != nil {
				fmt.Println("Error unmarshaling array:", err)
				return err
			}

			fmt.Printf("Key: %s, Value: %d\n", key, len(numbers))

			int64Key, err := strconv.ParseInt(key, 10, 64)
			if err != nil {
				fmt.Printf("Error converting key %s to int64: %v\n", key, err)
				return err
			}
			prepareRtaMap(int64Key, numbers)
		}

		if cursor == 0 {
			break
		}
	}

	fmt.Println("Successfully init RtaMap")
	return nil
}

func prepareRtaMap(rtaId int64, numbers []int) {
	common.RtaMapInst().InitByOneRta(rtaId, numbers)
}

func InitIDMap() {

}
