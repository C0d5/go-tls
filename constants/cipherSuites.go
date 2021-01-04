package constants

import (
	"crypto/sha512"
	"hash"
)

type CipherSuiteInfo struct {
	Code [2]byte
	HashingAlgorithm func() hash.Hash
}

type CipherSuitesConverter struct {
	CipherSuites map[string]CipherSuiteInfo
}

// TODO create struct to include algorithm
func MakeCipherSuites() CipherSuitesConverter {
	cipherSuitesConverter := CipherSuitesConverter{}
	cipherSuitesConverter.CipherSuites = map[string]CipherSuiteInfo{
		//"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256":   {0xcc, 0xa8},
		//"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256": {0xcc, 0xa9},
		//"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":         {0xc0, 0x2f},
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384": {
			Code: [2]byte{0xc0, 0x30},
			HashingAlgorithm: sha512.New384,
		},
		//"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":       {0xc0, 0x2b},
		//"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":       {0xc0, 0x2c},
		//"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":            {0xc0, 0x13},
		//"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":          {0xc0, 0x09},
		//"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":            {0xc0, 0x14},
		//"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":          {0xc0, 0x0a},
		//"TLS_RSA_WITH_AES_128_GCM_SHA256":               {0x00, 0x9c},
		//"TLS_RSA_WITH_AES_256_GCM_SHA384":               {0x00, 0x9d},
		//"TLS_RSA_WITH_AES_128_CBC_SHA":                  {0x00, 0x2f},
		//"TLS_RSA_WITH_AES_256_CBC_SHA":                  {0x00, 0x35},
		//"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":           {0xc0, 0x12},
		//"TLS_RSA_WITH_3DES_EDE_CBC_SHA":                 {0x00, 0x0a},
	}

	return cipherSuitesConverter
}

func (converter CipherSuitesConverter) GetByteCodeForSuite(suite string) [2]byte {
	return converter.CipherSuites[suite].Code
}

func (converter CipherSuitesConverter) GetSuiteForByteCode(suite [2]byte) string {
	for k, v := range converter.CipherSuites {
		if v.Code == suite {
			return k
		}
	}

	return ""
}

func (converter CipherSuitesConverter) GetSuiteInfoForByteCode(suite [2]byte) *CipherSuiteInfo {
	for _, v := range converter.CipherSuites {
		if v.Code == suite {
			return &v
		}
	}

	return nil
}

func (converter CipherSuitesConverter) GetSuiteByteCodes(suites []string) []byte {
	var suiteByteCodes []byte

	for _, s := range suites {
		code := converter.CipherSuites[s].Code
		suiteByteCodes = append(suiteByteCodes, code[:]...)
	}

	return suiteByteCodes
}

func (converter CipherSuitesConverter) GetAllSuites() []string {
	var suites []string

	for k := range converter.CipherSuites {
		suites = append(suites, k)
	}

	return suites
}

var GCipherSuites = MakeCipherSuites()
