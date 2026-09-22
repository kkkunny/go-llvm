package ir

import (
	"github.com/kkkunny/go-llvm"
)

// NewFunc 按 Go 函数签名声明/定义函数。
//
// 签名映射规则见 llvm.FnSignatureOf；不支持的类型返回 llvm.ErrUnsupported。
func (m *Module) NewFunc[F any](name string) (Func[F], error) {
	sig, goTy, err := llvm.FnSignatureOf[F](m.ctx)
	if err != nil {
		return Func[F]{}, err
	}
	return Func[F]{fn: m.NewFunction(name, sig), goType: goTy}, nil
}
