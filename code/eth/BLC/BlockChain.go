package BLC

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
)

//数据库名称
const dbName = "block.db"

//表名称
const blockTableName = "blocks"

type BlockChain struct {
	//Blocks []*Block
	DB  *bolt.DB //数据库对象
	Tip []byte   //最新区块的哈希值
}

//判断数据库文件是否存在
func isDbExist() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}

	return true
}

//初始化区块链
func CreateBlockChainWithGenesisBlock(address string) *BlockChain {

	if isDbExist() {
		fmt.Println("数据库已存在...")
		os.Exit(1)
	}
	//1.打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panicf("create db [%s] failed %v\n", dbName, err)

	}

	var genesisBlock *Block

	//从创建桶，保存至数据库中
	db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b == nil {
			b, err = tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				log.Panicf("create bucket [%s] failed %v\n", blockTableName, err)
			}

		}

		//生成一个coinbase交易
		txCoinbase := NewCoinbaseTransaction(address)

		//生成创世区块
		genesisBlock = CreateGenesisBlock([]*Transaction{txCoinbase})

		//fmt.Printf("genesisBlock:%v\n", genesisBlock)

		//存储
		err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
		if err != nil {
			log.Panicf("insert the genesis block failed %v\n", err)
		}

		//存储最新区块的哈希
		err = b.Put([]byte("lasestHash"), genesisBlock.Hash)
		if err != nil {
			log.Panicf("save the hash of genesis block failed %v\n", err)
		}

		return nil
	})

	return &BlockChain{db, genesisBlock.Hash}
}

//添加区块到区块链中
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	// newBlock := NewBlock(bc.Blocks[len(bc.Blocks)-1].Height+1,
	// 	bc.Blocks[len(bc.Blocks)-1].Hash, data)
	// bc.Blocks = append(bc.Blocks, newBlock)

	bc.DB.Update(func(tx *bolt.Tx) error {
		//1.打开数据库桶
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			//2.获取最新区块哈希
			blockBytes := b.Get(bc.Tip)
			//3.反序化
			lasestBlock := UnSerializeBlock(blockBytes)
			//4.新建区块
			newBlock := NewBlock(lasestBlock.Height+1, lasestBlock.Hash, txs)
			//5.存入数据库
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panicf("insert the new block to db failed %v\n", err)
			}
			//更新数据最新区块哈希
			err = b.Put([]byte("lasestHash"), newBlock.Hash)
			if err != nil {
				log.Panicf("update the lasest hast to db failed %v\n", err)
			}

			bc.Tip = newBlock.Hash
		}
		return nil
	})
}

func (bc *BlockChain) PrintChain() {
	//读取数据库
	bcIter := bc.Iterator() //获取迭代对象

	for {
		fmt.Println("-------------------------------")

		curBlock := bcIter.Next()
		//输出区块链
		fmt.Printf("\tTimeStamp:%v\n", curBlock.TimeStamp)
		fmt.Printf("\tHash:%x\n", curBlock.Hash)
		fmt.Printf("\tPrevBlockHash:%x\n", curBlock.PrevBlockHash)
		fmt.Printf("\tHeight:%d\n", curBlock.Height)
		fmt.Printf("\tTxs:%v\n", curBlock.Txs)
		for _, tx := range curBlock.Txs {
			fmt.Printf("\t\ttx-hash: %x\n", tx.TxHash)
			fmt.Printf("\t\t输入...\n")
			for _, vin := range tx.Vins {
				fmt.Printf("\t\t\tvin-txHash: %x\n", vin.TxHash)
				fmt.Printf("\t\t\tvin-vout: %d\n", vin.Vout)
				fmt.Printf("\t\t\tvin-sriptsig: %s\n", vin.ScriptSig)
			}
			fmt.Printf("\t\t输出...\n")
			for _, vout := range tx.Vouts {
				fmt.Printf("\t\t\tvout-txValue: %d\n", vout.Value)
				fmt.Printf("\t\t\tvout-ScriptPubKey: %s\n", vout.ScriptPubKey)
			}
		}
		fmt.Printf("\tNonce:%d\n", curBlock.Nonce)

		//退出条件
		var prevHashInt big.Int
		prevHashInt.SetBytes(curBlock.PrevBlockHash)
		if big.NewInt(0).Cmp(&prevHashInt) == 0 {
			break
		}
	}
}

//获取blockchain对象
func GetBlockChainObject() *BlockChain {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panicf("open bolt db failed %v \n", err)
	}

	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			tip = b.Get([]byte("lasestHash"))
		}

		return nil
	})

	if err != nil {
		log.Panicf("get lasestHash failed %v\n", err)
	}

	return &BlockChain{db, tip}

}

//实现挖矿
//通过接受交易 生成区块
func (bc *BlockChain) MineNewBlock(from, to, amount []string) {
	//生成交易
	var txs []*Transaction

	//多笔交易
	for index, address := range from {
		//string转int
		value, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTansaction(address, to[index], value, bc, txs)
		txs = append(txs, tx)
	}

	var block *Block
	//从数据库中获取最新区块
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			lasestHash := b.Get([]byte("lasestHash"))
			//fmt.Printf("lasestHash: %x\n", lasestHash)
			//反序列化
			blockBytes := b.Get(lasestHash)

			block = UnSerializeBlock(blockBytes)
		}

		return nil
	})

	//通过数据库中最新的区块生成新区块
	block = NewBlock(block.Height+1, block.Hash, txs)

	//持久化新生成的区块至数据库中
	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			err := b.Put(block.Hash, block.Serialize())
			if err != nil {
				log.Panicf("update the new block to db failed %v\n", err)
			}

			//更新最新区块
			err = b.Put([]byte("lasestHash"), block.Hash)
			if err != nil {
				log.Panicf("save lasestHash failed %v", err)
			}

			bc.Tip = block.Hash
		}

		return nil
	})
}

