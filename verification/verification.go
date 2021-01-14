package verification

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"

	// "github.com/decred/dcrd/dcrec/secp256k1"
	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/common/hexutil"
	// "github.com/ethereum/go-ethereum/crypto"
	// "golang.org/x/crypto/sha3"

	b64 "encoding/base64"
)

type Verify struct {
	pubKey  *secp256k1.PublicKey
	privKey *secp256k1.PrivateKey
}
type AccountInfo struct {
	privateKey string
	publicKey  string
}

func ImportPubKey(key string) (*Verify, error) {
	pubKeyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	pubKey, err := secp256k1.ParsePubKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	// var verify Verify
	// verify.pubKey = pubKey

	return &Verify{
		pubKey: pubKey,
	}, nil
}

func ImportPrivKey(key string) (*Verify, error) {
	privKeyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	privKey, pubKey := secp256k1.PrivKeyFromBytes(privKeyBytes)

	// var loginVerify LoginVerify
	// loginVerify.privKey = privKey
	// loginVerify.pubKey = pubKey

	return &Verify{
		privKey: privKey,
		pubKey:  pubKey,
	}, nil
}

func (v *Verify) Encrypt(plaintext string) (string, error) {
	if v.pubKey == nil {
		err := errors.New("public key required.")
		return "", err
	}

	ciphertext, err := secp256k1.Encrypt(v.pubKey, []byte(plaintext))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ciphertext), nil
}
func (v *Verify) EncryptWithCheck(account string, plaintext string) (string, error) {
	if v.pubKey == nil {
		err := errors.New("public key required.")
		return "", err
	}

	err := v.VerifyAccountWithPubKey(account)
	if err != nil {
		return "", err
	}

	ciphertext, err := secp256k1.Encrypt(v.pubKey, []byte(plaintext))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ciphertext), nil
}
func (v *Verify) EncryptBase64(plaintext string) (string, error) {
	if v.pubKey == nil {
		err := errors.New("public key required.")
		return "", err
	}

	ciphertext, err := secp256k1.Encrypt(v.pubKey, []byte(plaintext))
	if err != nil {
		return "", err
	}

	return b64.StdEncoding.EncodeToString([]byte(ciphertext)), nil
}

func (v *Verify) Decrypt(ct string) (string, error) {
	if v.pubKey == nil {
		err := errors.New("private key required.")
		return "", err
	}

	ciphertext, err := hex.DecodeString(ct)
	if err != nil {
		return "", err
	}

	plaintext, err := secp256k1.Decrypt(v.privKey, ciphertext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
func (v *Verify) DecryptWithCheck(account string, ct string) (string, error) {
	if v.pubKey == nil {
		err := errors.New("private key required.")
		return "", err
	}

	err := v.VerifyAccountWithPubKey(account)
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(ct)
	if err != nil {
		return "", err
	}

	plaintext, err := secp256k1.Decrypt(v.privKey, ciphertext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (v *Verify) VerifyAccountWithPubKey(account string) error {
	pubKeyBytes := v.pubKey.SerializeUncompressed()

	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubKeyBytes[1:]) // remove EC prefix 04

	buf := hash.Sum(nil)
	addr := common.HexToAddress(hex.EncodeToString(buf[12:]))
	// fmt.Println("addr:", strings.ToLower(addr.Hex()))

	if strings.Compare(strings.ToLower(account), strings.ToLower(addr.Hex())) != 0 { // account=addr.Hex() is 0
		return errors.New("verify failure.")
	}

	return nil
}

// func VerifyAccountPubkey(account, pubkey string) error {
// 	// 驗證公鑰與帳號是否一致
// 	pubKeyBytes := pubKey.SerializeUncompressed()

// 	hash := sha3.NewLegacyKeccak256()
// 	hash.Write(pubKeyBytes[1:]) // remove EC prefix 04

// 	buf := hash.Sum(nil)
// 	addr := common.HexToAddress(hex.EncodeToString(buf[12:]))
// 	// fmt.Println("addr:", strings.ToLower(addr.Hex()))

// 	if strings.Compare(strings.ToLower(account), strings.ToLower(addr.Hex())) != 0 { // account=addr.Hex() is 0
// 		return errors.New("verify failure.")
// 	}

// 	return nil
// }

func GenerateKeyPair() (*AccountInfo, error) {
	secKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, errors.New("Generate Private Key failure.")
	}

	privateKeyECDSA := secKey.ToECDSA()
	privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed ot get public key")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	return &AccountInfo{
		privateKey: hexutil.Encode(privateKeyBytes)[2:],
		publicKey:  hexutil.Encode(publicKeyBytes)[2:],
	}, nil
}
func (ac *AccountInfo) PrivateKey() string {
	return ac.privateKey
}
func (ac *AccountInfo) PublicKey() string {
	return ac.publicKey
}
