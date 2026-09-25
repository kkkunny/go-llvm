package llvm

import "github.com/kkkunny/go-llvm/internal/binding"

// Linkage 链接类型
type Linkage binding.LLVMLinkage

// Linkage 取值对应 LLVM 链接类型（binding.LLVMLinkage），决定符号的可见范围
// 以及链接器合并同名定义时的取舍规则。
const (
	LinkageExternal            = Linkage(binding.LLVMExternalLinkage)            // 外部链接：符号对外可见，可由其他模块引用（默认）
	LinkageAvailableExternally = Linkage(binding.LLVMAvailableExternallyLinkage) // 外部可用：定义仅供编译期优化，不写入目标文件，可随时丢弃
	LinkageLinkOnceAny         = Linkage(binding.LLVMLinkOnceAnyLinkage)         // 链接一次（任意）：链接时保留一份副本，可被更强定义覆盖，未被引用则丢弃
	LinkageLinkOnceODR         = Linkage(binding.LLVMLinkOnceODRLinkage)         // 链接一次（ODR）：链接时保留一份副本，仅允许被等价定义替换
	LinkageLinkOnceODRAutoHide = Linkage(binding.LLVMLinkOnceODRAutoHideLinkage) // 链接一次（ODR、自动隐藏）：已废弃
	LinkageWeakAny             = Linkage(binding.LLVMWeakAnyLinkage)             // 弱链接（任意）：链接时保留一份，可被非弱定义覆盖，未被引用则丢弃
	LinkageWeakODR             = Linkage(binding.LLVMWeakODRLinkage)             // 弱链接（ODR）：仅允许被等价定义替换
	LinkageAppending           = Linkage(binding.LLVMAppendingLinkage)           // 追加链接：仅用于全局数组，链接时按顺序合并所有定义
	LinkageInternal            = Linkage(binding.LLVMInternalLinkage)            // 内部链接：仅本模块可见（类似 static）
	LinkagePrivate             = Linkage(binding.LLVMPrivateLinkage)             // 私有链接：仅本模块可见且不进入符号表
	LinkageExternalWeak        = Linkage(binding.LLVMExternalWeakLinkage)        // 外部弱引用：符号缺失时解析为 null，而非链接错误
	LinkageCommon              = Linkage(binding.LLVMCommonLinkage)              // 公共（暂定）定义：链接时合并多个同名定义，常用于未初始化全局变量
)

// Visibility 符号可见性
type Visibility binding.LLVMVisibility

// Visibility 取值对应 LLVM 符号可见性（binding.LLVMVisibility），
// 控制符号是否导出到动态符号表以及能否被其他模块覆盖。
const (
	VisibilityDefault   = Visibility(binding.LLVMDefaultVisibility)   // 默认：符号对外可见，可被其他模块引用
	VisibilityHidden    = Visibility(binding.LLVMHiddenVisibility)    // 隐藏：符号不导出，仅在本模块内可见
	VisibilityProtected = Visibility(binding.LLVMProtectedVisibility) // 受保护：符号对外可见但不可被覆盖，引用绑定到本模块定义
)

// DLLStorageClass DLL 存储类
type DLLStorageClass binding.LLVMDLLStorageClass

// DLLStorageClass 取值对应 LLVM DLL 存储类（binding.LLVMDLLStorageClass），
// 用于 Windows 平台的 DLL 导入/导出标记。
const (
	DLLStorageDefault = DLLStorageClass(binding.LLVMDefaultStorageClass)   // 默认：非 DLL 导入/导出
	DLLStorageImport  = DLLStorageClass(binding.LLVMDLLImportStorageClass) // 从 DLL 导入
	DLLStorageExport  = DLLStorageClass(binding.LLVMDLLExportStorageClass) // 导出到 DLL
)

// UnnamedAddr 匿名地址语义
type UnnamedAddr binding.LLVMUnnamedAddr

// UnnamedAddr 取值对应 LLVM 匿名地址语义（binding.LLVMUnnamedAddr），
// 表示全局值的地址是否可被忽略，以便与等价符号合并。
const (
	UnnamedAddrNone   = UnnamedAddr(binding.LLVMNoUnnamedAddr)     // 地址有意义：不可与其他等价符号合并
	UnnamedAddrLocal  = UnnamedAddr(binding.LLVMLocalUnnamedAddr)  // 地址在模块内不重要：可与本模块内的等价符号合并
	UnnamedAddrGlobal = UnnamedAddr(binding.LLVMGlobalUnnamedAddr) // 地址全局不重要：可与任意等价符号合并
)

// ThreadLocalMode 线程局部模式
type ThreadLocalMode binding.LLVMThreadLocalMode

// ThreadLocalMode 取值对应 LLVM 线程局部存储模型（binding.LLVMThreadLocalMode），
// 对应 IR 的 thread_local 属性。
const (
	ThreadLocalNone         = ThreadLocalMode(binding.LLVMNotThreadLocal)         // 非线程局部变量（默认）
	ThreadLocalGeneral      = ThreadLocalMode(binding.LLVMGeneralDynamicTLSModel) // general-dynamic：通用动态 TLS，运行时经 __tls_get_addr 解析
	ThreadLocalLocalDynamic = ThreadLocalMode(binding.LLVMLocalDynamicTLSModel)   // local-dynamic：同一动态模块内多个 TLS 变量合并解析
	ThreadLocalInitialExec  = ThreadLocalMode(binding.LLVMInitialExecTLSModel)    // initial-exec：模块加载时分配，经固定偏移访问，不适用于运行时动态加载
	ThreadLocalLocalExec    = ThreadLocalMode(binding.LLVMLocalExecTLSModel)      // local-exec：可执行文件内直接静态偏移访问，效率最高
)
