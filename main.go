package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
)

var (
	Config     ProxyConfig
	Log        *logrus.Logger
	configFile = flag.String("c", "./conf.yaml", "配置文件：conf.yaml")
)

func onExitSignal() {
	signalChan := make(chan os.Signal)
	// 监听系统服务退出信号
	signal.Notify(signalChan, syscall.SIGUSR1, syscall.SIGTERM, syscall.SIGINT, os.Kill)
	for {
		signal := <-signalChan
		log.Println("Get Signal:%v\r\n", signal)
		switch signal {
		case syscall.SIGTERM, syscall.SIGINT, os.Kill:
			log.Fatal("系统退出。。。")
		}
	}
}
func main() {

	flag.Parse()
	fmt.Println("Start Proxy...")

	// 解析配置
	parseConfigFile(*configFile)

	// 初始化日志模块
	initLogger()

	// 初始化代理的服务
	initBackendSvrs(Config.Backend)

	// 系统退出信号监听
	go onExitSignal()

	// 初始化状态服务
	initStats()

	// 初始化代理服务
	initProxy()

}
