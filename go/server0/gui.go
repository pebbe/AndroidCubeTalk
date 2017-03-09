package main

import (
	"github.com/nsf/gothic"

	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	tk *gothic.Interpreter
)

func gui() {

	tk = gothic.NewInterpreter("namespace eval go {}")

	x(tk.Eval(`

wm title . "` + tclquote(filepath.Base(os.Args[0])) + `"

frame .sizes
pack .sizes
label .sizes.l -text {size of cubes, width / height / depth:}
set cubew 1
set cubeh 1
set cubed 1
entry .sizes.w -textvariable cubew
entry .sizes.h -textvariable cubeh
entry .sizes.d -textvariable cubed
button .sizes.b -text {submit} -command {go::cubesize $cubew $cubeh $cubed}
pack .sizes.l .sizes.w .sizes.h .sizes.d .sizes.b -side left

set fAval ` + fmt.Sprint(cubes[0].face) + `
frame .fA
pack .fA
label .fA.l -text {face A:}
entry .fA.e  -textvariable fAval
button .fA.b -text {submit} -command {go::face A $fAval}
pack .fA.l .fA.e .fA.b -side left

set fBval ` + fmt.Sprint(cubes[1].face) + `
frame .fB
pack .fB
label .fB.l -text {face B:}
entry .fB.e  -textvariable fBval
button .fB.b -text {submit} -command {go::face B $fBval}
pack .fB.l .fB.e .fB.b -side left

set fCval ` + fmt.Sprint(cubes[2].face) + `
frame .fC
pack .fC
label .fC.l -text {face C:}
entry .fC.e  -textvariable fCval
button .fC.b -text {submit} -command {go::face C $fCval}
pack .fC.l .fC.e .fC.b -side left

set fDval ` + fmt.Sprint(cubes[3].face) + `
frame .fD
pack .fD
label .fD.l -text {face D:}
entry .fD.e  -textvariable fDval
button .fD.b -text {submit} -command {go::face D $fDval}
pack .fD.l .fD.e .fD.b -side left

set fEval ` + fmt.Sprint(cubes[4].face) + `
frame .fE
pack .fE
label .fE.l -text {face E:}
entry .fE.e  -textvariable fEval
button .fE.b -text {submit} -command {go::face E $fEval}
pack .fE.l .fE.e .fE.b -side left

set fFval ` + fmt.Sprint(cubes[5].face) + `
frame .fF
pack .fF
label .fF.l -text {face F:}
entry .fF.e  -textvariable fFval
button .fF.b -text {submit} -command {go::face F $fFval}
pack .fF.l .fF.e .fF.b -side left

button .cA -text {recenter A} -command {go::recenter A}
button .cB -text {recenter B} -command {go::recenter B}
button .cC -text {recenter C} -command {go::recenter C}
button .cD -text {recenter D} -command {go::recenter D}
button .cE -text {recenter E} -command {go::recenter E}
button .cF -text {recenter F} -command {go::recenter F}
pack .cA .cB .cC .cD .cE .cF

frame .r
pack .r
label .r.l -text {global nod enhance:}
set nodvalue 1
entry .r.e -textvariable nodvalue
button .r.b -text {submit} -command {go::globalnod $nodvalue}
pack .r.l .r.e .r.b -side left

# Don't use command exit, because that will kill Go as well
button .q -text {exit} -command {destroy .}
pack .q -side left

`))

	x(tk.RegisterCommand("go::cubesize", func(w, h, d string) {
		chCmd <- "cubesize " + w + " " + h + " " + d
	}))

	x(tk.RegisterCommand("go::recenter", func(s string) {
		chCmd <- "recenter " + s
	}))

	x(tk.RegisterCommand("go::globalnod", func(s string) {
		chCmd <- "globalnod " + s
	}))

	x(tk.RegisterCommand("go::face", func(uid, idx string) {
		chCmd <- "face " + uid + " " + idx
	}))

	<-tk.Done
}

func tclquote(s string) string {
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	s = strings.Replace(s, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	s = strings.Replace(s, "$", "\\$", -1)
	return s
}
