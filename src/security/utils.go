package security

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"golang.org/x/crypto/bcrypt"
)

// EncryptData шифрует данные с использованием PGP.
func EncryptData(data string, publicKey string) (string, error) {
	// Загрузка публичного ключа
	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(publicKey))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	writer, err := armor.Encode(&buf, "PGP MESSAGE", nil)
	if err != nil {
		return "", err
	}

	plaintext, err := openpgp.Encrypt(writer, entityList, nil, nil, nil)
	if err != nil {
		return "", err
	}
	defer plaintext.Close()

	if _, err := plaintext.Write([]byte(data)); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// GenerateHMAC генерирует HMAC для данных.
func GenerateHMAC(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// HashCVV хеширует CVV с использованием bcrypt.
func HashCVV(cvv string) (string, error) {
	hashedCVV, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedCVV), nil
}

// IsValidCardNumber проверяет, является ли номер карты валидным по алгоритму Луна.
func IsValidCardNumber(number string) bool {
	var sum int
	shouldDouble := false

	// Проходим по цифрам номера карты с конца
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0') // Преобразуем символ в цифру

		if shouldDouble {
			digit *= 2
			if digit > 9 {
				digit -= 9 // Если результат больше 9, вычитаем 9
			}
		}

		sum += digit
		shouldDouble = !shouldDouble // Переключаем флаг
	}

	return sum%10 == 0 // Номер валиден, если сумма делится на 10
}
