// Package llvm 提供系统安装版 LLVM 的 Go 封装。
//
// 包布局：
//
//	llvm        核心词汇：Kind、Type[T]、Value[T]、常量、Context、错误与生命周期
//	llvm/ir     IR 构建：Module、Function、Block、Builder、指令
//	llvm/target 目标机器与代码生成
//	llvm/jit    ORC LLJIT 执行引擎
//	llvm/pass   优化管线
//
// 类型安全：值与类型以种类级泛型 Value[T]/Type[T] 表达，类别专属操作放在角色包装上
// （如 IntType.Bits()、Alloca.SetAlign()）。
//
// 错误处理：运行时可失败操作返回 error；程序员错误 panic(*Error)，可用 Catch 收敛。
package llvm
