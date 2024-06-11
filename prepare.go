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
)

const (
	IDTableName = "ID_MAPPING"
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

	rows, err := db.Query("SELECT uid,imei_md5,oaid,idfa,android_id FROM " + IDTableName)
	if err != nil {
		return err
	}
	common.IDMapInst().CleanMap()
	for rows.Next() {
		var item common.IDUpdateRequest
		err := rows.Scan(&item.UserID, &item.IMEIMD5, &item.OAID, &item.IDFA, &item.AndroidIDMD5)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return err
		}

		common.IDMapInst().UpdateByMySqlWithoutLock(item)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Error scanning rows:", err)
		return err
	}

	return nil
}
