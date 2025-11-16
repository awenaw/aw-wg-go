# tun-platform

## 目的
- 解析 wireguard-go 如何抽象不同平台的 TUN/TAP 设备。
- 记录 Windows、Linux、macOS 各自的驱动调用差异与封装策略。

## 原理简介
TUN 模块通过 `tun` 包封装：
1. Linux/macOS 走原生 TUN FD；Windows 借助 wintun DLL 交互。
2. `tun.Device` 提供统一的 `Read`, `Write`, `Events` 接口，供 `device` 调用。
3. 适配层额外处理 MTU、IPv6/IPv4、接口状态等。
本目录将把系统调用/WinAPI 与 Go 抽象的映射关系整理清楚。
