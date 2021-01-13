package BLC

import (
	"fmt"
	"os"
)

//查询余额
func (cli *CLI) getBalance(from string) {
	if !isDbExist() {
		fmt.Println("数据库不存在。")
		os.Exit(1)
	}

	//查找该地址UTXO
	bc := GetBlockChainObject()
	defer bc.DB.Close()
	amount := bc.getBalance(from)
	fmt.Printf("the balance of %s is %d", from, amount)
}
