package encrypt

import (
	"fmt"
	"io"
)

type EncryptType string

const (
	EncryptTypeZip EncryptType = "zip"
	EncryptTypeGPG EncryptType = "gpg"
)

type FileEncrypter interface {
	Encrypt(w io.Writer, r io.Reader) error
}

//TODO: Functional options patternでオプションを渡せるようにする
func CreateEncrypter(t EncryptType) (FileEncrypter, error) {
	switch t {
	case EncryptTypeZip:
		return &ZipEncrypter{}, nil
	case EncryptTypeGPG:
		return &PGPEncrypter{}, nil
	}

	return nil, fmt.Errorf("Unknown EncryptType: %s", t)
}
