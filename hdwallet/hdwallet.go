package hdwallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

var hardenedKeyStart = 0x80000000

type AccountInfo struct {
	name       string
	address    string
	privateKey string
	publicKey  string
}

type Wallet struct {
	mnemonic  string
	masterKey *hdkeychain.ExtendedKey
	seed      []byte
}

// NewMnemonic fixed to use the 128bits
func NewMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

// NewFromMnemonic returns a new wallet from a BIP-39 mnemonic.
func NewFromMnemonic(mnemonic string, passphrase string) (*Wallet, error) {
	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("mnemonic is invalid")
	}

	seed, err := newSeedFromMnemonic(mnemonic, passphrase)
	if err != nil {
		return nil, err
	}

	wallet, err := newWallet(seed)
	if err != nil {
		return nil, err
	}
	wallet.mnemonic = mnemonic

	return wallet, nil
}
func NewFromPrivateKey(key string, passphrase string) (*Wallet, error) {
	var wallet Wallet
	masterKey, err := hdkeychain.NewKeyFromString(key)
	if err != nil {
		return nil, err
	}
	wallet.masterKey = masterKey

	return &wallet, nil
}

func (w *Wallet) GenerateAccount(name string, idx int) (*AccountInfo, error) {
	var err error
	var path []int
	switch name {
	case "BTC":
		path = []int{hardenedKeyStart + 44, hardenedKeyStart, hardenedKeyStart, 0, idx} // m/44'/0'/0'/0/0
	case "USDT":
		path = []int{hardenedKeyStart + 44, hardenedKeyStart + 200, hardenedKeyStart, 0, idx} // m/44'/200'/0'/0/0
	case "ETH":
		path = []int{hardenedKeyStart + 44, hardenedKeyStart + 60, hardenedKeyStart, 0, idx} // m/44'/60'/0'/0/0
	case "DDMX":
		path = []int{hardenedKeyStart + 44, hardenedKeyStart + 338, hardenedKeyStart, 0, idx} // m/44'/338'/0'/0/0
	case "GPE":
		path = []int{hardenedKeyStart + 44, hardenedKeyStart + 608, hardenedKeyStart, 0, idx} // m/44'/608'/0'/0/0
	}

	key := w.masterKey
	for _, n := range path {
		key, err = key.Child(uint32(n))
		if err != nil {
			return nil, err
		}
	}

	var acc AccountInfo
	switch name {
	case "BTC", "USDT":
		acc, err = newBitcoinSeriesAccount(key)
	case "ETH", "DDMX", "GPE":
		acc, err = newEthereumSeriesAccount(key)
	}
	if name == "GPE" {
		addr := acc.address
		acc.address = addr[0:2] + "gpe" + addr[2:]
	}
	return &AccountInfo{
		name:       name,
		address:    acc.address,
		privateKey: acc.privateKey,
		publicKey:  acc.publicKey,
	}, nil
}

func (ac *AccountInfo) Address() string {
	return ac.address
}
func (ac *AccountInfo) PrivateKey() string {
	return ac.privateKey
}
func (ac *AccountInfo) PublicKey() string {
	return ac.publicKey
}

func newWallet(seed []byte) (*Wallet, error) {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		masterKey: masterKey,
		seed:      seed,
	}, nil
}
func newSeedFromMnemonic(mnemonic string, passphrase string) ([]byte, error) {
	// NewSeedFromMnemonic returns a BIP-39 seed based on a BIP-39 mnemonic.
	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}

	return bip39.NewSeedWithErrorChecking(mnemonic, passphrase)
}
func newBitcoinSeriesAccount(key *hdkeychain.ExtendedKey) (AccountInfo, error) {
	var rtn AccountInfo

	addr, _ := key.Address(&chaincfg.MainNetParams)

	privateKey, err := key.ECPrivKey()
	if err != nil {
		return rtn, err
	}

	publicKey, err := key.ECPubKey()
	if err != nil {
		return rtn, err
	}
	publicKeyBytes := publicKey.SerializeCompressed()

	wif, _ := btcutil.NewWIF(privateKey, &chaincfg.MainNetParams, true)

	return AccountInfo{
		address:    addr.String(),
		privateKey: wif.String(),
		publicKey:  hex.EncodeToString(publicKeyBytes),
	}, nil
}
func newEthereumSeriesAccount(key *hdkeychain.ExtendedKey) (AccountInfo, error) {
	var rtn AccountInfo

	privateKey, err := key.ECPrivKey()
	if err != nil {
		return rtn, err
	}
	privateKeyECDSA := privateKey.ToECDSA()
	privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return rtn, errors.New("{\"failed ot get public key\"}")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	addr := crypto.PubkeyToAddress(*publicKeyECDSA)

	return AccountInfo{
		address:    addr.Hex(),
		privateKey: hexutil.Encode(privateKeyBytes)[2:],
		publicKey:  hexutil.Encode(publicKeyBytes)[2:],
	}, nil
}
