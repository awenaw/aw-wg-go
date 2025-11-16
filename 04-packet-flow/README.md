# packet-flow

## 目的
- 逐帧跟踪数据包自 TUN 进入用户态、加密、封装、经 UDP 发送的全过程。
- 反向路径同样记录，并附带关键函数/文件指引。

## 原理简介
数据平面主要由以下路径组成：
1. `RoutineReadFromTUN` 读取 IP 包，匹配 Peer，排入加密队列。
2. `RoutineEncryption` 使用当前发送 Keypair 进行 ChaCha20-Poly1305 封装。
3. `RoutineWriteToUDP` 将封包发送至对端。接收路径反向执行并写回 TUN。
本目录将配合时序图、栈追踪帮助理解性能热点与背压点。
