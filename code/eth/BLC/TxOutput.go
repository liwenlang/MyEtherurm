package BLC

type TxOutput struct {
	Value        int    //金额
	ScriptPubKey string //UTXO拥有者
}

//验证当前UTXO是否属于指定地址
func (txOutput *TxOutput) CheckPubKey(address string) bool {
	return txOutput.ScriptPubKey == address
}
