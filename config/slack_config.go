package config

// NOTE: tokenは必ず環境変数を介して渡すことにするのでここには含めない

type SlackConfig struct {
	// Administrators are slack app administrators. they can accept/reject permission claim.
	Administrators []string
	// notify streaming channel
	NotifyChannel string
	// audit log streaming channel
	AuditChannel string
}
