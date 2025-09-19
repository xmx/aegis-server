## 编译

```shell
GOEXPERIMENT=jsonv2 sh build.sh

# 交叉编译
GOOS=windows GOARCH=amd64 GOEXPERIMENT=jsonv2 sh build.sh
```

## 开发

> go1.25.0+
>
> mongodb 8.0.0+

```shell
git clone https://github.com/xmx/aegis-common.git
git clone https://github.com/xmx/aegis-control.git
git clone https://github.com/xmx/aegis-agent.git
git clone https://github.com/xmx/aegis-broker.git
git clone https://github.com/xmx/aegis-server.git

go work init
go work use aegis-common aegis-control aegis-agent aegis-broker aegis-server
```

- [ageis-server](https://github.com/xmx/aegis-server): 服务端程序。

- [aegis-broker](https://github.com/xmx/aegis-broker): 代理、调度端程序。

- [aegis-agent](https://github.com/xmx/aegis-agent): agent 程序。

- [aegis-control](https://github.com/xmx/aegis-control): server, broker 共用的代码。

- [aegis-common](https://github.com/xmx/aegis-common): server, broker, agent 共用的代码。
