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
