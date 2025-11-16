# core-device

## 目的
- 深入阅读 `device` 结构体、队列、工作协程调度与配置接口。
- 搞清楚 `device.NewDevice`、`Start`、`Reconfigure` 等入口如何串联。

## 原理简介
`device` 包是 wireguard-go 的中枢：
1. 通过 `device.Device` 聚合 TUN、UDP、定时器、密钥等资源。
2. `device.Routine*` 系列协程形成数据管道：从内核 TUN -> 加密 -> UDP，以及反向流程。
3. 配置通过 UAPI（userspace API）传入，最终落在 peer/allowed IP 等结构。
分析这些文件可帮助理解 Go 版如何在无内核模块的情况下复现 WireGuard 行为。
