package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
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
	// 传输后执行该命令
	Cmd  string
	Port string
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
	fmt.Println("### -> 传输完成")
	runCmd(*config)
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

func runCmd(conf Config) {
	if len(conf.Cmd) == 0 {
		return
	}
	// 配置连接
	config := &ssh.ClientConfig{}
	config.SetDefaults()
	config.User = conf.User
	config.Auth = []ssh.AuthMethod{ssh.Password(conf.Password)}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}
	client, err := ssh.Dial("tcp", conf.Ip+":"+conf.Port, config)
	if nil != err {
		fmt.Println(err)
		return
	}
	fmt.Println("### -> 正在连接服务器...")
	// 创建与远程服务器的会话
	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("### -> 服务器连接成功!")
	fmt.Println("### -> 执行: ", conf.Cmd)
	defer session.Close()
	cmdReader, err := session.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	logScanner := bufio.NewScanner(cmdReader)
	_logChan := make(chan []byte, 300)
	go func(logScan *bufio.Scanner, logChan chan<- []byte) {
		for logScan.Scan() {
			_logChan <- []byte(logScan.Text())
		}
	}(logScanner, _logChan)
	if err = session.Start(conf.Cmd); err != nil {
		fmt.Println(err)
		return
	}
	for log := range _logChan {
		_log := string(log)
		fmt.Println(_log)
	}
}
