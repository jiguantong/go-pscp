package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	User      string
	Password  string
	Ip        string
	Remotedir string
	Localdir  string
}

func main() {
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

	config := loadConf()
	cmd := exec.Command(puttyPath+"pscp.exe", "-r", "-pw",
		config.Password, config.Localdir,
		config.User+"@"+config.Ip+":"+
			config.Remotedir)
	fmt.Println(cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
}

func loadConf() *Config {
	var c = new(Config)
	ymlFile, err := ioutil.ReadFile("./pscp.yml")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err := yaml.Unmarshal(ymlFile, c); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return c
}
