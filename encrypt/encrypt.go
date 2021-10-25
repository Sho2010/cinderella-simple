package encrypt

type FileEncrypter interface {
	Encrypt(path string, src []byte) error
}
