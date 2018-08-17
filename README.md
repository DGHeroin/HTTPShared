# HTTPShared

[![Build Status](https://travis-ci.org/DGHeroin/HTTPShared.svg?branch=master)](https://travis-ci.org/DGHeroin/HTTPShared)
[![codecov](https://codecov.io/gh/DGHeroin/HTTPShared/branch/master/graph/badge.svg)](https://codecov.io/gh/DGHeroin/HTTPShared)

### 写入
写入一个值, 键为message 值为hello (["message"] = "hello"), 需要用 HTTP PUT 操作
```
curl http://127.0.0.1:9999/v1/keys/message -X PUT -d hello
```

### 读取
读取key为message的值
```
curl http://127.0.0.1:9999/v1/keys/message
```

### 监听值的变化
监听一次变化
```
curl http://127.0.0.1:9999/v1/keys/message?wait=true
```
一直监听变化
```
curl http://127.0.0.1:9999/v1/keys/message?wait=true&r=true
```
一直监听变化(每次返回json之前, 先写入4个字节int数据用来标识接下来json的长度)
```
curl http://127.0.0.1:9999/v1/keys/message?wait=true&r=true&h=true
```
