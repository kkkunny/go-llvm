package llvm

import "sync/atomic"

// Lifetime 生命周期令牌；ir.Module 持有并传给每个值，实现 use-after-free 检测
type Lifetime struct{ dead atomic.Bool }

// NewLifetime 创建存活的生命周期令牌
func NewLifetime() *Lifetime { return &Lifetime{} }

// Kill 标记失效（Close 时调用）
func (l *Lifetime) Kill() { l.dead.Store(true) }

// Alive 是否存活
func (l *Lifetime) Alive() bool { return !l.dead.Load() }
