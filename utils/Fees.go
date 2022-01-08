package utils

import "github.com/libsv/go-bt"

var Fees = []*bt.Fee{
	{
		FeeType: "standard",
		MiningFee: bt.FeeUnit{
			Satoshis: 500,
			Bytes:    1000,
		},
		RelayFee: bt.FeeUnit{
			Satoshis: 250,
			Bytes:    1000,
		},
	},
	{
		FeeType: "data",
		MiningFee: bt.FeeUnit{
			Satoshis: 500,
			Bytes:    1000,
		},
		RelayFee: bt.FeeUnit{
			Satoshis: 250,
			Bytes:    1000,
		},
	},
}
