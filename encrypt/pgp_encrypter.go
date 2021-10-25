// WIP
package encrypt

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

// NOTE:  https://gist.github.com/r10r/1254039c940c426656e5d217216e0eec

type PGPEncrypter struct {
	privateKey *packet.PrivateKey
}

func (*PGPEncrypter) Encrypt(path string) error {
	return nil
}

func NewPGPEncrypter() (PGPEncrypter, error) {
	// ImplementMe
	return PGPEncrypter{}, nil
}

func decodeFromFile(path string) (*packet.PublicKey, error) {

	in, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer in.Close()

	// open ascii armored public key
	block, err := armor.Decode(in)
	if err != nil {
		return nil, err
	}

	if block.Type != openpgp.PublicKeyType {
		return nil, fmt.Errorf("Invalid public key file")
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		return nil, err
	}
	key, ok := pkt.(*packet.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Invalid public key file")
	}
	return key, nil
}

func decodePublicKey(rawKey string) (*packet.PublicKey, error) {
	// open ascii armored public key
	block, err := armor.Decode(strings.NewReader(rawKey))
	if err != nil {
		return nil, err
	}

	if block.Type != openpgp.PublicKeyType {
		return nil, fmt.Errorf("Invalid public key file")
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		return nil, err
	}
	key, ok := pkt.(*packet.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Invalid public key file")
	}
	return key, nil
}

// func decodePublicKey(filename string) (*packet.PublicKey, error) {
//
// 	// open ascii armored public key
// 	in, err := os.Open(filename)
// 	defer in.Close()
// 	if err != nil {
// 		return nil, fmt.Errorf("Error: %s", "TODO")
// 	}
//
// 	block, err := armor.Decode(in)
//
// 	if block.Type != openpgp.PublicKeyType {
// 		return nil, fmt.Errorf("Error: %s", "TODO")
// 	}
//
// 	reader := packet.NewReader(block.Body)
// 	pkt, err := reader.Next()
//
// 	key, ok := pkt.(*packet.PublicKey)
// 	if !ok {
// 		return nil, fmt.Errorf("Error: %s", "TODO")
// 	}
// 	return key, nil
// }
