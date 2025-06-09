package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

func readPassword() (pass_hash string) {
	var pass string
	if password_flag != nil && *password_flag != "" {
		pass = *password_flag
	} else {
		if non_interactive_flag != nil && !*non_interactive_flag {
			fmt.Println("Enter password: ")
		}

		b, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			// Check if the error is a Ctrl+C (SIGINT) interrupt
			if err.Error() == "interrupt" {
				fmt.Println("\nPassword entry interrupted. Exiting...")
				os.Exit(1) // Gracefully exit after Ctrl+C
			} else {
				// Handle other types of errors
				fmt.Println("Error reading password:", err)
				os.Exit(1) // Exit on any other error
			}
		}
		pass = string(b)
	}
	if pass == "" {
		fmt.Println("password cannot be blank")
		return // explicitly exit the program
	}
	hash := md5.Sum([]byte(pass))
	return hex.EncodeToString(hash[:])
}

/*
	 Credit to
		https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes

for the encryption code below
*/
func encrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case,
	// we add it as a prefix to the encrypted data.
	// The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

/*
	Credit to
		https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes

for the encryption code below
*/
func decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		if strings.Contains(err.Error(), "message authentication failed") {
			fmt.Println("wrong password")
			return
		} else {
			panic(err)
		}
	}

	return string(plaintext)
}
