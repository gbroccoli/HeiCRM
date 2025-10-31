package password

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// Password : base structure Password Manager
type Password struct{}

// GenerateHash : func for generate password hash
func (p *Password) GenerateHash(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// CheckHash : func for check hash password
func (p *Password) CheckHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (p *Password) GeneratePassword() string {
	length := 24

	lowercase := "abcdefghijklmnopqrstuvwxyz"
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	symbols := "!@#$%^&*()-_=+"

	allChars := lowercase + uppercase + digits + symbols

	password := make([]byte, length)
	password[0] = lowercase[mustRandomInt(len(lowercase))]
	password[1] = lowercase[mustRandomInt(len(lowercase))]
	password[2] = lowercase[mustRandomInt(len(lowercase))]
	password[3] = lowercase[mustRandomInt(len(lowercase))]

	for i := 4; i < length; i++ {
		password[i] = allChars[mustRandomInt(len(allChars))]
	}

	for i := length - 1; i > 0; i-- {
		j := mustRandomInt(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password)
}

func mustRandomInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}
