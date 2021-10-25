package claim

type EncryptType string

const (
	EncryptTypeZip EncryptType = "zip"
	EncryptTypeGPG EncryptType = "gpg"
)

type Claim struct {
	Subject     string
	Name        string
	SlackID     string
	Email       string
	EncryptType EncryptType
}

