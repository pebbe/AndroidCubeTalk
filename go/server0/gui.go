package main

import (
	"github.com/nsf/gothic"

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
label .sizes.l -text {size of cubes, width / height / depth}
set cubew 1
set cubeh 1
set cubed 1
entry .sizes.w -textvariable cubew
entry .sizes.h -textvariable cubeh
entry .sizes.d -textvariable cubed
button .sizes.b -text {submit} -command {go::cubesize $cubew $cubeh $cubed}
pack .sizes.l .sizes.w .sizes.h .sizes.d .sizes.b -side left

button .cA -text {recenter A} -command {go::recenter A}
button .cB -text {recenter B} -command {go::recenter B}
button .cC -text {recenter C} -command {go::recenter C}
button .cD -text {recenter D} -command {go::recenter D}
button .cE -text {recenter E} -command {go::recenter E}
button .cF -text {recenter F} -command {go::recenter F}
pack .cA .cB .cC .cD .cE .cF -expand yes -fill both

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
