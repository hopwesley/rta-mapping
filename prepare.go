package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hopwesley/rta-mapping/common"
	"strconv"
	"strings"
	"time"
)

const (
	IDTableName         = "ad_id_mapping"
	DefaultMaxIDMapSize = 20_000_000
	RatRedisKeyPrefix   = "ad:bytedance:"
	ReadSizeOnce        = 1 << 10
)

type RedisCfg struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
}

func (c *RedisCfg) String() string {
	s := "\n========redis config========"
	s += "\nAddress:" + c.Addr
	s += "\n============================"
	return s
}

func InitRtaMap(cfg *RedisCfg) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,     //"localhost:6379", // Redis地址
		Password: cfg.Password, // 密码（如果有的话）
	})
	ctx := context.Background()
	defer rdb.Close()
	var start = time.Now()

	keys, err := rdb.Keys(ctx, RatRedisKeyPrefix+"*").Result()
	if err != nil {
		fmt.Println("redis read keys err:", err)
		return err
	}
	if len(keys) == 0 {
		return fmt.Errorf("no data in redis")
	}
	fmt.Printf("Found %d keys in Redis:\n", len(keys))

	for _, key := range keys {

		card, err := rdb.SCard(ctx, key).Result()
		if err != nil {
			fmt.Printf("failed to get cardinality for key %s: %v\n", key, err)
			return err
		}
		fmt.Printf("key: %s, cardinality: %d\n", key, card)
		if card == 0 {
			continue
		}

		rtaIDStr, found := strings.CutPrefix(key, RatRedisKeyPrefix)
		if !found {
			return fmt.Errorf("invalid rta id")
		}
		rid, err := strconv.ParseInt(rtaIDStr, 10, 64)
		if err != nil {
			return err
		}

		var cursor uint64
		var counter = 0
		for {
			userIDStr, nextCursor, err := rdb.SScan(ctx, key, cursor, "", ReadSizeOnce).Result()
			if err != nil {
				fmt.Printf("failed to scan members for key %s: %v\n", key, err)
				return err
			}

			var userIDs []int
			for _, uidStr := range userIDStr {
				uid, err := strconv.Atoi(uidStr)
				if err != nil {
					fmt.Println("invalid user id found:", uidStr)
					return err
				}
				userIDs = append(userIDs, uid)
			}
			counter += len(userIDStr)
			common.RtaMapInst().InitByOneRtaWithoutLock(rid, userIDs)

			cursor = nextCursor
			fmt.Printf("\rrtaid[%s] loaded: %d progress: %.2f%% time used: %s",
				rtaIDStr, counter, float32(counter)*100/float32(card), time.Since(start))

			if cursor == 0 {
				fmt.Printf("\nkey:%s read finished!\n", key)
				break
			}
		}
	}

	fmt.Println("Successfully init RtaBitMap")
	return nil
}

type MysqlCfg struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Limit    int64  `json:"limit"`
}

func (c *MysqlCfg) String() string {
	s := "\n============mysql config========"
	s += "\nuser name:" + c.UserName
	s += "\nhost:" + c.Host
	s += "\nport:" + c.Port
	s += "\ndatabase:" + c.Database
	s += fmt.Sprintf("\nlimit:%d", c.Limit)
	s += "\n================================"
	return s

}

func (c *MysqlCfg) ToDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.UserName, c.Password, c.Host, c.Port, c.Database)
}
func (c *MysqlCfg) OpenDatabase() (*sql.DB, error) {
	if c.Limit <= 0 {
		c.Limit = DefaultMaxIDMapSize
	}
	dsn := c.ToDsn() //"username:password@tcp(localhost:3306)/yourdatabase"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("open mysql database failed:", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("connect to database failed:", err)
		return nil, err
	}
	fmt.Println("Successfully connected to the database")
	return db, nil
}

func InitIDMap(cfg *MysqlCfg) error {
	db, err := cfg.OpenDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT uid, imei_md5, oaid, idfa, android_id FROM %s LIMIT ?", IDTableName)
	rows, err := db.Query(query, cfg.Limit)

	if err != nil {
		fmt.Println("data query failed:", err)
		return err
	}

	var counter = int64(0)
	var start = time.Now()
	for rows.Next() {
		var item common.JsonRequest

		err := rows.Scan(&item.UserID, &item.IMEIMD5, &item.OAID, &item.IDFA, &item.AndroidIDMD5)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return err
		}

		counter++
		common.IDMapInst().UpdateByMySqlWithoutLock(item)
		if counter%10000 == 0 {
			fmt.Printf("\rUpdated id map size: %d progress: %.2f%% time used: %s", counter, float32(counter)*100/float32(cfg.Limit), time.Since(start))
			//bts, _ := json.Marshal(item)
			//fmt.Println(string(bts))
		}
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Error scanning rows:", err)
		return err
	}
	fmt.Println("Successfully init id map")
	return nil
}
