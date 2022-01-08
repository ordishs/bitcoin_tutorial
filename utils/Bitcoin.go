package utils

import "github.com/ordishs/go-bitcoin"

var Bitcoin *bitcoin.Bitcoind

func init() {
	var err error
	// Connect to local bitcoin regtest node.
	Bitcoin, err = bitcoin.New("localhost", 18332, "bitcoin", "bitcoin", false)
	if err != nil {
		panic(err)
	}
}
