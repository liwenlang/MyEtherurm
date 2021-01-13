package BLC

//UTXO结构管理
type UTXO struct {
	TxHash []byte    //UTXO对应的交易哈希
	Index  int       //所在输出列表索引
	Output *TxOutput //输出
}
