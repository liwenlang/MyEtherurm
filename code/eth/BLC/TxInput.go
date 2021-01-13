package BLC

//交易输入管理

//输入结构
type TxInput struct {
	TxHash    []byte //交易哈希
	Vout      int    //引用上一笔交易的输出索引（coinbase -1）
	ScriptSig string //用户名
}

//验证引用的地址是否匹配
func (txInput *TxInput) CheckPubKey(address string) bool {
	return txInput.ScriptSig == address
}
