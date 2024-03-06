package service

import (
	"supernet.tools/tcp-proxy-server/config"
)

// General purpose encryptor for data transfer manipulations...
type Encryptor interface {
	// In place method, will encrypt input data
	Encrypt(data []byte)

	// In place method will decrypt input data
	Decrypt(data []byte)
}

// Provides dummy implementation of Encryptor
func NewDummyEncryptor(conf *config.AppConf) Encryptor {
	return &dummyEncryptor{
		cipher: []byte(conf.SecretKey),
	}
}

// dummy encryptor using XOR operator
type dummyEncryptor struct {
	cipher []byte
}

func (enc *dummyEncryptor) Encrypt(data []byte) {
	enc.encryptDecrypt(data)
}

func (enc *dummyEncryptor) Decrypt(data []byte) {
	enc.encryptDecrypt(data)
}

func (enc *dummyEncryptor) encryptDecrypt(data []byte) {
	keyLen := len(enc.cipher)
	for i := 0; i < len(data); i++ {
		data[i] = data[i] ^ enc.cipher[i%keyLen]
	}
}
