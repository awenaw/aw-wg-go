# protocol-overview

## 目的
- 汇总 WireGuard 白皮书、RFC 草案与 wireguard-go 源码中涉及的高层协议细节。
- 建立握手、数据传输、Keepalive、重试等 state machine 的鸟瞰图。

## 原理简介
WireGuard 在 Go 实现中严格复刻 NoiseIK 模式：
1. 以 Curve25519+ChaCha20-Poly1305 组合完成密钥协商与数据平面加密。
2. Handshake 消息通过 UDP 承载，包含 Initiation、Response、Cookie 交换。
3. 每个 Peer 维护发送/接收加密计数器，通过 `device.Routines` 驱动。
该目录用于把协议状态与源码文件（如 `device/peer.go`, `device/noise.go`）建立映射。

## 运行条件
- 已安装 Go 1.22（或兼容版本），并能通过 `go version` 正常输出。
- 当前工作目录为 `c:\Users\kisst\prj\github\aw-wg-go`，本章节在 `01-protocol-overview/` 下。

## 运行示例
```powershell
cd 01-protocol-overview
go run .
```
输出日志即为 `main.go` 中的 NoiseIK 握手/传输示例 Transcript，可用于对照 README 的流程描述。

## 入门者可以学到什么
- **协议角色的心智模型**：`main.go` 通过 `Peer` 结构体展示 Initiator/Responder 的最小关键信息（名字、公钥），帮助理解 WireGuard 只有静态密钥、没有传统证书体系的设计。
- **三类核心消息**：`HandshakeInitiation`、`HandshakeResponse`、`TransportData` 会在输出中以 `Sender -> Receiver | Kind | fields` 形式呈现，方便对照官方白皮书中的字段。
- **NoiseIK 的 CK/SK 演化**：`keyDerivations` 日志按顺序列出 CK0 ~ CK3 与 SKi/SKr 的推导，让初学者能快速掌握“MixHash / MixKey / KDF2”在握手中的作用。
- **代码模块划分**：`Simulation` 结构用脚本化的方式串联握手步骤，和 wireguard-go 里 `device/noise-handshake.go` 的状态机一一对应，便于后续跳到真实源码查阅。
- **运行→观察→推导的闭环**：通过 `go run .` 即可看到全流程输出，先理解高层概念，再去源码中寻找实现细节（如 peer 管理、定时器、UDP 协程）。

> 建议：先阅读本 README 的“原理简介”，运行示例观察日志，再转向 `main.go` 注释或官方文档深入理解。必要时可把输出与抓包工具（Wireshark/wg CLI）的结果对比，加深对消息格式与计数器的认识。
