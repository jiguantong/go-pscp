package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	viper.SetConfigName("pscp")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("配置文件读取异常: " + err.Error())
		os.Exit(1)
	}
	var puttyPath string
	path := strings.Split(os.Getenv("PATH"), ";")
	for _, v := range path {
		if strings.Contains(v, "PuTTY") {
			puttyPath = v
			break
		}
	}
	if len(puttyPath) == 0 {
		log.Printf("PATH 中未找到 PuTTY, 请先安装 PuTTY\ndownload: %s\n", "https://the.earth.li/~sgtatham/putty/0.73/w64/putty-64bit-0.73-installer.msi")
		os.Exit(1)
	}

	localDir := viper.GetString("local.dir")
	host := new(Host)
	host.Dir = viper.GetString("host.dir")
	host.User = "root"
	host.Ip = viper.GetString("host.ip")
	host.Pwd = viper.GetString("host.password")
	cmd := exec.Command(puttyPath+"pscp.exe", "-r", "-pw",
		host.Pwd, localDir,
		host.User+"@"+host.Ip+":"+
			host.Dir)
	fmt.Println(cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