//获取所有已花费输出
func (bc *BlockChain) SpentOutputs(address string) map[string][]int {
	spentTxOutputs := make(map[string][]int)
	bcIter := bc.Iterator()
	for {
		block := bcIter.Next()
		for _, tx := range block.Txs {
			if !tx.IsCoinbaseTransaction() {
				for _, in := range tx.Vins {
					if in.CheckPubKey(address) {
						key := hex.EncodeToString(in.TxHash)
						spentTxOutputs[key] = append(spentTxOutputs[key], in.Vout)
					}
				}
			}
		}

		//退出条件
		var prevHashInt big.Int
		prevHashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&prevHashInt) == 0 {
			break
		}
	}
	return spentTxOutputs
}

//查找指定地址的UTXO
/*
	遍历查找每一个交易每一个输出
	判断满足以下条件：
	1.属于传入的地址
	2.是否未被花费
		2.1遍历所有区块，记录所有已花费TUXO
		2.2检查每一项vout是否在vin被引用
*/
func (bc *BlockChain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	//获取迭代器
	bcIter := bc.Iterator()
	spentTxOutputs := bc.SpentOutputs(address)
	var unUTXOs []*UTXO
	//缓存迭代
	//查找已花费输出
	for _, tx := range txs {
		if !tx.IsCoinbaseTransaction() {
			for _, in := range tx.Vins {
				if in.CheckPubKey(address) {
					key := hex.EncodeToString(in.TxHash)
					spentTxOutputs[key] = append(spentTxOutputs[key], in.Vout)
				}
			}
		}
	}

	//同一个区块多笔交易，交易哈希一样，此处处理方式和数据库不一致
	for _, tx := range txs {
		curTxHash := hex.EncodeToString(tx.TxHash)
	workCache:
		for index, vout := range tx.Vouts {
			if vout.CheckPubKey(address) { //////////////1//////
				if len(spentTxOutputs) != 0 {
					//标记
					var isUtxoTx bool //判断交易是被其他交易所引用
					for txHash, indexArr := range spentTxOutputs {
						//txHash：交易哈希
						//indexArr:vout索引列表
						if txHash == curTxHash { ////////////2////////////
							isUtxoTx = true
							var isSpentUtxo bool
							for i := range indexArr {
								if index == i { //////////////3///////////////
									isSpentUtxo = true
									continue workCache
								}
							}
							if !isSpentUtxo {
								utxo := &UTXO{tx.TxHash, index, vout}
								unUTXOs = append(unUTXOs, utxo)
							}
						}

					}
					if !isUtxoTx {
						utxo := &UTXO{tx.TxHash, index, vout}
						unUTXOs = append(unUTXOs, utxo)
					}
				} else {
					utxo := &UTXO{tx.TxHash, index, vout}
					unUTXOs = append(unUTXOs, utxo)
				}
			}
		}
	}

	//数据库迭代
	for {
		block := bcIter.Next()
		//遍历每个区块交易
		for _, tx := range block.Txs {
			curTxHash := hex.EncodeToString(tx.TxHash)
		work:
			for index, vout := range tx.Vouts {
				if vout.CheckPubKey(address) {
					if len(spentTxOutputs) != 0 {
						//标记
						var isSpentOutput bool
						for txHash, indexArr := range spentTxOutputs {
							//txHash：交易哈希
							//indexArr:vout索引列表
							for i := range indexArr {
								if txHash == curTxHash && index == i {
									isSpentOutput = true
									continue work
								}
							}
						}
						if !isSpentOutput {
							utxo := &UTXO{tx.TxHash, index, vout}
							unUTXOs = append(unUTXOs, utxo)
						}
					} else {
						utxo := &UTXO{tx.TxHash, index, vout}
						unUTXOs = append(unUTXOs, utxo)
					}
				}
			}
		}

		//退出条件
		var prevHashInt big.Int
		prevHashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&prevHashInt) == 0 {
			break
		}
	}

	return unUTXOs
}

//查询余额
func (bc *BlockChain) getBalance(from string) int {
	var amount int

	utxos := bc.UnUTXOs(from, []*Transaction{})
	for _, txOutput := range utxos {
		amount += txOutput.Output.Value
	}
	return amount
}

//查找指定地址的可用UTXO,超过amount中断，并更新UTXO数量
//txs []*Transaction 缓存交易列表
func (bc *BlockChain) FindSpendableUTXO(from string,
	amount int, txs []*Transaction) (int, map[string][]int) {
	spendableUTXOs := make(map[string][]int)
	var value int
	utxos := bc.UnUTXOs(from, txs)

	//遍历UTXO
	for _, utxo := range utxos {
		value += utxo.Output.Value
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXOs[hash] = append(spendableUTXOs[hash], utxo.Index)

		if value >= amount {
			break
		}
	}

	if value < amount {
		fmt.Printf("地址[%s] 余额不足，当前余额 [%d], 转账金额 [%d]\n", from, value, amount)
		os.Exit(1)
	}

	return value, spendableUTXOs
}
