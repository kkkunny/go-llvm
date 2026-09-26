; ModuleID = 'eh'
source_filename = "eh"

@ti = external global ptr

declare i32 @pers(i32)

declare i32 @g(i32)

declare void @h()

define i32 @itanium(i32 %0) personality ptr @pers {
entry:
  %v = invoke i32 @g(i32 %0)
          to label %cont unwind label %lpad

cont:                                             ; preds = %entry
  ret i32 %v

lpad:                                             ; preds = %entry
  %lp = landingpad { ptr, i32 }
          cleanup
          catch ptr @ti
  resume { ptr, i32 } %lp
}

define void @funclets() personality ptr @pers {
entry:
  invoke void @h()
          to label %cont unwind label %dispatch

cont:                                             ; preds = %entry
  ret void

dispatch:                                         ; preds = %entry
  %cs = catchswitch within none [label %hnd] unwind to caller

hnd:                                              ; preds = %dispatch
  %cp = catchpad within %cs [ptr @ti]
  catchret from %cp to label %done

done:                                             ; preds = %hnd
  ret void
}

define void @cleanup() personality ptr @pers {
entry:
  invoke void @h()
          to label %cont unwind label %clean

cont:                                             ; preds = %entry
  ret void

clean:                                            ; preds = %entry
  %clp = cleanuppad within none []
  cleanupret from %clp unwind to caller
}
