##go-upstream
服务代理中间件，支持后端服务集群配置，根据HASH配置选择后端机器进行代理。

##安装方法
```
go get github.com/fbbin/go-upstream
```

##Examples

```
bind: 0.0.0.0:9800
wait_queue_len: 100
max_conn: 50
timeout: 5
failover: 3
stats: 0.0.0.0:8090
backend:
    - 192.168.163.184:20001
    - 192.168.163.184:20001
    - 192.168.163.184:20001

log:
    level: "info"
    path: "./logs/proxy.log"
```