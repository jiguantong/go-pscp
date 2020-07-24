package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io"
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

type Option struct {
	Key   string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}
type Component struct {
	Options []Option `xml:"option"`
}
type XmlConfig struct {
	XMLName xml.Name  `xml:"project"`
	C       Component `xml:"component"`
}

func main() {
	var puttyPath string
	path := strings.Split(os.Getenv("PATH"), ";")
	for _, v := range path {
		if strings.Contains(v, "PuTTY") {
			puttyPath = v + "pscp.exe"
			break
		}
	}
	if len(os.Args) > 2 {
		puttyPath = os.Args[2]
	} else if len(puttyPath) == 0 {
		fmt.Println("PATH 中未找到 pscp, 可通过第二个参数指定pscp路径")
		os.Exit(1)
	}

	config := loadConf()
	cmd := exec.Command(puttyPath, "-r", "-pw",
		config.Password, config.Localdir,
		config.User+"@"+config.Ip+":"+
			config.Remotedir)
	//fmt.Println(cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, _ := cmd.StdinPipe()
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}
	io.WriteString(stdin, "y")
	stdin.Close()
	cmd.Wait()
	//fmt.Println("### -> 传输完成")
	fmt.Println("### => Push complete")
	runCmd(*config)
}

func loadConf() *Config {
	var c = new(Config)
	filePath := "./pscp.yml"
	if len(os.Args) > 1 {
		// 指定, 读取xml配置
		filePath = os.Args[1]
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer file.Close()
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		var xmlC = new(XmlConfig)
		if err := xml.Unmarshal(data, xmlC); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		optionMap := make(map[string]string)
		for _, o := range xmlC.C.Options {
			optionMap[o.Key] = o.Value
		}
		c.Ip = optionMap["ip"]
		c.User = optionMap["user"]
		c.Port = optionMap["port"]
		c.Password = optionMap["pwd"]
		c.Localdir = optionMap["localDir"]
		c.Remotedir = optionMap["remoteDir"]
		c.Cmd = optionMap["cmd"]
	} else {
		// 未指定文件路径, 读取当前目录下的yml
		ymlFile, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if err := yaml.Unmarshal(ymlFile, c); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
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
	//fmt.Println("### -> 正在连接服务器...")
	fmt.Println("### => Connecting to server...")
	client, err := ssh.Dial("tcp", conf.Ip+":"+conf.Port, config)
	if nil != err {
		fmt.Println(err)
		return
	}
	// 创建与远程服务器的会话
	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("### -> 执行: ", conf.Cmd)
	fmt.Println("### => Run cmd: ", conf.Cmd)
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
