

## server 端

- key/value 语义
- 入参：
    - id：数据的唯一 id
    - body：数据的具体内容
- 返回值：
    - 位置元数据

## client 端

- 负责管理元数据：文件 -> 位置元数据


假定：

1. 只有一个目录，那就是根目录。
2. 每个 Server 一次性只能写一个文件，重启之后换文件。
3. 文件不能大于 1M
4. 不能连续写。
5. 每次读都是全读

## 编译部署

1. 编译

```
make
```

确认 server，client 这两个二进制编译成功。

2. 启动 server

```
./start_servers.sh
```

这个脚本会在本地拉起 6 个进程，对应模拟 6 个节点的 server 。
> 如果出现问题，查看对应日志。

3. 启动 client

```
./start_client.sh
```

拉起一个 client 进程，用于承接 /mnt/hellofs 挂载点的请求。
> 如果出现问题，查看对应日志。

## 关闭、卸载

1. 卸载挂载点，kill client

```
./stop_client.sh
```

2. kill server

```
./stop_servers.sh
```
