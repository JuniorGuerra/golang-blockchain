package wallet

import (
	"bytes"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadJsonWallet()

	return &wallets, err
}

func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses

}

func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := string(wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// Deprecated: use LoadJsonWallet
func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		fmt.Println("Wallet file not exists")
		return err
	}

	var wallets Wallets

	fileContent, err := os.ReadFile(walletFile)

	if err != nil {
		return err
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))

	err = decoder.Decode(&wallets)

	if err != nil {
		return err
	}

	ws.Wallets = wallets.Wallets

	return nil
}

func (ws *Wallets) LoadJsonWallet() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		fmt.Println("Wallet file not exists")
		return err
	}

	var wallets Wallets

	fileContent, err := os.ReadFile(walletFile)

	if err != nil {
		return err
	}

	walletsMaps := make(map[string]map[string][]byte)
	err = json.Unmarshal(fileContent, &walletsMaps)
	if err != nil {
		return err
	}

	wallets.Wallets = make(map[string]*Wallet)

	for address, keys := range walletsMaps {

		privateKey := keys["private"]

		privateKeyByte, err := x509.ParseECPrivateKey(privateKey)
		if err != nil {
			log.Panic(err)
		}

		publicKey := keys["public"]
		wallets.Wallets[address] = &Wallet{
			PrivateKey: *privateKeyByte,
			PublicKey:  publicKey,
		}
	}

	ws.Wallets = wallets.Wallets

	return nil
}

// Deprecated: use SaveJsonWallet
func (ws *Wallets) SaveFile() {
	var content bytes.Buffer
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, content.Bytes(), 0644)

	if err != nil {
		log.Panic(err)
	}

}

func (ws Wallets) SaveJsonWallet() {
	walletsMaps := make(map[string]map[string][]byte)
	for address, wallet := range ws.Wallets {
		privateKey, err := x509.MarshalECPrivateKey(&wallet.PrivateKey)
		if err != nil {
			log.Panic(err)
		}
		walletsMaps[address] = map[string][]byte{
			"private": privateKey,
			"public":  wallet.PublicKey,
		}
	}

	jsonData, err := json.MarshalIndent(walletsMaps, "", "  ")
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, jsonData, 0644)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Wallets saved in JSON format to", walletFile)
}
