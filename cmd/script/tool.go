package main

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func keygen() {
	// Generate a new private key.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	// Get the private key in its hexadecimal representation
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := fmt.Sprintf("0x%x", privateKeyBytes)
	fmt.Println("Private Key:", privateKeyHex)

	// Extract the public key in ECDSA format
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	// Encode the public key in uncompressed form prefixed with 0x04
	publicKeyBytes := elliptic.Marshal(crypto.S256(), publicKeyECDSA.X, publicKeyECDSA.Y)
	publicKeyHex := fmt.Sprintf("0x%x", publicKeyBytes[1:]) // 0x04 prefix is usually implied, not added by elliptic.Marshal
	fmt.Println("Public Key:", publicKeyHex)

	// Get the public address derived from the public key
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("Public Address:", address)
}

func decryptMsg() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your privatekey: ")
	privateKeyHex, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln("can't read privatekey: ", err)
	}
	privateKeyHex = strings.ReplaceAll(privateKeyHex, "\n", "")
	privateKeyHex = privateKeyHex[2:] // skip '0x'

	fmt.Print("Enter hashed_msg: ")
	encryptedMessageHex, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln("can't read privatekey: ", err)
	}
	encryptedMessageHex = strings.ReplaceAll(encryptedMessageHex, "\n", "")

	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}
	fmt.Println(len(privateKeyBytes))
	privateKeyECDSA, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to create ECDSA private key: %v", err)
	}

	privateKeyECIES := ecies.ImportECDSA(privateKeyECDSA)

	encryptedMessage, err := hex.DecodeString(encryptedMessageHex[2:]) // skip '0x'
	if err != nil {
		log.Fatalf("Failed to decode encrypted message: %v", err)
	}

	decryptedMessage, err := privateKeyECIES.Decrypt(encryptedMessage, nil, nil)
	if err != nil {
		log.Fatalf("Failed to decrypt message: %v", err)
	}

	fmt.Println("Decrypted Message:", string(decryptedMessage))
}

func main() {
	keygenCmd := flag.CommandLine.Bool("keygen", false, "Generate a new key pair")
	decryptCmd := flag.CommandLine.Bool("decrypt", false, "Decrypt a message")
	flag.Parse()

	if !*keygenCmd && !*decryptCmd {
		fmt.Println("Error: Please specify a command: -keygen or -decrypt")
		os.Exit(1)
	}

	if *keygenCmd {
		keygen()
	} else if *decryptCmd {
		decryptMsg()
	}
}
