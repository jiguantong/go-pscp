package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("配置文件读取异常")
		os.Exit(1)
	}
	localDir := viper.GetString("local.dir")
	host := new(Host)
	host.Dir = viper.GetString("host.dir")
	host.User = "root"
	host.Ip = viper.GetString("host.ip")
	host.Pwd = viper.GetString("host.password")
	fmt.Println(host)
	fmt.Println("localDir: " + localDir)
	cl := "pscp -r -pw " + host.Pwd + " " + localDir + " " + host.User + "@" + host.Ip + ":" + host.Dir
	cmd := exec.Command("cmd", "/C", cl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
