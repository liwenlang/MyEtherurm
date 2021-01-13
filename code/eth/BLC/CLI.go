package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

//对blockchain的命令行进行管理

//client对象
type CLI struct {
}

//用法展示

func PrintUsage() {
	fmt.Println("Usage:")

	//初始化区块链
	fmt.Println("\tcreateblockchain -address ADDRESS-- 创建区块链")
	//添加区块
	//fmt.Println("\taddblock -data DATA --添加区块")
	//打印完整区块信息
	fmt.Println("\tprintChain -- 输出区块链信息")
	//添加转账
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT --发起转账")
	fmt.Println("\t\t-from FROM --转账源地址")
	fmt.Println("\t\t-to TO --转账目标地址")
	fmt.Println("\t\t-amount AMOUNT --转账金额")

	//查询余额
	fmt.Println("\tgetbalance -address FROM --查询指定地址余额")

}

func (cli *CLI) addBlock(txs []*Transaction) {
	if !isDbExist() {
		fmt.Println("数据库不存在。")
		os.Exit(1)
	}

	blockchain := GetBlockChainObject()
	blockchain.AddBlock(txs)
	defer blockchain.DB.Close()
}

//命令行运行函数
func (cli *CLI) Run() {
	//检查参数
	IsVaildArgs()

	//新建相关命令
	createBLCWithGenesisBlockCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	sendTxCmd := flag.NewFlagSet("send", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	//数据参数处理
	flagAddBlockArg := addBlockCmd.String("data", "send 100 btc to player", "增加区块链数据")
	//指定矿工地址
	flagCreateBlockChainArg := createBLCWithGenesisBlockCmd.String("address", "liwenlang", "指定接受系统奖励的矿工地址")
	var strDef = "[\"lwl\"]"
	flagSendTxFromArg := sendTxCmd.String("from", strDef, "转账源地址")
	strDef = "[\"wenlang\"]"
	flagSendTxToArg := sendTxCmd.String("to", strDef, "转账目标地址")
	strDef = "[\"88\"]"
	flagSendTxAmountArg := sendTxCmd.String("amount", strDef, "转账金额")

	flagGetBalanceArg := getBalanceCmd.String("address", "liwenlang", "查询余额")
	//判断命令
	switch os.Args[1] {
	case "getbalance":
		if err := getBalanceCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse getBalanceCmd failed %v", err)
		}
	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse addBlockCmd failed %v", err)
		}
	case "printchain":
		if err := printChainCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse printChainCmd failed %v", err)
		}
	case "createblockchain":
		if err := createBLCWithGenesisBlockCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse createBLCWithGenesisBlockCmd failed %v", err)
		}
	case "send":
		if err := sendTxCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse send Transaction failed %v", err)
		}
	default:
		PrintUsage()
		os.Exit(1)
	}

	//输出区块链信息
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	//创建区块链
	if createBLCWithGenesisBlockCmd.Parsed() {
		if *flagCreateBlockChainArg == "" {
			fmt.Println("地址不能为空")
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockChain(*flagCreateBlockChainArg)
	}
	//添加区块命令
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.addBlock([]*Transaction{})
	}

	if sendTxCmd.Parsed() {
		if *flagSendTxFromArg == "" {
			fmt.Println("源地址不能为空")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendTxToArg == "" {
			fmt.Println("目标地址不能为空")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendTxAmountArg == "" {
			fmt.Println("金额不能为空")
			PrintUsage()
			os.Exit(1)
		}
		// fmt.Printf("转账交易：From %s To %s %s Btc", *flagSendTxFromArg,
		// 	*flagSendTxToArg, *flagSendTxAmountArg)

		fmt.Printf("\tFROM:[%v]\n", JsonToSlice(*flagSendTxFromArg))
		fmt.Printf("\tTO:[%v]\n", JsonToSlice(*flagSendTxToArg))
		fmt.Printf("\tAMOUNT:[%v]\n", JsonToSlice(*flagSendTxAmountArg))

		cli.sendTx(JsonToSlice(*flagSendTxFromArg), JsonToSlice(*flagSendTxToArg), JsonToSlice(*flagSendTxAmountArg))
	}

	//创建区块链
	if getBalanceCmd.Parsed() {
		if *flagGetBalanceArg == "" {
			fmt.Println("地址不能为空")
			PrintUsage()
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceArg)
	}
}
