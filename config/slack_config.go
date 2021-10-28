package config

// NOTE: tokenは必ず環境変数を介して渡すことにするのでここには含めない

type SlackConfig struct {
	NotifyChannel string
	AuditChannel  string
}
