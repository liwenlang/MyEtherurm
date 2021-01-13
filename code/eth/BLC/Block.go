package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	TimeStamp     int64          //区块时间戳
	Hash          []byte         //当前哈希
	PrevBlockHash []byte         //上一个区块哈希
	Height        int64          //区块高度
	Txs           []*Transaction //交易列表
	Nonce         int64          //随机数

}

//新建区块
func NewBlock(height int64, prevBlockHash []byte, txs []*Transaction) *Block {
	var block Block
	block = Block{
		TimeStamp:     time.Now().Unix(),
		Hash:          nil,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Txs:           txs,
	}

	//生成Hash
	//block.SetHash()

	//通过POW生成Hash
	pow := NewProofOfWork(&block)
	//执行工作量证明算法
	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return &block
}

//计算区块哈希
func (block *Block) SetHash() {

	timeStampBytes := IntToHex(block.TimeStamp)
	heightBytes := IntToHex(block.Height)
	blockBytes := bytes.Join([][]byte{
		heightBytes,
		timeStampBytes,
		block.PrevBlockHash,
		block.HashTransaction(),
	}, []byte{})
	hash := sha256.Sum256(blockBytes)
	block.Hash = hash[:]
}

//初始化创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(1, nil, txs)
}

//区块结构序列化
func (block *Block) Serialize() []byte {
	//gob
	var buffer bytes.Buffer
	//新建编码对象
	encoder := gob.NewEncoder(&buffer)
	//编码
	err := encoder.Encode(block)
	if nil != err {
		log.Panicf("serialize the block to []byte failed %v\n", err)
	}

	return buffer.Bytes()
}

//区块结构饭序列化
func UnSerializeBlock(blockBytes []byte) *Block {
	var block Block
	decode := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decode.Decode(&block)
	if err != nil {
		log.Panicf("deserialize the []bte to block failed %v\n", err)
	}
	return &block
}

//把区块中所有的交易序列化
func (block *Block) HashTransaction() []byte {
	var txHashes [][]byte
	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}
