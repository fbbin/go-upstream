package main

import (
	"net"
	"time"
)

// 初始化代理服务
func initProxy() {
	Log.Infof("Proxying %s -> %s\n", Config.Bind, Config.Backend)
	server, err := net.Listen("tcp", Config.Bind)
	if err != nil {
		Log.Fatal(err)
	}
	// 等待的队列长度
	waitQueue := make(chan net.Conn, Config.WaitQueueLen)
	// 最大并发连接
	connPools := make(chan bool, Config.MaxConn)
	for i := 0; i < Config.MaxConn; i++ {
		connPools <- true
	}
	// 等待连接处理
	go waitConn(waitQueue, connPools)
	// 接收连接并抛给管道处理
	for {
		conn, err := server.Accept()
		if err != nil {
			Log.Error(err)
			continue
		}
		Log.Infof("Received connection from %s.\n", conn.RemoteAddr())
		waitQueue <- conn
	}
}

// 连接数控制
func waitConn(waitQueue chan net.Conn, connPools chan bool) {
	for conn := range waitQueue {
		// 接收一个链接，连接池释放一个
		<-connPools
		go func(conn net.Conn) {
			handleConn(conn)
			// 链接处理完毕，增加
			connPools <- true
			Log.Infof("Closed connection from %s.\n", conn.RemoteAddr())
		}(conn)
	}
}

// 处理连接
func handleConn(conn net.Conn) {
	defer conn.Close()
	// 根据链接哈希选择机器
	proxySvr, ok := getBackendSvr(conn)
	if !ok {
		return
	}
	// 链接远程代理服务器
	remote, err := net.Dial("tcp", proxySvr.identify)
	if err != nil {
		Log.Error(err)
		proxySvr.failTimes++
		return
	}
	// 等待双向连接完成
	complete := make(chan bool, 2)
	oneSwitch := make(chan bool, 1)
	otherSwitch := make(chan bool, 1)
	// 将当前客户端链接发送的数据发送给远程被代理的服务器
	go transaction(conn, remote, complete, oneSwitch, otherSwitch)
	// 将远程服务返回的数据返回给客户端
	go transaction(remote, conn, complete, otherSwitch, oneSwitch)
	<-complete
	<-complete
	remote.Close()
}

// 数据交换传输（从from读数据，再写入to）
func transaction(from, to net.Conn, complete, oneSwitch, otherSwitch chan bool) {
	var err error
	var read int
	bytes := make([]byte, 1024)
	for {
		select {
		case <-otherSwitch:
			complete <- true
			return
		default:
			timeOutSec := time.Duration(Config.Timeout) * time.Second
			// 设置超时时间
			from.SetReadDeadline(time.Now().Add(timeOutSec))
			read, err = from.Read(bytes)
			if err != nil {
				complete <- true
				oneSwitch <- true
				return
			}
			// 设置超时时间
			to.SetWriteDeadline(time.Now().Add(timeOutSec))
			_, err = to.Write(bytes[:read])
			if err != nil {
				complete <- true
				oneSwitch <- true
				return
			}
		}
	}
}
