; ModuleID = 'vec'
source_filename = "vec"

define i32 @vec_demo(<4 x i32> %0) {
entry:
  %ins = insertelement <4 x i32> %0, i32 42, i32 1
  %sh = shufflevector <4 x i32> %ins, <4 x i32> %ins, <4 x i32> zeroinitializer
  %ex = extractelement <4 x i32> %sh, i32 0
  ret i32 %ex
}
