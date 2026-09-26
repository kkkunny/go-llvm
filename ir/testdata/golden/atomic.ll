; ModuleID = 'atomic'
source_filename = "atomic"

define void @atomics(ptr %0) {
entry:
  fence seq_cst
  %old = atomicrmw add ptr %0, i32 1 monotonic, align 4
  store atomic i32 1, ptr %0 release, align 4
  %cur = load atomic i32, ptr %0 acquire, align 4
  %cx = cmpxchg ptr %0, i32 %cur, i32 2 syncscope("singlethread") acquire monotonic, align 4
  ret void
}
