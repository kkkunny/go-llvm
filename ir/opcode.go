package ir

import (
	"strconv"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/internal/binding"
)

// Op 指令操作码；转发 binding.LLVMOpcode（仿 llvm.Linkage 先例）
type Op binding.LLVMOpcode

// Op 取值对应 LLVM 指令操作码（binding.LLVMOpcode）；每个取值对应 LangRef 中的一类指令，
// OpUserOp1/OpUserOp2 除外（保留给 pass 内部使用）。未知操作码不会 panic：String 输出 op<N>，
// 角色转换（AsLoad 等）返回 false，随后以不透明的 Value[DynT] 处理。
const (
	OpRet            = Op(binding.LLVMRet)            // 从函数返回
	OpBr             = Op(binding.LLVMBr)             // 无条件或条件分支
	OpSwitch         = Op(binding.LLVMSwitch)         // 多路分支：匹配 case 则跳转，否则跳 default
	OpIndirectBr     = Op(binding.LLVMIndirectBr)     // 间接跳转：跳往操作数列表中的某个基本块（计算 goto）
	OpInvoke         = Op(binding.LLVMInvoke)         // 可展开（unwind）的调用，正常返回与异常路径分开（终结指令）
	OpUnreachable    = Op(binding.LLVMUnreachable)    // 不可达点；执行到此处为未定义行为
	OpFNeg           = Op(binding.LLVMFNeg)           // 浮点取负（一元运算）
	OpAdd            = Op(binding.LLVMAdd)            // 整数加法（可带 nsw/nuw）
	OpFAdd           = Op(binding.LLVMFAdd)           // 浮点加法
	OpSub            = Op(binding.LLVMSub)            // 整数减法（可带 nsw/nuw）
	OpFSub           = Op(binding.LLVMFSub)           // 浮点减法
	OpMul            = Op(binding.LLVMMul)            // 整数乘法（可带 nsw/nuw）
	OpFMul           = Op(binding.LLVMFMul)           // 浮点乘法
	OpUDiv           = Op(binding.LLVMUDiv)           // 无符号整数除法（除零为未定义行为，可带 exact）
	OpSDiv           = Op(binding.LLVMSDiv)           // 有符号整数除法（除零或 INT_MIN/-1 为未定义行为，可带 exact）
	OpFDiv           = Op(binding.LLVMFDiv)           // 浮点除法
	OpURem           = Op(binding.LLVMURem)           // 无符号整数取余（除零为未定义行为）
	OpSRem           = Op(binding.LLVMSRem)           // 有符号整数取余（结果符号与被除数一致；除零或 INT_MIN % -1 为未定义行为）
	OpFRem           = Op(binding.LLVMFRem)           // 浮点取余（fmod 语义，结果符号与被除数一致）
	OpShl            = Op(binding.LLVMShl)            // 左移：低位补 0；移位量不小于位宽时得 poison
	OpLShr           = Op(binding.LLVMLShr)           // 逻辑右移：高位补 0；移位量不小于位宽时得 poison
	OpAShr           = Op(binding.LLVMAShr)           // 算术右移：高位补符号位；移位量不小于位宽时得 poison
	OpAnd            = Op(binding.LLVMAnd)            // 按位与
	OpOr             = Op(binding.LLVMOr)             // 按位或
	OpXor            = Op(binding.LLVMXor)            // 按位异或
	OpAlloca         = Op(binding.LLVMAlloca)         // 在当前函数的栈帧上分配内存，返回指向它的指针
	OpLoad           = Op(binding.LLVMLoad)           // 从内存读取值
	OpStore          = Op(binding.LLVMStore)          // 向内存写入值
	OpGetElementPtr  = Op(binding.LLVMGetElementPtr)  // 计算聚合体元素或指针目标元素的地址（只算地址，不访存）
	OpTrunc          = Op(binding.LLVMTrunc)          // 整数截断到更窄类型
	OpZExt           = Op(binding.LLVMZExt)           // 整数零扩展（高位补 0）
	OpSExt           = Op(binding.LLVMSExt)           // 整数符号扩展（高位复制符号位）
	OpFPToUI         = Op(binding.LLVMFPToUI)         // 浮点转无符号整数（舍入到零；超出目标范围得到 poison）
	OpFPToSI         = Op(binding.LLVMFPToSI)         // 浮点转有符号整数（舍入到零；超出目标范围得到 poison）
	OpUIToFP         = Op(binding.LLVMUIToFP)         // 无符号整数转浮点
	OpSIToFP         = Op(binding.LLVMSIToFP)         // 有符号整数转浮点
	OpFPTrunc        = Op(binding.LLVMFPTrunc)        // 浮点截断到更窄类型（可能丢失精度）
	OpFPExt          = Op(binding.LLVMFPExt)          // 浮点扩展到更宽类型
	OpPtrToInt       = Op(binding.LLVMPtrToInt)       // 指针转整数（捕获地址与 provenance）
	OpIntToPtr       = Op(binding.LLVMIntToPtr)       // 整数转指针（能否安全解引用取决于该整数是否源自 ptrtoint）
	OpBitCast        = Op(binding.LLVMBitCast)        // 位级重解释转换，要求源与目标位宽相同
	OpAddrSpaceCast  = Op(binding.LLVMAddrSpaceCast)  // 指针地址空间转换
	OpICmp           = Op(binding.LLVMICmp)           // 整数或指针比较，返回 i1
	OpFCmp           = Op(binding.LLVMFCmp)           // 浮点比较（谓词分有序/无序），返回 i1
	OpPHI            = Op(binding.LLVMPHI)            // φ 节点：按前驱基本块选择到达该块时生效的值
	OpCall           = Op(binding.LLVMCall)           // 函数调用（非终结指令）
	OpSelect         = Op(binding.LLVMSelect)         // 条件选择：条件为真取第一个值，否则取第二个（向量条件逐元素选择）
	OpUserOp1        = Op(binding.LLVMUserOp1)        // 保留给 pass 内部使用的自定义操作 1（正常 IR 中不应出现）
	OpUserOp2        = Op(binding.LLVMUserOp2)        // 保留给 pass 内部使用的自定义操作 2（正常 IR 中不应出现）
	OpVAArg          = Op(binding.LLVMVAArg)          // 从可变参数列表（va_list）中取下一个实参
	OpExtractElement = Op(binding.LLVMExtractElement) // 按下标从向量中取出元素
	OpInsertElement  = Op(binding.LLVMInsertElement)  // 按下标向向量中写入元素，返回新向量
	OpShuffleVector  = Op(binding.LLVMShuffleVector)  // 按掩码重排两个向量的元素，构造新向量
	OpExtractValue   = Op(binding.LLVMExtractValue)   // 按下标从聚合体（结构体/数组）中取出成员
	OpInsertValue    = Op(binding.LLVMInsertValue)    // 按下标向聚合体中写入成员，返回新聚合体
	OpFence          = Op(binding.LLVMFence)          // 内存栅栏：与原子操作配合建立跨线程的 happens-before 关系
	OpAtomicCmpXchg  = Op(binding.LLVMAtomicCmpXchg)  // 原子比较交换：返回旧值与是否交换成功
	OpAtomicRMW      = Op(binding.LLVMAtomicRMW)      // 原子读改写：返回修改前的旧值
	OpResume         = Op(binding.LLVMResume)         // 恢复异常的栈展开（终结指令）
	OpLandingPad     = Op(binding.LLVMLandingPad)     // 异常着陆点：unwind 进入本块时的首指令（旧式 EH）
	OpCleanupRet     = Op(binding.LLVMCleanupRet)     // 退出 cleanup pad（WinEH 终结指令）
	OpCatchRet       = Op(binding.LLVMCatchRet)       // 离开 catch pad 并跳转到普通基本块（WinEH 终结指令）
	OpCatchPad       = Op(binding.LLVMCatchPad)       // WinEH catch 结构的入口 pad
	OpCleanupPad     = Op(binding.LLVMCleanupPad)     // WinEH cleanup 结构的入口 pad
	OpCatchSwitch    = Op(binding.LLVMCatchSwitch)    // WinEH 终结指令：把异常分派给一组 handler pad
	OpCallBr         = Op(binding.LLVMCallBr)         // 可分支的调用：用于带间接跳转标签的内联汇编（终结指令）
	OpFreeze         = Op(binding.LLVMFreeze)         // 冻结 poison/undef：返回一个此后不再变化的值，并阻止 poison 传播
	OpPtrToAddr      = Op(binding.LLVMPtrToAddr)      // 指针转整数地址（ptrtoaddr）：只取索引位、不捕获 provenance
)

