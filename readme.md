### IM-System with Golang
___生成并执行可执行文件创建客户端___
```bash
    go build -o server main.go server.go user.go
    ./server
``` 
___创建客户端___
```bash
    nc 127.0.0.1 8888
```