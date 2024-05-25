## TUIC-PROTOCOL-GO

本协议已经有前人栽树，请参考[TUIC Protocol](https://github.com/EAimTY/tuic/blob/tuic-5.0.0/SPEC.md)，作为学习HTTP3，用Golang实现了一遍，仅供娱乐。

根据协议，将代理转发过程分为以下几个部分：

### 1、认证
- 客户端打开一个单向流(unidirectional_stream)，发送 Authenticate 命令，其中包含客户端的 UUID 和 Token
- 服务器接收到 Authenticate 命令后，验证 Token 的有效性。如果 Token 有效，则连接认证成功，可以进行其他中继任务
- 如果服务器在接收到 Authenticate 命令之前收到其他命令，它应该只接受命令头部分并暂停。连接认证成功后，服务器应该恢复所有暂停的任务

### 2、TCP中继（可选）
- 客户端打开一个双向流(bidirectional_stream)，发送 Connect 命令，其中包含目标地址信息
- 客户端在命令头传输完成后，可以立即开始使用该流进行 TCP 中继，无需等待服务器的响应
- 服务器接收到 Connect 命令后，会打开一个到目标地址的 TCP 流，在 TCP 流建立后，服务器可以开始在 TCP 流和双向流之间传输数据

### 3、UDP中继（可选）
- TUIC 通过在客户端和服务器之间同步 UDP 会话 ID(associate ID)来实现 0-RTT 全锥型 UDP 转发
- 客户端和服务器都应该为每个 QUIC 连接创建一个 UDP 会话表,将每个 associate ID 映射到一个关联的 UDP socket
- 客户端生成一个 16 位无符号整数作为 associate ID。如果客户端想要使用服务器的同一个 UDP socket 发送数据包，Packet 命令中附加的 associate ID 应该保持一致
- 服务器在接收到 Packet 命令时，应检查附加的 associate ID 是否已经与一个 UDP socket 关联。如果没有，服务器应该为该 associate ID 分配一个 UDP socket。服务器将使用这个 UDP socket 发送客户端请求的 UDP 数据包,同时接受来自任何目的地的 UDP 数据包，为其添加 Packet 命令头，然后发送回客户端
- 一个 UDP 数据包可以被分片成多个 Packet 命令。字段 PKT_ID、FRAG_TOTAL 和 FRAG_ID 用于标识和重组分片的 UDP 数据包
- 客户端可以通过QUIC 单向流(quic模式)或 QUIC datagram(native模式)发送 Packet 命令
- 服务器在从一个 UDP 中继会话(associate ID)接收到第一个 Packet 时，应该使用相同的模式发回 Packet 命令
- 客户端可以通过 QUIC 单向流发送 Dissociate 命令来解除 UDP 会话关联。服务器将移除 UDP 会话并释放关联的 UDP socket

### 4、 心跳
当有任何正在进行的中继任务时，客户端应该定期通过 QUIC datagram 发送 Heartbeat 命令,以保持 QUIC 连接的活跃状态。

