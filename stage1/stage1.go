package stage1

import (
	"context"
	"crypto/ecdsa"
	"dapp_stage1/stage1/counter"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetBlockInfo 任务 1：区块链读写 任务目标
func GetBlockInfo() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/ee21fca4bf0c434bb94402e53f4def84")
	if err != nil {
		log.Fatal(err)
	}
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("blockHash:", block.Hash())
	fmt.Println("blockTime:", block.Time())
	fmt.Println("blockNumber:", block.Transactions().Len())
}

// Transaction 任务 2：合约代码生成 任务目标
func Transaction() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/ee21fca4bf0c434bb94402e53f4def84")
	pk := "793cceee38dd634db97a45ebee6dd112e18d6e2b68c42023eb9a25d21060e5e3"
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	fromAddress := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))
	value := big.NewInt(10000000000000000) // 0.01eth
	tipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	toAddress := common.HexToAddress("0x8fBB7D0e008301e72DE134D45E94cD50DC6DA23b")
	chainId, _ := client.ChainID(context.Background())
	nonce, err := client.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	tx := types.NewTx(&types.DynamicFeeTx{
		GasTipCap: tipCap,
		GasFeeCap: gasPrice,
		Nonce:     nonce,
		To:        &toAddress,
		Value:     value,
		Gas:       uint64(21000),
	})
	signTx, err := types.SignTx(tx, types.NewLondonSigner(chainId), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("signTx:", signTx.Hash().Hex())
}

func Increment() {
	contractAddress := "0xFdDB8a61339208EEDda482DE873160dea9b396a3"
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/ee21fca4bf0c434bb94402e53f4def84")
	if err != nil {
		log.Fatal(err)
	}
	contract, err := counter.NewCounter(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	pk := "793cceee38dd634db97a45ebee6dd112e18d6e2b68c42023eb9a25d21060e5e3"
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatal(err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	tx, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}
	transaction, err := contract.Increment(tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("transaction:", transaction.Hash().Hex())
}

func GetCount() {
	contractAddress := "0xFdDB8a61339208EEDda482DE873160dea9b396a3"
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/ee21fca4bf0c434bb94402e53f4def84")
	if err != nil {
		log.Fatal(err)
	}
	contract, err := counter.NewCounter(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	optCall := &bind.CallOpts{Context: context.Background()}
	count, err := contract.Count(optCall)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("count:", count)
}
