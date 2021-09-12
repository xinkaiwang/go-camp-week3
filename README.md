# 本周作业

## 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

## 示例

```
xinkai@mbp15s:go-camp-week3$ go build server.go 
xinkai@mbp15s:go-camp-week3$ ./server 
```

从另一个窗口
```
xinkai@mbp15s:go-camp-week3$ curl http://localhost:8080/ping
Hi there, handler=h1, req=ping
xinkai@mbp15s:go-camp-week3$ curl http://localhost:8082/pong
Hi there, handler=h2, req=pong
xinkai@mbp15s:go-camp-week3$ curl http://localhost:8080/close
Hi there, handler=h1, req=close
```

这时server 会退出
```
closing server h1 by API
signal go routine exit
closing server 1 by graceful shutdown
closing server 2 by graceful shutdown
exit reason: err=http: Server closed
```

另外，CTR-C 也能让server 优雅退出。
```
xinkai@mbp15s:go-camp-week3$ ./server 
^Csignal recrived: interrupt
closing server 1 by graceful shutdown
closing server 2 by graceful shutdown
exit reason: err=signal
```