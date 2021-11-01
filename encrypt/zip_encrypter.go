package encrypt

import (
	"fmt"
	"io"
	"log"

	"github.com/sethvargo/go-password/password"
	"github.com/yeka/zip"
)

type ZipEncrypter struct {
	Password string
}

func (enc *ZipEncrypter) Encrypt(w io.Writer, r io.Reader) error {

	if len(enc.Password) == 0 {
		fmt.Println("!!!Zip password blank")
		passwd, err := password.Generate(32, 10, 10, false, false)
		if err != nil {
			return fmt.Errorf("password generate fail: %w", err)
		}
		enc.Password = passwd
	}

	zipw := zip.NewWriter(w)
	defer zipw.Close()

	w, err := zipw.Encrypt("cinderella_config", enc.Password, zip.AES256Encryption)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(w, r)

	if err != nil {
		log.Fatal(err)
	}
	zipw.Flush()
	return nil
}
