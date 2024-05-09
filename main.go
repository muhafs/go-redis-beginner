package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ctx = context.Background()

type Member struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func main() {
	// connect to database
	dsn := "root:root@tcp(127.0.0.1:3306)/redisdb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("cannot conntect to database:", err)
		return
	}

	// retrieve data from db
	var members []Member
	db.Model(&Member{}).Find(&members)

	// initialize redis connection
	rClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// encode data to json
	membersbyte, err := json.Marshal(members)
	if err != nil {
		fmt.Println("couldn't encode data into json:", err)
		return
	}

	// store data into redis cache
	err = rClient.Set(ctx, "members", membersbyte, 0).Err()
	if err != nil {
		fmt.Println("cannot set value:", err)
		return
	}

	// retrieve data from redis cache
	val, err := rClient.Get(ctx, "members").Result()
	if err != nil {
		fmt.Println("cannot get value:", err)
		return
	}

	// decode data from json
	var result []Member
	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		fmt.Println("couldn't decode data from json:", err)
		return
	}

	// show result
	fmt.Println(result[0])
}
