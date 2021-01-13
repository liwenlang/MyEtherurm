package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) sendTx(from, to, amount []string) {
	//.\bc.exe  send -from "[\"lwl\",\"wenlang\"]"  -to "[\"wenlang\",\"lwl\"]" -amount "[\"5\",\"2\"]"
	//.\bc.exe  send -from "[\"lwl\"]" -to "[\"wenlang\"]" -amount "[\"1\"]"
	//.\bc.exe  send -from lwl" -to wenlang -amount 1
	if !isDbExist() {
		fmt.Println("数据库不存在。")
		os.Exit(1)
	}

	blockchain := GetBlockChainObject()
	defer blockchain.DB.Close()

	if len(from) != len(to) || len(from) != len(amount) {
		fmt.Println("交易参数输入有误，请检查一致性")
		os.Exit(1)
	}

	blockchain.MineNewBlock(from, to, amount)
}
