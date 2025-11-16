# tooling-labs

## 目的
- 收集学习过程中用到的 Go 小实验、抓包脚本、可视化工具。
- 记录调试 wireguard-go 时的命令组合与常见陷阱。

## 原理简介
围绕 Go 工具链与网络调试：
1. 使用 `go test`, `go tool trace`, `pprof` 观察性能。
2. 借助 `tcpdump`, `wireshark`, `wg` CLI 佐证协议行为。
3. 自制实验（如模拟握手、构造 UDP 流）验证对源码的理解。
该目录强调“实践—验证—再阅读”的闭环。
