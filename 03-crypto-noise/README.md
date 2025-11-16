# crypto-noise

## 目的
- 专注 `noise-protocol` 相关的 Go 结构，如 `HandshakeState`, `Keypair`。
- 梳理密钥生命周期：生成、轮换、过期与重协商。

## 原理简介
wireguard-go 内置 NoiseIK：
1. 使用 `x/crypto/chacha20poly1305`、`blake2s`、`x25519` 等库构建 AEAD。
2. Handshake 由 `device/noise-protocol`、`device/noise-handshake` 实现，密钥滚动通过定时器触发。
3. `Keypair` 切换遵循 `Reject-After-Messages/Time` 策略，保障前向安全。
本目录将记录相关结构体字段含义与调用示意。
