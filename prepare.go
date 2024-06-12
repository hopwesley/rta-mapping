package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hopwesley/rta-mapping/common"
	"strconv"
	"time"
)

const (
	IDTableName         = "ad_id_mapping"
	DefaultMaxIDMapSize = 20_000_000
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
	defer rdb.Close()

	var cursor uint64
	var rtaIDs []string
	for {
		var err error
		keys, err := rdb.Keys(ctx, "ad:bytedance*").Result()
		//rtaIDs, cursor, err = rdb.Scan(ctx, cursor, "", 0).Result()
		if err != nil {
			fmt.Println("redis scan err:", err)
			return err
		}
		fmt.Println("rtaIDS len:", len(rtaIDs))
		for _, key := range keys {
			jsonUserIDs, err := rdb.Get(ctx, key).Result()
			if err != nil {
				fmt.Println("redis get err:", err)
				return err
			}
			var userIDs []int
			err = json.Unmarshal([]byte(jsonUserIDs), &userIDs)
			if err != nil {
				fmt.Println("Error unmarshalling array:", err)
				return err
			}

			fmt.Printf("rtaID: %s, Value Len: %d\n", key, len(userIDs))
			if len(userIDs) == 0 {
				continue
			}
			rid, err := strconv.ParseInt(key, 10, 64)
			if err != nil {
				return err
			}
			common.RtaMapInst().InitByOneRtaWithoutLock(rid, userIDs)

			fmt.Printf("init rta[%s] successfully\n", key)
		}

		if cursor == 0 {
			fmt.Println("cursor is zero now")
			break
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
