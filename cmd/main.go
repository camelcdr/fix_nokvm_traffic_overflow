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

func readInstanceTrafficLimits(filepath string) map[string]int64 {
	fmt.Println("开始读取nokvm.data文件...")
	trafficLimits := make(map[string]int64)

	json, err := os.ReadFile(filepath)
	fmt.Println("成功读取nokvm.data文件.")

	guests := gjson.Get(string(json), "Guests")
	if err != nil {
		fmt.Println("读取nokvm.data文件失败", err)
		return trafficLimits
	}
	guests.ForEach(func(key, value gjson.Result) bool {
		s := value.String()
		name := gjson.Get(s, "name").String()
		trafficLimit := gjson.Get(s, "traffic_limit").Int()
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

	var filepath = "/home/nokvm/nokvm.data"
	trafficLimits := readInstanceTrafficLimits(filepath)
	count := 0
	for name, traffic := range trafficLimits {
		fmt.Println(count, "\tName:", name, "traffic:", traffic)
		count++
	}

	//correctDb(passwd)
}
