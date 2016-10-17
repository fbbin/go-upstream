package main

import (
	"fmt"
	"net/http"
)

// 查询监控信息的接口
func statsHandler(w http.ResponseWriter, r *http.Request) {
	_str := ""
	for _, v := range BackendSvrs {
		_str += fmt.Sprintf("Server:%s FailTimes:%d isUp:%t\n", v.identify, v.failTimes, v.isLive)
	}
	fmt.Fprintf(w, "%s", _str)
}

// 初始化监控服务地址
func initStats() {
	Log.Infof("Start monitor on addr %s", Config.Stats)

	go func() {
		http.HandleFunc("/stats", statsHandler)
		http.ListenAndServe(Config.Stats, nil)
	}()
}
