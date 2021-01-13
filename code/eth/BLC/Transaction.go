package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

//交易管理

type Transaction struct {
	TxHash []byte //交易哈希

	Vins  []*TxInput  //输入列表
	Vouts []*TxOutput //输出列表
}

//实现coinbase交易
func NewCoinbaseTransaction(address string) *Transaction {

	txInput := &TxInput{[]byte{}, -1, "system reward"}
	txOutput := &TxOutput{10, address}

	txConinbase := &Transaction{
		nil,
		[]*TxInput{txInput},
		[]*TxOutput{txOutput},
	}

	txConinbase.HashTransaction()

	return txConinbase
}

//生成交易哈希（交易序列hua)
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panicf("serialize transaction failed %v\n", err)
	}

	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

//生成普通交易
func NewSimpleTansaction(from string, to string, amount int, bc *BlockChain, txs []*Transaction) *Transaction {
	var txInputs []*TxInput
	var txOutputs []*TxOutput

	//可花费UTXO
	money, spendableUTXOMap := bc.FindSpendableUTXO(from, amount, txs)
	//fmt.Printf("money : %d\n", money)

	//输入
	for txHash, indexArr := range spendableUTXOMap {
		txHashBytes, err := hex.DecodeString(txHash)
		if err != nil {
			log.Panicf("decode string to []byte failed %v\n", err)
		}
		//索引列表
		for _, index := range indexArr {
			txInput := &TxInput{txHashBytes, index, from}
			txInputs = append(txInputs, txInput)
		}
	}

	//输出
	txOutput := &TxOutput{amount, to}
	txOutputs = append(txOutputs, txOutput)
	//输出找零
	if amount < money {
		txOutput = &TxOutput{money - amount, from}
		txOutputs = append(txOutputs, txOutput)
	} else {
		log.Panicf("%s余额不足, 余额：%d, 查找的数量和：%d\n", from, amount, money)
	}

	tx := Transaction{nil, txInputs, txOutputs}
	tx.HashTransaction()
	return &tx

}

//判断指定的交易是否coinbase
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return tx.Vins[0].Vout == -1 && len(tx.TxHash) == 0
}
