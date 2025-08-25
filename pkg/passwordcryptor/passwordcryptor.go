package passwordcryptor

import "golang.org/x/crypto/bcrypt"

const hashCost = 12

type PasswordCryptor struct {
}

func (pc PasswordCryptor) Crypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)

	return string(hash), err
}

func (pc PasswordCryptor) CheckHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(hash)) == nil
}
