const bsv = require('bsv')

// First we need a private key. This will be a new random private key each time we run.
const privateKey = bsv.PrivateKey('3421d46da4b9b2f0e7f00235edde586696d5c8bd982a0112b54213ca3bbe49f4')
console.log(`Private key:             ${privateKey.toString()}`)

// Get the public key associated with this private key.  Print this out in both uncompressed and compressed formats.
const publicKey = privateKey.publicKey
console.log(`Compressed public key:   ${publicKey.toString()}`)

const publicKeyHash = bsv.crypto.Hash.sha256ripemd160(publicKey.toBuffer())
console.log(`Public key hash:         ${publicKeyHash.toString('hex')}`)


// Non-mainnet address for use below.
const address = bsv.Address.fromPublicKey(bsv.PrivateKey().publicKey, 'testnet')

console.log(`Testnet public key hash: ${address.toObject().hash}`)
console.log(`Testnet address:         ${address.toString()}`)

// TODO - the following UTXO needs to be generated like in the Golang example.
const utxo = {
  txid: "5c2120cb72e197d34767a45e9bbd1f3d308809ad09d8b9547002781a070ab664",
  vout: 1,
  scriptPubKey: '76a914a279c46801c5691bb27aca2ac79c6c173fde6b7888ac',
  satoshis: 1e8
}

const tx = bsv.Transaction()
  .from(utxo)
  .to(address, utxo.satoshis / 2)
  .change(address)
  
// The following line is how to add an OP_RETURN output to the transaction.
// tx.addData(['Hello world', 'text/plain'])

console.log(`Unsigned transaction:    ${tx.toString()}`)

tx.sign(privateKey)

console.log(`Signed transaction:      ${tx.toString()}`)
