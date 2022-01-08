package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/libsv/go-bt"
	"github.com/libsv/go-bt/bscript"
)

type Utxo struct {
	TxID          string `json:"txid"`
	Vout          uint32 `json:"vout"`
	LockingScript string `json:"lockingScript"`
	Satoshis      uint64 `json:"satoshis"`
}

func (u *Utxo) String() string {
	b, err := json.MarshalIndent(&u, "", "  ")
	if err != nil {
		return fmt.Sprintf("%+v", *u)
	}
	return string(b)
}

func GetUtxo(publicKey *bsvec.PublicKey) (*Utxo, error) {
	lockingScript, err := bscript.NewP2PKHFromPubKeyBytes(publicKey.SerializeCompressed())
	if err != nil {
		return nil, err
	}

	address, err := bscript.NewAddressFromPublicKey(publicKey, false) // false means "not mainnet"
	if err != nil {
		return nil, err
	}

	// Tell the node to send 1 BSV to our address.
	fundingTxID, err := Bitcoin.SendToAddress(address.AddressString, 1.0)
	if err != nil {
		return nil, err
	}

	// Ask the node for the raw hex of the transaction that it has just created.
	fundingTxHex, err := Bitcoin.GetRawTransactionHex(fundingTxID)
	if err != nil {
		return nil, err
	}

	// Parse the raw hex to create a fundingTx object.
	fundingTx, err := bt.NewTxFromString(*fundingTxHex)
	if err != nil {
		return nil, err
	}

	// Find the output that we was paid to our address...
	var vout int = -1

	for i, out := range fundingTx.Outputs {
		if out.LockingScript.ToString() == lockingScript.ToString() {
			vout = i
			break
		}
	}

	if vout == -1 {
		return nil, errors.New("Did not find a UTXO")
	}

	return &Utxo{
		TxID:          fundingTxID,
		Vout:          uint32(vout),
		LockingScript: lockingScript.ToString(),
		Satoshis:      fundingTx.Outputs[vout].Satoshis,
	}, nil
}
