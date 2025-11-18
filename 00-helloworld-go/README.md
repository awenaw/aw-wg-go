# 00 - 基础积木与学习路线提纲

本目录用于通过最小可运行示例，逐步掌握理解 WireGuard 所需的密码学与协议基础。下面是一个**只列提纲、不展开细节**的学习路线。

## 1. Go 与运行环境基础
- Go 模块与 `go.mod` / `go.sum` 的作用  
- `package main` 与 `go run` / `go build`  
- 日志与调试输出：`fmt.Printf`、分步骤打印

## 2. 椭圆曲线 Diffie-Hellman（X25519）
- Curve25519 / X25519 的基本概念（只需知道用途）  
- 使用 `golang.org/x/crypto/curve25519`：  
  - 生成密钥对（私钥 32 字节、公钥）  
  - 使用 `ScalarMult` 计算共享密钥  
- 双方只交换公钥即可协商出相同共享密钥的流程  
- Demo：`00-helloworld-go/curve25519dh/main.go`

## 3. AEAD 对称加密（AES-GCM / ChaCha20-Poly1305）
- AEAD 的概念：同时提供**机密性 + 完整性**  
- AES-256-GCM 的基本使用方式（只关注 API）  
- ChaCha20-Poly1305 的基本使用方式（WireGuard 实际使用）  
- nonce 的含义与随机生成方式  
- 如何把共享密钥（X25519 输出）当成 AEAD 的密钥使用  
- Demo：在 `curve25519dh/main.go` 中用共享密钥加解密一条消息

## 4. HKDF：从一个秘密派生多把密钥
- IKM / salt / info / OKM 的直观理解  
- 为什么需要“从一个共享秘密派生多把用途不同的密钥”  
- 使用 `golang.org/x/crypto/hkdf` + `sha256` 派生固定长度 key  
- 使用不同 `info` 派生不同方向 / 不同用途的密钥  
- Demo：`00-helloworld-go/hkdfdemo/main.go`

## 5. 串联：从 DH 到会话密钥
- 流程总览：  
  1. X25519 计算共享密钥  
  2. 把共享密钥作为 HKDF 的 IKM  
  3. 用 HKDF 派生出：  
     - 发起方 → 响应方的 AEAD key  
     - 响应方 → 发起方的 AEAD key  
  4. 使用派生出的 key 进行 ChaCha20-Poly1305 加解密  
- 计划中的 Demo：`dh + hkdf + chacha` 一体化示例

## 6. 协议层：从积木到 Noise / WireGuard
- Noise 协议框架的直观认识（不做深度数学）：  
  - static key / ephemeral key / PSK 的角色  
  - 消息 1/2/3 的大致结构  
- WireGuard 使用的 Noise 模式（IKpsk2）  
- “链密钥（chaining key）”与会话密钥如何在握手中演变

## 7. 源码层：wireguard-go 的结构（预览提纲）
- 握手相关代码（handshake）：  
  - 使用 X25519 / HKDF / ChaCha20-Poly1305 的位置  
  - 如何保存 peer 的静态/会话密钥  
- 数据通道相关代码（device / tunnel）：  
  - 发包：上层 → 加密 → UDP  
  - 收包：UDP → 解密 → 上层  
- 状态机与 key rotation：  
  - 何时触发重新握手  
  - 如何同时维护旧/新两套会话密钥  

---

后续目录（`01-...`、`02-...` 等）可以在此提纲的基础上，分别深入：  
- 具体握手消息格式  
- peer / device 的内部结构  
- 重放防护与计数器窗口  
- 操作系统层面的接口与集成方式

