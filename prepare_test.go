package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hopwesley/rta-mapping/common"
	"strconv"
	"strings"
	"testing"
)

var password string

func init() {
	flag.StringVar(&password, "password", "", "password")
}

func TestInitIDMap(t *testing.T) {
	cfg := &MysqlCfg{
		UserName: "marketing",
		Host:     "rm-bp177891b06z9yj8dfo.mysql.rds.aliyuncs.com",
		Port:     "3306",
		Database: "marketing",
		Limit:    20_000_000,
	}
	cfg.Password = password
	err := InitIDMap(cfg)
	if err != nil {
		t.Fatal(err)
	}
	common.PrintMemUsage()
}

func TestMysqlPing(t *testing.T) {

	dsn := "marketing:" + password + "@tcp(rm-bp177891b06z9yj8dfo.mysql.rds.aliyuncs.com:3306)/marketing"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("open mysql database failed:", err)
		t.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("connect to database failed:", err)
		t.Fatal(err)
	}
	fmt.Println("Successfully connected to the database")

	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", IDTableName)
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Total rows in %s: %d\n", IDTableName, count)
}

func TestMysqlQuerySomeData(t *testing.T) {
	dsn := "marketing:" + password + "@tcp(rm-bp177891b06z9yj8dfo.mysql.rds.aliyuncs.com:3306)/marketing"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("open mysql database failed:", err)
		t.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT uid, imei_md5, oaid, idfa, android_id FROM %s LIMIT ?", IDTableName)
	rows, err := db.Query(query, 10)

	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var item common.JsonRequest

		err := rows.Scan(&item.UserID, &item.IMEIMD5, &item.OAID, &item.IDFA, &item.AndroidIDMD5)
		if err != nil {
			t.Fatal(err)
		}
		bts, _ := json.Marshal(item)
		fmt.Println(string(bts))
	}
}

func TestInitRtaMap(t *testing.T) {
	cfg := &RedisCfg{
		Addr:     "47.99.198.186:6600",
		Password: password,
	}
	err := InitRtaMap(cfg)
	if err != nil {
		t.Fatal(err)
	}
	common.PrintMemUsage()
}

func TestRedisPing(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "47.99.198.186:6600", // Redis地址
		Password: password,
	})
	ctx := context.Background()
	defer rdb.Close()

	// 检查Redis中是否有数据
	keys, err := rdb.Keys(ctx, "ad:bytedance:*").Result()
	if err != nil {
		t.Fatalf("failed to get keys: %v", err)
	}

	if len(keys) == 0 {
		fmt.Println("No data found in Redis.")
		return
	}

	fmt.Printf("Found %d keys in Redis:\n", len(keys))
	for _, key := range keys {
		// 获取集合的大小
		card, err := rdb.SCard(ctx, key).Result()
		if err != nil {
			fmt.Printf("failed to get cardinality for key %s: %v\n", key, err)
		} else {
			fmt.Printf("key: %s, cardinality: %d\n", key, card)
		}
		rtaIDStr, found := strings.CutPrefix(key, RatRedisKeyPrefix)
		fmt.Println("rta id string:", found, rtaIDStr)
		rid, err := strconv.ParseInt(rtaIDStr, 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("rid integer:", rid)
		// 使用 SSCAN 分批获取集合成员
		var cursor uint64
		for {
			members, nextCursor, err := rdb.SScan(ctx, key, cursor, "", 1000).Result()
			if err != nil {
				fmt.Printf("failed to scan members for key %s: %v\n", key, err)
				break
			}

			if len(members) > 0 {
				fmt.Printf("key: %s, members: %v\n", key, members)
			}
			cursor = nextCursor
			if cursor == 0 {
				break
			}
		}
	}
}

func TestRedisQuerySomeData(t *testing.T) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "47.99.198.186:6600", // Redis地址
		Password: password,
	})
	ctx := context.Background()
	defer rdb.Close()

	// 检查Redis中是否有数据
	keys, err := rdb.Keys(ctx, "ad:bytedance:*").Result()
	if err != nil {
		t.Fatalf("failed to get keys: %v", err)
	}

	if len(keys) == 0 {
		fmt.Println("No data found in Redis.")
		return
	}

	fmt.Printf("Found %d keys in Redis:\n", len(keys))
	for _, key := range keys {
		// 获取集合的大小
		card, err := rdb.SCard(ctx, key).Result()
		if err != nil {
			fmt.Printf("failed to get cardinality for key %s: %v\n", key, err)
		} else {
			fmt.Printf("key: %s, cardinality: %d\n", key, card)
		}
		rtaIDStr, found := strings.CutPrefix(key, RatRedisKeyPrefix)
		fmt.Println("rta id string:", found, rtaIDStr)
		rid, err := strconv.ParseInt(rtaIDStr, 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("rid integer:", rid)
		// 使用 SSCAN 分批获取集合成员
		var cursor uint64
		members, nextCursor, err := rdb.SScan(ctx, key, cursor, "", 1000).Result()
		if err != nil {
			fmt.Printf("failed to scan members for key %s: %v\n", key, err)
			break
		}

		if len(members) > 0 {
			fmt.Printf("key: %s, members: %v\n", key, members)
		}
		cursor = nextCursor
	}
}
