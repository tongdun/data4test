package biz

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"regexp"
	"strings"
)

// RsaEncryptPassword 使用 RSA 公钥（PEM）对明文做 PKCS1_v1_5 加密并 Base64 输出，
// 与 JSEncrypt / Python rsa 库的 rsa_pwd 行为一致。
// publicKeyPem 支持 PKIX(SPKI) 与 PKCS1 两种公钥格式。
func RsaEncryptPassword(publicKeyPem, plaintext string) (encrypted string, err error) {
	block, _ := pem.Decode([]byte(publicKeyPem))
	if block == nil {
		err = E("error.rsa_encrypt_public_key_invalid")
		return
	}

	var pub *rsa.PublicKey
	pubAny, errParse := x509.ParsePKIXPublicKey(block.Bytes)
	if errParse == nil {
		var ok bool
		pub, ok = pubAny.(*rsa.PublicKey)
		if !ok {
			err = E("error.rsa_encrypt_public_key_invalid")
			return
		}
	} else {
		pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			err = E("error.rsa_encrypt_public_key_invalid")
			return
		}
	}

	ciphertext, errTmp := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(plaintext))
	if errTmp != nil {
		err = E("error.rsa_encrypt_failed")
		return
	}

	encrypted = base64.StdEncoding.EncodeToString(ciphertext)
	return
}

// GetRsaEncryptOutput 遍历 df.Action，对 type 为 `rsa_encrypt:明文:结果变量名` 的动作
// 计算加密串，返回 `map[结果变量名]加密串`。明文与公钥中 `{变量}` 占位符从 depOutVars 解析。
func (df DataFile) GetRsaEncryptOutput(depOutVars map[string][]interface{}) (rsaOut map[string]string, err error) {
	rsaOut = make(map[string]string)
	if len(df.Action) == 0 {
		return
	}

	varReg := regexp.MustCompile(`\{([-a-zA-Z0-9_]+)\}`)
	for _, item := range df.Action {
		if !strings.HasPrefix(item.Type, "rsa_encrypt") {
			continue
		}

		parts := strings.SplitN(item.Type, ":", 3)
		if len(parts) != 3 {
			err = E("error.rsa_encrypt_value_invalid")
			Logger.Error("%s", err)
			return
		}

		plaintext := parts[1]
		outputVar := parts[2]
		publicKey := Interface2Str(item.Value)
		if len(plaintext) == 0 || len(outputVar) == 0 || len(publicKey) == 0 {
			err = E("error.rsa_encrypt_value_invalid")
			Logger.Error("%s", err)
			return
		}

		plaintext = resolveRsaVar(plaintext, depOutVars, varReg)
		publicKey = resolveRsaVar(publicKey, depOutVars, varReg)

		encrypted, errTmp := RsaEncryptPassword(publicKey, plaintext)
		if errTmp != nil {
			err = errTmp
			Logger.Error("%s", err)
			return
		}
		rsaOut[outputVar] = encrypted
	}

	return
}

// resolveRsaVar 将字符串中 `{变量}` 占位符替换为依赖变量的值，字面量原样返回。
func resolveRsaVar(s string, depOutVars map[string][]interface{}, reg *regexp.Regexp) string {
	match := reg.FindAllStringSubmatch(s, -1)
	if len(match) == 0 {
		return s
	}
	for _, m := range match {
		raw := m[0]
		name := m[1]
		if vals, ok := depOutVars[name]; ok && len(vals) > 0 {
			s = strings.Replace(s, raw, Interface2Str(vals[0]), -1)
		}
	}
	return s
}