// opNames 操作码短名（String 与调试输出用；未知码输出 op<N>）
var opNames = map[Op]string{
	OpRet: "ret", OpBr: "br", OpSwitch: "switch", OpIndirectBr: "indirectbr",
	OpInvoke: "invoke", OpUnreachable: "unreachable", OpFNeg: "fneg",
	OpAdd: "add", OpFAdd: "fadd", OpSub: "sub", OpFSub: "fsub",
	OpMul: "mul", OpFMul: "fmul", OpUDiv: "udiv", OpSDiv: "sdiv", OpFDiv: "fdiv",
	OpURem: "urem", OpSRem: "srem", OpFRem: "frem", OpShl: "shl",
	OpLShr: "lshr", OpAShr: "ashr", OpAnd: "and", OpOr: "or", OpXor: "xor",
	OpAlloca: "alloca", OpLoad: "load", OpStore: "store", OpGetElementPtr: "getelementptr",
	OpTrunc: "trunc", OpZExt: "zext", OpSExt: "sext", OpFPToUI: "fptoui",
	OpFPToSI: "fptosi", OpUIToFP: "uitofp", OpSIToFP: "sitofp", OpFPTrunc: "fptrunc",
	OpFPExt: "fpext", OpPtrToInt: "ptrtoint", OpIntToPtr: "inttoptr",
	OpBitCast: "bitcast", OpAddrSpaceCast: "addrspacecast", OpICmp: "icmp",
	OpFCmp: "fcmp", OpPHI: "phi", OpCall: "call", OpSelect: "select",
	OpUserOp1: "userop1", OpUserOp2: "userop2", OpVAArg: "va_arg",
	OpExtractElement: "extractelement", OpInsertElement: "insertelement",
	OpShuffleVector: "shufflevector", OpExtractValue: "extractvalue",
	OpInsertValue: "insertvalue", OpFence: "fence", OpAtomicCmpXchg: "cmpxchg",
	OpAtomicRMW: "atomicrmw", OpResume: "resume", OpLandingPad: "landingpad",
	OpCleanupRet: "cleanupret", OpCatchRet: "catchret", OpCatchPad: "catchpad",
	OpCleanupPad: "cleanuppad", OpCatchSwitch: "catchswitch", OpCallBr: "callbr",
	OpFreeze: "freeze", OpPtrToAddr: "ptrtoaddr",
}

