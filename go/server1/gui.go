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
button .q -text {exit} -command {go::finish; exit}
pack .q -side left
`))

	x(tk.RegisterCommand("go::finish", finish))

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
