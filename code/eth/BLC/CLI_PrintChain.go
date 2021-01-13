package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) printChain() {
	if !isDbExist() {
		fmt.Println("数据库不存在。")
		os.Exit(1)
	}

	blockchain := GetBlockChainObject()
	blockchain.PrintChain()
	defer blockchain.DB.Close()
}
