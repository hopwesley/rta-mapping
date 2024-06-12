package main

import (
	"database/sql"
	"flag"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
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

func TestInitRtaMap(t *testing.T) {
	cfg := &RedisCfg{
		Addr: "47.99.198.186:6600",
	}
	cfg.Password = password
	err := InitRtaMap(cfg)
	if err != nil {
		fmt.Println("connect to redis failed:", err)
		t.Fatal(err)
	}
}
