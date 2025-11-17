// Package main 演示一次使用 Curve25519 的
// Diffie-Hellman（DH）密钥交换过程，并用共享密钥做对称加密
// （同时演示 AES-GCM 和 ChaCha20-Poly1305 两种 AEAD）。
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
	// 第三方密码学库，提供 Curve25519 椭圆曲线相关运算。
	"golang.org/x/crypto/curve25519"
)

func main() {
	fmt.Println("=== 步骤 1：生成双方的 Curve25519 密钥对 ===")

	// 1) 生成 Initiator（发起者）的密钥对
	iPriv, iPub := mustGenerateKeyPair()

	// 2) 生成 Responder（响应者）的密钥对
	rPriv, rPub := mustGenerateKeyPair()

	fmt.Printf("Initiator 公钥: %x\n", iPub)
	fmt.Printf("Responder 公钥: %x\n\n", rPub)

	fmt.Println("=== 步骤 2：双方计算共享密钥（Diffie-Hellman） ===")

	// 3) 双方使用“自己的私钥 + 对方的公钥”计算共享密钥
	var dhInitiator, dhResponder [32]byte
	curve25519.ScalarMult(&dhInitiator, &iPriv, &rPub)
	curve25519.ScalarMult(&dhResponder, &rPriv, &iPub)

	// 4) 打印两边算出的共享秘密，理论上应该完全一致
	fmt.Printf("Initiator 计算的共享密钥: %x\n", dhInitiator)
	fmt.Printf("Responder 计算的共享密钥: %x\n", dhResponder)
	fmt.Println("是否相等 (保证握手成功):", dhInitiator == dhResponder)

	fmt.Println("\n=== 步骤 3：使用共享密钥做 AES-256-GCM 加密 / 解密 ===")

	// 5) 用共享密钥做一次对称加密示例（AES-GCM）
	// 注意：真实协议里通常会先对 DH 输出做 KDF（比如 HKDF），
	// 这里为了示例简单，直接把 32 字节 DH 输出当作 AES-256 密钥使用。
	plaintext := []byte("hello from Initiator, 加密测试")

	fmt.Printf("原始明文: %s\n", plaintext)

	ciphertext, err := encryptWithSharedKey(&dhInitiator, plaintext)
	if err != nil {
		panic(fmt.Errorf("encrypt error: %w", err))
	}
	fmt.Printf("AES-GCM 密文 (hex): %x\n", ciphertext)

	// Responder 端拿自己的 dhResponder，解密得到明文
	decrypted, err := decryptWithSharedKey(&dhResponder, ciphertext)
	if err != nil {
		panic(fmt.Errorf("decrypt error: %w", err))
	}
	fmt.Printf("AES-GCM 解密得到: %s\n", string(decrypted))

	fmt.Println("\n=== 步骤 4：使用共享密钥做 ChaCha20-Poly1305 加密 / 解密 ===")

	// 6) 用共享密钥做一次 ChaCha20-Poly1305 加密 / 解密示例
	chachaCiphertext, err := encryptWithChaCha(&dhInitiator, plaintext)
	if err != nil {
		panic(fmt.Errorf("chacha encrypt error: %w", err))
	}
	fmt.Printf("ChaCha20-Poly1305 密文 (hex): %x\n", chachaCiphertext)

	chachaDecrypted, err := decryptWithChaCha(&dhResponder, chachaCiphertext)
	if err != nil {
		panic(fmt.Errorf("chacha decrypt error: %w", err))
	}
	fmt.Printf("ChaCha20-Poly1305 解密得到: %s\n", string(chachaDecrypted))
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

// encryptWithSharedKey 使用共享密钥（32 字节）做 ChaCha20-Poly1305 加密。
// 返回值中同样把 nonce+密文拼在一起，方便传输。
func encryptWithChaCha(shared *[32]byte, plaintext []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(shared[:])
	if err != nil {
		return nil, fmt.Errorf("new chacha20poly1305: %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	mustFillRandom(nonce)

	ciphertext := aead.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}

// decryptWithChaCha 使用共享密钥解密 encryptWithChaCha 的输出。
func decryptWithChaCha(shared *[32]byte, data []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(shared[:])
	if err != nil {
		return nil, fmt.Errorf("new chacha20poly1305: %w", err)
	}

	nonceSize := aead.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("chacha20poly1305 open: %w", err)
	}
	return plaintext, nil
}

// encryptWithSharedKey 使用共享密钥（32 字节）做 AES-256-GCM 加密。
// 返回值中已经把 nonce+密文拼在一起，方便传输。
func encryptWithSharedKey(shared *[32]byte, plaintext []byte) ([]byte, error) {
	// 1) 用共享密钥创建 AES-256 block
	block, err := aes.NewCipher(shared[:])
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	// 2) 包一层 GCM（带认证标签的模式）
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	// 3) 生成随机 nonce
	nonce := make([]byte, aesgcm.NonceSize())
	mustFillRandom(nonce)

	// 4) Seal: 输出 = nonce || ciphertext||tag
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}

// decryptWithSharedKey 使用共享密钥解密 encryptWithSharedKey 的输出。
func decryptWithSharedKey(shared *[32]byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(shared[:])
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("gcm open: %w", err)
	}
	return plaintext, nil
}
