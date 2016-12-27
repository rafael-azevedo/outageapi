package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type OracleDB struct {
	Username    string
	Password    string
	HostName    string
	Port        string
	ServiceName string
}

type EncryptOutput struct {
	Key      string
	Password string
}

type EnvKeys struct {
	OracleDB  string
	RethingDB string
}

func NewOracleDB(cfg *viper.Viper, e *EnvKeys) (*OracleDB, error) {
	var o *OracleDB
	e.OracleDB = os.Getenv("ORACLEDB")
	oracle := cfg.Sub("oracle")
	err := oracle.Unmarshal(&o)
	if err != nil {
		log.Printf("unable to decode config into oracle struct, %#v", err)
		return o, err
	}
	if e.OracleDB != "" && o.Password != "" {
		o.Decrypt(e)
	}
	return o, nil
}

// This function encrypts your password with the supplied key that can be decrypted by the API and other tools in this software
func Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	return ciphertext, nil
}

func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GenBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a  base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenString(s int) (string, error) {
	b, err := GenBytes(s)
	return base64.StdEncoding.EncodeToString(b), err
}

func DecodeKey64(s string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return decoded, err
	}
	keyHexString := hex.EncodeToString(decoded)
	return []byte(keyHexString), nil
}

func (e *EncryptOutput) EncryptWithKey64(key string, password string) error {
	e.Key = key
	decoded, err := DecodeKey64(key)
	if err != nil {
		return err
	}
	encrptedPass, err := Encrypt(decoded, []byte(password))
	if err != nil {
		return err
	}
	e.Password = hex.EncodeToString(encrptedPass)
	return nil
}

func (e *EncryptOutput) GenKeyAndEncrypt(password string) error {
	key, err := GenString(16)
	if err != nil {
		return err
	}
	e.Key = key
	decoded, err := DecodeKey64(key)
	if err != nil {
		return err
	}
	encrptedPass, err := Encrypt(decoded, []byte(password))
	if err != nil {
		return err
	}
	e.Password = hex.EncodeToString(encrptedPass)
	return nil
}

func (o *OracleDB) Decrypt(e *EnvKeys) error {
	key, err := DecodeKey64(e.OracleDB)
	if err != nil {
		return err
	}
	pbytes, err := hex.DecodeString(o.Password)
	if err != nil {
		return err
	}
	password, err := Decrypt(key, pbytes)
	if err != nil {
		return err
	}
	o.Password = string(password)
	return nil
}

//LastID returns the highest number from a list of files whose names are integers in OUTAGELOGDIR enviroment variable directory, if there is no log directory it creates it and returns 0
func LastID() (int, error) {
	hiFi := 0
	files, err := ioutil.ReadDir(os.Getenv("OUTAGELOGDIR"))
	if err != nil {
		err = errors.New("Could not read log directory")
		return hiFi, err
	}

	for _, file := range files {
		fName := strings.Split(file.Name(), ".")
		fNum, isNumeric := IsNumeric(fName[0])
		switch {
		case isNumeric:
			if fNum > hiFi {
				hiFi = fNum
			}
		}

	}
	return hiFi, nil
}

//IsNumeric returns the int from the string and a bool true if the string is numeric and false if it is not
func IsNumeric(s string) (int, bool) {
	i64, err := strconv.ParseInt(s, 10, 64)
	i := int(i64)
	return i, err == nil
}
