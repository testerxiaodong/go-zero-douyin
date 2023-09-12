package utils

import "github.com/zeromicro/go-zero/core/threading"

var testGoImpl TestSafeGo

type TestSafeGo struct {
	ignoreGo bool
}

func NewTestGo() *TestSafeGo {
	return &TestSafeGo{}
}

// IgnoreGo 在单测里面要执行这个，用来忽略 go
func IgnoreGo() {
	testGoImpl.ignoreGo = true
}

// RecoverGo 不忽略 go
func RecoverGo() {
	testGoImpl.ignoreGo = false
}

func (t *TestSafeGo) RunSafe(f func()) {
	if testGoImpl.ignoreGo {
		return
	}

	// 正常的业务逻辑，还是正常 go 一个协程出去执行业务逻辑
	threading.GoSafe(f)
}
