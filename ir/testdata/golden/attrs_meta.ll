; ModuleID = 'attrs_meta'
source_filename = "attrs_meta"

$mycomdat = comdat nodeduplicate

@llvm.global_ctors = appending global [1 x { i32, ptr, ptr }] [{ i32, ptr, ptr } { i32 65535, ptr @ctor, ptr null }]
@glob = global i32 0, comdat($mycomdat)

; Function Attrs: noinline
declare fastcc noundef i32 @g(ptr align 8) #0

define void @ctor() {
entry:
  %v = call i32 @g(ptr null), !my.kind !2
  ret void
}

attributes #0 = { noinline "my-attr"="v1" }

!llvm.module.flags = !{!0}
!my.md = !{!1}

!0 = !{i32 4, !"my.flag", !"v"}
!1 = !{!"x"}
!2 = !{!"tag"}
