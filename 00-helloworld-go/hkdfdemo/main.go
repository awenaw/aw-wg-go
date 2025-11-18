// Package main 演示 HKDF 的基本用法，
// 模拟“用一段原始密钥材料派生出多把会话密钥”的过程。
// 这和 WireGuard 里用 DH 输出 + HKDF 派生不同用途的密钥是同一种模式。
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

func main() {
	fmt.Println("=== HKDF 示例：从一段密钥材料派生多把会话密钥 ===")

	// 步骤 1：准备一段“原始密钥材料”（IKM）
	// 在真实协议里，这通常是：
	// - DH 计算出来的共享秘密，或者
	// - 预共享密钥（PSK）
	var ikm [32]byte
	mustFillRandom(ikm[:])
	fmt.Printf("步骤 1：原始密钥材料 (IKM): %x\n\n", ikm)

	// 步骤 2：选择 salt
	// 在协议里可以是：
	// - 上一次握手的链密钥（chaining key）
	// - 或者某个固定常量（比如“Noise_WireGuard_1”之类）
	salt := []byte("wg-demo-hkdf-salt")
	fmt.Printf("步骤 2：salt（可视作上一轮状态）: %x (文本: %q)\n\n", salt, salt)

	// 步骤 3：使用不同的 info 字段派生出不同用途的密钥
	// info 通常用来区分“这把 key 用来干什么”，例如：
	// - "initiator->responder key"
	// - "responder->initiator key"
	infoKey1 := []byte("initiator->responder key")
	infoKey2 := []byte("responder->initiator key")

	key1 := mustHKDF(ikm[:], salt, infoKey1, 32)
	key2 := mustHKDF(ikm[:], salt, infoKey2, 32)

	fmt.Println("步骤 3：通过 HKDF 派生会话密钥：")
	fmt.Printf(" key1 (info = %q): %x\n", infoKey1, key1)
	fmt.Printf(" key2 (info = %q): %x\n\n", infoKey2, key2)

	// 步骤 4：演示“只改 info，派生出的 key 就完全不同”
	infoKey1b := []byte("initiator->responder key - alt")
	key1b := mustHKDF(ikm[:], salt, infoKey1b, 32)

	fmt.Println("步骤 4：更换 info 后重新派生：")
	fmt.Printf(" key1  (info = %q): %x\n", infoKey1, key1)
	fmt.Printf(" key1b (info = %q): %x\n", infoKey1b, key1b)
	fmt.Println(" key1 与 key1b 是否相等:", equalBytes(key1, key1b))
}

// mustHKDF 使用 HKDF-SHA256 从 ikm + salt 派生出 length 字节的 key。
// 出错时直接 panic，方便示例代码简化错误处理。
func mustHKDF(ikm, salt, info []byte, length int) []byte {
	h := hkdf.New(sha256.New, ikm, salt, info)
	okm := make([]byte, length)
	if _, err := io.ReadFull(h, okm); err != nil {
		panic(fmt.Errorf("hkdf read: %w", err))
	}
	return okm
}

// mustFillRandom 使用 crypto/rand 生成指定长度的随机字节。
func mustFillRandom(buf []byte) {
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("rand read: %w", err))
	}
}

// equalBytes 简单比较两个字节切片是否相等，用于打印结果。
func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
