package mlog

// 模块名
func FieldMod(modName string) Field {
	return String("module", modName)
}

// 服务地址
func FieldAddr(addr string) Field {
	return String("addr", addr)
}
