// Package main 演示一次使用 Curve25519 的
// Diffie-Hellman（DH）密钥交换过程。
package main

import (
	"crypto/rand"
	"fmt"

	// 第三方密码学库，提供 Curve25519 椭圆曲线相关运算。
	"golang.org/x/crypto/curve25519"
)

func main() {
	// 1) 生成 Initiator（发起者）的密钥对
	iPriv, iPub := mustGenerateKeyPair()

	// 2) 生成 Responder（响应者）的密钥对
	rPriv, rPub := mustGenerateKeyPair()

	// 3) 双方使用“自己的私钥 + 对方的公钥”计算共享密钥
	var dhInitiator, dhResponder [32]byte
	curve25519.ScalarMult(&dhInitiator, &iPriv, &rPub)
	curve25519.ScalarMult(&dhResponder, &rPriv, &iPub)

	// 4) 打印两边算出的共享秘密，理论上应该完全一致
	fmt.Printf("DH from Initiator: %x\n", dhInitiator)
	fmt.Printf("DH from Responder: %x\n", dhResponder)
	fmt.Println("Equal:", dhInitiator == dhResponder)
}

// mustGenerateKeyPair 生成一对 Curve25519 密钥对：
// priv 为私钥，pub 为公钥。
// 如果随机数生成失败，直接 panic。
func mustGenerateKeyPair() ([32]byte, [32]byte) {
	var priv, pub [32]byte

	// 使用安全随机数填充 32 字节私钥
	mustFillRandom(priv[:])

	// ScalarBaseMult 使用 Curve25519 的基点生成公钥
	curve25519.ScalarBaseMult(&pub, &priv)
	return priv, pub
}

// mustFillRandom 使用 crypto/rand 生成指定长度的随机字节。
// 出现错误时同样直接 panic，方便示例代码简化错误处理。
func mustFillRandom(buf []byte) {
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("rand read: %w", err))
	}
}
