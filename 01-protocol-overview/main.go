package main

import (
	"fmt"
	"strings"
)

// Peer 表示 WireGuard 中的一个节点，只保留静态公钥用来描述 NoiseIK 里的身份绑定。
type Peer struct {
	Name      string
	PublicKey string
}

// Message 抽象握手/传输阶段的消息内容；用 map 记录字段便于拓展。
type Message struct {
	Kind     string
	Sender   Peer
	Receiver Peer
	Fields   map[string]string
}

// Simulation 通过“脚本”串联一次 NoiseIK 流程，并保存文字版 Transcript。
type Simulation struct {
	Initiator      Peer
	Responder      Peer
	Transcript     []string
	keyDerivations []string
}

// aw-newSimulation 创建两个 peer；你只要知道 Initiator/Responder 已经准备好了。
func newSimulation() *Simulation {
	return &Simulation{
		Initiator: Peer{Name: "Initiator", PublicKey: "init-pub-AAAAAAAAAAAAAAAAAAAA"},
		Responder: Peer{Name: "Responder", PublicKey: "resp-pub-BBBBBBBBBBBBBBBBBBBB"},
	}
}

// logf/derive 负责追踪旁白与密钥派生顺序，方便和 README 对照。
func (s *Simulation) logf(format string, args ...any) {
	s.Transcript = append(s.Transcript, fmt.Sprintf(format, args...))
}

func (s *Simulation) derive(label, material string) {
	s.keyDerivations = append(s.keyDerivations, fmt.Sprintf("[%s] => %s", label, material))
}

// send/reply 对应 Initiator 与 Responder 发出的消息。
func (s *Simulation) send(kind string, fields map[string]string) {
	msg := Message{
		Kind:     kind,
		Sender:   s.Initiator,
		Receiver: s.Responder,
		Fields:   fields,
	}
	s.describeMessage(msg)
}

func (s *Simulation) reply(kind string, fields map[string]string) {
	msg := Message{
		Kind:     kind,
		Sender:   s.Responder,
		Receiver: s.Initiator,
		Fields:   fields,
	}
	s.describeMessage(msg)
}

// describeMessage 将字段拼成“Sender -> Receiver | Kind | key=value”格式，便于阅读。
func (s *Simulation) describeMessage(msg Message) {
	var fieldParts []string
	for k, v := range msg.Fields {
		fieldParts = append(fieldParts, fmt.Sprintf("%s=%s", k, v))
	}
	s.logf("%s -> %s | %s | %s", msg.Sender.Name, msg.Receiver.Name, msg.Kind, strings.Join(fieldParts, ", "))
}

// run 展示 NoiseIK 的三步流程，并同步列出 CK/SK 的派生顺序。
func (s *Simulation) run() {
	s.logf("=== NoiseIK 握手/传输流程示例 ===")

	// Step 1: Initiator 生成临时密钥并封装第一条握手消息。
	s.logf("1) Initiator 生成 Ephemeral 公钥并组装 Handshake Initiation")
	s.derive("CK0", "Hash(ResponderStatic || Prologue)")
	s.derive("CK1", "MixHash(EphemeralInitiator)")
	s.send("HandshakeInitiation", map[string]string{
		"ephemeral": "init-ephemeral-1234",
		"static":    "Enc(ResponderPub, InitiatorStatic)",
		"timestamp": "T0",
	})

	// Step 2: Responder 验证消息并生成返回握手。
	s.logf("2) Responder 验证 Initiation，派生会话密钥")
	s.derive("CK2", "MixKey(ResponderStatic, EphemeralInit)")
	s.derive("SKi", "KDF2(CK2)")
	s.reply("HandshakeResponse", map[string]string{
		"ephemeral": "resp-ephemeral-5678",
		"empty":     "MACs",
	})

	// Step 3: 双方获得对称密钥后转入数据平面，示例一条数据与 Keepalive。
	s.logf("3) 双方完成对称密钥并切换到数据平面")
	s.derive("CK3", "MixKey(EphemeralResp, EphemeralInit)")
	s.derive("SKr", "KDF2(CK3)")
	s.send("TransportData", map[string]string{
		"counter": "0",
		"payload": "IP packet (encrypted)",
	})
	s.reply("TransportData", map[string]string{
		"counter": "0",
		"payload": "Keepalive",
	})

	// 展示派生轨迹，帮助联想到 wireguard-go 中 keypair 切换的依据。
	s.logf("=== 关键派生 ===")
	for _, entry := range s.keyDerivations {
		s.logf(entry)
	}
}

func main() {
	// main 只负责运行脚本并打印 Transcript，便于和 README 文字对照。
	sim := newSimulation()
	sim.run()
	for _, line := range sim.Transcript {
		fmt.Println(line)
	}
}
