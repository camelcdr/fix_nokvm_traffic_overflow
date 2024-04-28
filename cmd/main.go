package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/tidwall/gjson"
	"os"
	"strings"
)

func readMysqlPwd() string {
	var dir = "/home/nokvm/go/src/node/bin/conf"
	var filename = "app.conf"
	filepath := dir + "/" + filename
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("打开配置文件失败：", err)
		panic("打开配置文件失败")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "db_pass") {
			parts := strings.SplitN(line, "=", 2)
			passwd := strings.TrimSpace(parts[1])
			return passwd
		}
	}
	panic("获取密码失败")
}

func readInstanceTrafficLimits() map[string]string {
	trafficLimits := make(map[string]string)
	var filepath = "/home/nokvm/nokvm.data"
	file, err := os.Open(filepath)
	if err != nil {
		return trafficLimits
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	json := scanner.Text()
	guests := gjson.Get(json, "Guests")
	guests.ForEach(func(key, value gjson.Result) bool {
		name := gjson.Get(json, "name").String()
		trafficLimit := gjson.Get(json, "traffic_limit").String()
		trafficLimits[name] = trafficLimit
		return true
	})
	return trafficLimits
}

func correctDb(passwd string) {
	db, err := sql.Open("mysql", "username:"+passwd+"@tcp(127.0.0.1:3306)/node")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

}

func main() {
	passwd := readMysqlPwd()
	fmt.Println("当前的nokvm被控的数据库密码为：", passwd)

	trafficLimits := readInstanceTrafficLimits()
	count := 0
	for name, traffic := range trafficLimits {
		fmt.Println(count, "\tName:", name, "traffic:", traffic)
		count++
	}

	//correctDb(passwd)
}
