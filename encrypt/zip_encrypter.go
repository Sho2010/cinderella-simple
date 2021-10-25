package encrypt

import (
	"bytes"
	"github.com/yeka/zip"
	"io"
	"log"
	"os"
)

type ZipEncrypter struct {
	Password string
}

func (enc *ZipEncrypter) Encrypt(dst string, src []byte) error {
	fzip, err := os.Create(dst)
	if err != nil {
		log.Fatalln(err)
	}
	zipw := zip.NewWriter(fzip)
	defer zipw.Close()

	w, err := zipw.Encrypt("cinderella_config", enc.Password, zip.AES256Encryption)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(w, bytes.NewReader(src))
	if err != nil {
		log.Fatal(err)
	}
	zipw.Flush()
	return nil
}
