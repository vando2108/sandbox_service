package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strings"
)

func ValidatePublickey(pubkey string) error {
	// mock
	pubkey = strings.ReplaceAll(pubkey, "0x", "")
	if len(pubkey) != 128 {
		return fmt.Errorf("Invalid key length.")
	}

	return nil
}

func PublickeyToString(publickey *ecdsa.PublicKey) string {
	x := publickey.X.Bytes()
	y := publickey.Y.Bytes()
	xy := append(x, y...)
	return hex.EncodeToString(xy)
}
