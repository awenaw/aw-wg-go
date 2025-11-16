# WireGuard Go 源码学习项目

## 目标
- 以 wireguard-go 官方实现为主线，梳理模块职责与调用链路。
- 将协议手册与 Go 代码互相映射，形成可查阅的学习档案。
- 为后续阅读、调试、实验记录提供统一的目录结构。

## 学习节奏
1. 自上而下先理解协议、角色与部署模型。
2. 再进入核心 `device` 组件，建立数据平面与控制平面的心智模型。
3. 最后通过实验与工具验证推导结果。

## 目录速览
- `protocol-overview`：WireGuard 协议与数据流速记。
- `core-device`：Go 版 `device` 包的结构拆解。
- `crypto-noise`：NoiseIK 握手与加解密链路。
- `packet-flow`：数据包进出内核/用户空间的路径追踪。
- `tun-platform`：跨平台 TUN/TAP 适配层解析。
- `tooling-labs`：实验脚本、可视化与调试技巧。
