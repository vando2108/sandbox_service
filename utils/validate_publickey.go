package utils

import (
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
