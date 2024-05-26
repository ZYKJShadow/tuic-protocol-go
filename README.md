## TUIC-PROTOCOL-GO

本协议已经有前人栽树，请参考[TUIC Protocol](https://github.com/EAimTY/tuic/blob/tuic-5.0.0/SPEC.md)
<br>
本项目是以学习为目的开发和开源，任何其他用途与作者无关<br>
Golang客户端：[tuic-client](https://github.com/ZYKJShadow/tuic-client)<br>
Golang服务器：[tuic-server](https://github.com/ZYKJShadow/tuic-server)

## 协议
TUIC 协议依赖于一个可以多路复用的 TLS 加密流。所有的中继任务都通过 Command 中的 Header 来协商。
### 协议版本
`0x05`
### Command
```plain
+-----+------+----------+
| VER | TYPE |   OPT    |
+-----+------+----------+
|  1  |  1   | Variable |
+-----+------+----------+
```
- VER: 协议版本，固定为`0x05`
- TYPE: Command类型
- OPT: Command参数，参考[Command 参数](#Command-参数)
### Command 类型

目前有五种Command类型：
- `0x00` Authenticate - 认证
- `0x01` Connect - 建立TCP中继
- `0x02` Packet - 传输UDP中继（分片）
- `0x03` Dissociate - 终止UDP中继
- `0x04` Heartbeat - 心跳

其中 Connect 和 Packet 携带有效载荷（流/数据包片段）
### Command 参数

#### `Authenticate`
```plain
+------+-------+
| UUID | TOKEN |
+------+-------+
|  16  |  32   |
+------+-------+
```
- `UUID` 客户端UUID
- `TOKEN` 客户端令牌，客户端的UUID作为`label`，客户端密码作为`context`，从TLS连接中生成的32位密钥

#### `Connect`

```plain
+----------+
|   ADDR   |
+----------+
| Variable |
+----------+
```

- `ADDR`: 目标地址信息，参考[地址信息](#地址信息)：

#### `Dissociate`

```plain
+----------+
| ASSOC_ID |
+----------+
|    2     |
+----------+
```
- `ASSOC_ID` - UDP中继的关联ID，参考[UDP中继](#3、UDP中继（可选）)

#### `Heartbeat`
```plain
+-+
| |
+-+
| |
+-+
```

### 地址信息
### `ADDR`
```plain
+------+-------------+----------+
| TYPE |   ADDRESS   |   PORT   |
+------+-------------+----------+
|  1   |   Variable  |    2     |
+------+-------------+----------+
```

- `TYPE` - 地址类型
- `ADDRESS` - 地址
- `PORT` - 端口号

地址类型：
- `0xff`：无
- `0x00`：域名（类型为域名时，ADDRESS部分的第一个字节表示域名长度）
- `0x01`：IPv4 地址
- `0x02`：IPv6 地址

## 整体流程

### 1、认证
- 客户端打开一个单向流(uni_stream)，发送`Authenticate`命令
- 服务器接收到`Authenticate`后，验证 Token 的有效性。如果 Token 有效，则连接认证成功，进行其他中继任务
- 如果服务器在接收到`Authenticate`命令之前收到其他命令，则只接受命令头部分并等待认证，认证成功后通知正在等待的协程继续进行任务

### 2、TCP中继
- 客户端打开一个双向流(bi_stream)，发送`Connect`命令
- 客户端在`Connect`命令发送成功之后，立刻对双向流和本地连接建立双向传输
- 服务器接收到`Connect`命令后，打开一个到目标地址的TCP流，建立成功后，服务器立即在TCP流和双向流之间传输数据、
- 有一方数据传输完毕，双方都关闭并释放流和连接

### 3、UDP中继
- `ASSOC_ID`由客户端生成
- 客户端可以通过QUIC 单向流(quic模式)或 QUIC datagram(native模式)发送`Packet`命令
- TUIC 通过在客户端和服务器之间同步UDP`ASSOC_ID`来实现0-RTT
- 客户端和服务器为每个 QUIC 连接创建一个 UDP 会话表，将每个`ASSOC_ID`映射到一个UDP socket
- 服务器为每一个`ASSOC_ID`分配一个UDP socket。服务器使用这个UDP socket发送客户端请求的 UDP 数据包，同时接收来自任何目的地的UDP数据包，并添加`Packet`命令头发送回客户端
- 客户端可以通过 QUIC 单向流发送 Dissociate 命令来解除关联

### 4、 心跳
客户端定期通过Datagram的方式发送`Heartbeat`命令

