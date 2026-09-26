; ModuleID = 'loop'
source_filename = "loop"

define i32 @sum_to(i32 %0) {
entry:
  %acc = alloca i32, align 4
  store i32 0, ptr %acc, align 4
  br label %loop

loop:                                             ; preds = %loop, %entry
  %i = phi i32 [ 0, %entry ], [ %next, %loop ]
  %cur = load i32, ptr %acc, align 4
  %next = add i32 %i, 1
  store i32 %next, ptr %acc, align 4
  %done = icmp sge i32 %next, %0
  br i1 %done, label %exit, label %loop

exit:                                             ; preds = %loop
  %result = load i32, ptr %acc, align 4
  ret i32 %result
}
