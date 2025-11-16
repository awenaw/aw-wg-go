# 01 - protocol-overview

## 章节定位
- 不是零基础课堂，默认你对网络安全/加密握手有概念，或者已经写过其它语言的安全通信代码。
- 目标是在 WireGuard 协议概念与 wireguard-go 源码之间搭建“翻译层”，帮助后续阅读核心 `device` 包时少踩坑。
- 如果你刚学完 `00-helloworld-go`，这里就是下一站：不再教语法，而是开始讨论协议与模块职责。

## 学习目标
- 用一段简化代码把 WireGuard 白皮书中的 Initiator/Responder/Transport 流程串起来。
- 理解 NoiseIK 模式下 CK/SK（Chain Key / Session Key）的派生顺序，并知道这些符号在源码里对应的字段。
- 学会如何把 README 中的流程示意映射到 `main.go` 的日志，再映射到真正的 wireguard-go 文件。

## 原理速记
1. WireGuard 的用户态实现沿用 NoiseIK：Curve25519 进行 Diffie-Hellman，ChaCha20-Poly1305 负责加解密，BLAKE2s 参与哈希链。
2. 握手只有两条消息：`HandshakeInitiation` 与 `HandshakeResponse`；一旦成功，双方立即切到 `TransportData`。
3. 每个 Peer 都维护发送/接收计数器与 Keypair，wireguard-go 中由 `device.Routines` 与定时器驱动轮换。

## 运行条件
- Windows PowerShell 或任一 shell，工作目录为 `c:\Users\kisst\prj\github\aw-wg-go`。
- 已安装 Go 1.22+，`go version` 可正常输出。

## 如何阅读 main.go
1. **先读结构体**：`Peer` 只保存名字和静态公钥；`Message` 用 map 存字段；`Simulation` 则收集日志与密钥演化。搞清楚这是为“讲故事”准备的数据模型。
2. **浏览辅助函数**：`logf`、`derive`、`send`、`reply`、`describeMessage` 负责生成 Transcript，每个注释都告诉你它们在 NoiseIK 中扮演的角色。
3. **逐行阅读 `run()`**：它是脚本，分三段——发起握手、响应握手、进入 Transport。每段先写旁白，再调用 `derive`、`send/reply`，与输出一一对应。
4. **运行并对照**：执行 `go run .`，观察终端输出的每一行都能在 `run()` 找到源头。字段顺序如果和注释不同，是因为 map 遍历无序，不影响含义。
5. **联想到真实源码**：当搞懂 Transcript 表示什么后，跳到 wireguard-go 的 `device/noise-handshake.go`、`device/peer.go` 查同名概念，完成“示例 → 真实实现”的映射。

## 运行示例
```powershell
cd 01-protocol-overview
go run .
```
输出的 Transcript 就是 `Simulation.run()` 写下的完整流程说明，可多跑几次并在注释里补充自己的理解。

## 本章节收获
- **角色心智模型**：明白 WireGuard 只有两个对等端，靠静态公钥确定身份。
- **三类消息/字段**：初步熟悉 Initiation/Response/Transport 包含什么元素，便于抓包或阅读协议文档。
- **CK/SK 演化链**：通过 `keyDerivations` 日志掌握链式派生顺序，为后续阅读 keypair 轮换打基础。
- **阅读路径**：知道先看 README → `main.go` → 真实源码的顺序，不会在复杂的 `device` 目录迷路。

> 如果觉得仍旧抽象，可以把 `Simulation` 想成“伪代码版 Wireshark”：先理解日志，再去真实代码里寻找同款函数。必要时不妨在 README 里追加自己的笔记或 TODO。
