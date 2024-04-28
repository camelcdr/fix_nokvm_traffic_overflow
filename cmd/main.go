package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readMysqlPwd(filepath string) string {
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

func main() {
	var dir = "/home/nokvm/go/src/node/bin/conf"
	var filename = "app.conf"
	filepath := dir + "/" + filename
	passwd := readMysqlPwd(filepath)
	fmt.Println("当前的nokvm被控的数据库密码为：", passwd)
}
