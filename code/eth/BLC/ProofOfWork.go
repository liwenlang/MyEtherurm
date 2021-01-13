package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//共识算法管理文件

//实现POW实例以及相关功能

//目标难度值
const targetBit = 16

//工作量证明的结构
type ProofOfWork struct {
	//需要共识验证的区块
	Block *Block
	//目标难度的哈希
	target *big.Int
}

//创建POW对象
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)

	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{Block: block, target: target}
}

//执行pow,比较hash，返回哈希值及碰撞次数
func (pow *ProofOfWork) Run() ([]byte, int64) {
	var nonce int64
	var hashInt big.Int
	var hash [32]byte

	for {
		//生成准备数据
		dataBytes := pow.prepareData(nonce)
		hash = sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		//检验是否符合条件
		if pow.target.Cmp(&hashInt) == 1 {
			break
		}

		nonce++
	}

	fmt.Printf("碰撞次数:%d\n", nonce)
	return hash[:], nonce
}

//生成准备数据
func (pow *ProofOfWork) prepareData(nonce int64) []byte {

	timeStampBytes := IntToHex(pow.Block.TimeStamp)
	heightBytes := IntToHex(pow.Block.Height)
	nonceBytes := IntToHex(nonce)
	blockBytes := bytes.Join([][]byte{
		heightBytes,
		timeStampBytes,
		pow.Block.PrevBlockHash,
		pow.Block.HashTransaction(),
		nonceBytes,
	}, []byte{})
	hash := sha256.Sum256(blockBytes)
	pow.Block.Hash = hash[:]

	return hash[:]
}
