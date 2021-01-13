package BLC

import (
	"log"

	"github.com/boltdb/bolt"
)

//区块链迭代器管理文件

//基本结构
type BlockChainIterator struct {
	DB      *bolt.DB //迭代目标
	CurHash []byte   //当前哈希
}

//创建迭代对象
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.DB, bc.Tip}
}

//实现迭代函数Next
func (bcIter *BlockChainIterator) Next() *Block {
	var block *Block

	err := bcIter.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			curBlockHash := b.Get(bcIter.CurHash)
			block = UnSerializeBlock(curBlockHash)

			bcIter.CurHash = block.PrevBlockHash
		}

		return nil
	})

	if nil != err {
		log.Panicf("iterator the db failed %v\n", err)
	}

	return block
}
