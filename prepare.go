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
		rtaIDs, cursor, err = rdb.Scan(ctx, cursor, "", 0).Result()
		if err != nil {
			return err
		}

		for _, rtaID := range rtaIDs {
			jsonUserIDs, err := rdb.Get(ctx, rtaID).Result()
			if err != nil {
				return err
			}
			var userIDs []int
			err = json.Unmarshal([]byte(jsonUserIDs), &userIDs)
			if err != nil {
				fmt.Println("Error unmarshalling array:", err)
				return err
			}

			fmt.Printf("rtaID: %s, Value Len: %d\n", rtaID, len(userIDs))
			if len(userIDs) == 0 {
				continue
			}
			rid, err := strconv.ParseInt(rtaID, 10, 64)
			if err != nil {
				return err
			}
			err = common.RtaMapInst().InitByOneRtaWithoutLock(rid, userIDs)
			if err != nil {
				return err
			}

			fmt.Printf("init rta[%s] successfully\n", rtaID)
		}

		if cursor == 0 {
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

func InitIDMap(cfg *MysqlCfg) error {

	dsn := cfg.ToDsn() //"username:password@tcp(localhost:3306)/yourdatabase"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("open mysql database failed:", err)
		return err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("connect to database failed:", err)
		return err
	}
	fmt.Println("Successfully connected to the database")
	if cfg.Limit <= 0 {
		cfg.Limit = DefaultMaxIDMapSize
	}
	query := fmt.Sprintf("SELECT uid, imei_md5, oaid, idfa, android_id FROM %s LIMIT ?", IDTableName)
	rows, err := db.Query(query, cfg.Limit)

	if err != nil {
		fmt.Println("data query failed:", err)
		return err
	}
	common.IDMapInst().CleanMap()
	var counter = int64(0)
	var start = time.Now()
	for rows.Next() {
		var item common.IDUpdateRequest

		err := rows.Scan(&item.UserID, &item.IMEIMD5, &item.OAID, &item.IDFA, &item.AndroidIDMD5)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return err
		}
		//bts, _ := json.Marshal(item)
		//fmt.Println(string(bts))
		counter++
		common.IDMapInst().UpdateByMySqlWithoutLock(item)
		if counter%1000 == 0 {
			fmt.Printf("\rUpdated id map size: %d progress: %.2f%% time used: %s", counter, float32(counter)*100/float32(cfg.Limit), time.Since(start))
		}
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Error scanning rows:", err)
		return err
	}
	fmt.Println("Successfully init id map")
	return nil
}