// String 操作码名（小写，与 LLVM IR 文本一致）
func (o Op) String() string {
	if s, ok := opNames[o]; ok {
		return s
	}
	return "op" + strconv.Itoa(int(o))
}

// IsTerminator 操作码是否为终结指令：必然位于基本块末尾并决定控制流走向。
// 与底层 LLVMIsATerminatorInst（见 [SuccessorCount] 的拒绝路径）一致；
// 判断"块是否已被终结"用 [Block.IsTerminating]，不要只查最后一条指令的操作码。
func (o Op) IsTerminator() bool {
	switch o {
	case OpRet, OpBr, OpSwitch, OpIndirectBr, OpInvoke, OpUnreachable,
		OpResume, OpCleanupRet, OpCatchRet, OpCatchSwitch, OpCallBr:
		return true
	default:
		return false
	}
}

// OpOf 指令的操作码；非指令值（常量/参数/全局等）返回 false
func OpOf(inst llvm.AnyValue) (Op, bool) {
	const op = "ir.OpOf"
	if inst == nil || !inst.Alive() {
		llvm.Panicf(llvm.ErrInvalidArg, op, "nil or dead value")
	}
	ref := inst.Ref()
	if binding.LLVMGetValueKind(ref) != binding.LLVMInstructionValueKind {
		return 0, false
	}
	return Op(binding.LLVMGetInstructionOpcode(ref)), true
}
