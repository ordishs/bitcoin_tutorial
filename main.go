package main

import (
	"fmt"
	"os"
	"tutorial/utils"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/libsv/go-bt"
	"github.com/libsv/go-bt/bscript"
	"github.com/libsv/go-bt/crypto"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "long" {
		longMain()
	} else {
		shortMain()
	}
}

func longMain() {
	// First we need a private key. This will be a new random private key each time we run.
	privateKey, err := bsvec.NewPrivateKey(bsvec.S256())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Private key:             %x\n", privateKey.Serialize())

	// Generate the public key associated with this private key.  Print this out in both uncompressed and compressed formats.
	publicKey := privateKey.PubKey()
	fmt.Printf("Uncompressed public key: %x\n", publicKey.SerializeUncompressed())
	fmt.Printf("Compressed public key:   %x\n", publicKey.SerializeCompressed())

	publicKeyHash := crypto.Hash160(publicKey.SerializeCompressed())
	fmt.Printf("Public key hash:         %x\n\n", publicKeyHash)

	// Calculate the P2PKH locking script which we will use to identify outputs assigned to this private key.
	lockingScript, err := bscript.NewP2PKHFromPubKeyBytes(publicKey.SerializeCompressed())
	if err != nil {
		panic(err)
	}
	fmt.Printf("lockingScript:           %x\n", *lockingScript)

	// Print out the ASM for this.
	asm, err := lockingScript.ToASM()
	if err != nil {
		panic(err)
	}
	fmt.Printf("lockingScript (ASM):     %s\n", asm)

	// Calculate mainnet address to illustrate the difference between mainnet and non-mainnet.
	mainnetAddress, err := bscript.NewAddressFromPublicKey(publicKey, true) // true means "mainnet"
	if err != nil {
		panic(err)
	}
	fmt.Printf("Mainnet public key hash: %s\n", mainnetAddress.PublicKeyHash)
	fmt.Printf("Mainnet address:         %s\n\n", mainnetAddress.AddressString)

	// Non-mainnet address for use below.
	address, err := bscript.NewAddressFromPublicKey(publicKey, false) // false means "not mainnet"
	if err != nil {
		panic(err)
	}
	fmt.Printf("Testnet public key hash: %s\n", address.PublicKeyHash)
	fmt.Printf("Testnet address:         %s\n\n", address.AddressString)

	utxo, err := utils.GetUtxo(publicKey)
	if err != nil {
		panic(err)
	}

	fmt.Printf("UTXO:\n%s\n\n", utxo)

	// Create a new empty transaction.
	tx := bt.NewTx()

	// Add the UTXO that we are spending.
	if err := tx.From(utxo.TxID, utxo.Vout, utxo.LockingScript, utxo.Satoshis); err != nil {
		panic(err)
	}

	// Split the amount in half.
	amount := utxo.Satoshis / 2

	// Add an output that pays 0.5 BSV to our address.
	if err := tx.PayTo(address.AddressString, amount); err != nil {
		panic(err)
	}

	// Tell go-bt to calculate appropriate fees and pay the balance (the change).  For this example we will pay this to ourselves as well.
	if err := tx.ChangeToAddress(address.AddressString, utils.Fees); err != nil {
		panic(err)
	}

	fmt.Printf("Unsigned transaction:    %x\n\n", tx.ToBytes())

	// Sign the transaction with our private key.
	if _, err = tx.SignAuto(&bt.InternalSigner{PrivateKey: privateKey, SigHashFlag: 0}); err != nil {
		panic(err)
	}

	fmt.Printf("Signed transaction:      %x\n\n", tx.ToBytes())

	// Send the raw transaction to our node.
	txID, err := utils.Bitcoin.SendRawTransaction(tx.ToString())
	if err != nil {
		panic(err)
	}

	fmt.Printf("New transaction created: %s\n\n", txID)
}

func shortMain() {
	privateKey, _ := bsvec.NewPrivateKey(bsvec.S256())
	publicKey := privateKey.PubKey()

	utxo, _ := utils.GetUtxo(publicKey)

	address, _ := bscript.NewAddressFromPublicKey(publicKey, false) // false means "not mainnet"

	tx := bt.NewTx()
	_ = tx.From(utxo.TxID, utxo.Vout, utxo.LockingScript, utxo.Satoshis)

	amount := utxo.Satoshis / 2

	_ = tx.PayTo(address.AddressString, amount)

	_ = tx.ChangeToAddress(address.AddressString, utils.Fees)

	_, _ = tx.SignAuto(&bt.InternalSigner{PrivateKey: privateKey, SigHashFlag: 0})

	txID, _ := utils.Bitcoin.SendRawTransaction(tx.ToString())

	fmt.Printf("New transaction created: %s\n\n", txID)
}
