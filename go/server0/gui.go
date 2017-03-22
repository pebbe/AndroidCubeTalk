package main

import (
	"github.com/nsf/gothic"

	"fmt"
	"os"
	"os/exec"
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

frame .cmd
pack .cmd
label .cmd.l -text {external command:}
set runcmd {mplayer --quiet cling.mp3}
entry .cmd.e -textvariable runcmd -width 40
button .cmd.b -text {run} -command {go::runcommand $runcmd}
pack .cmd.l .cmd.e .cmd.b -side left

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

frame .nAB
pack .nAB
label .nAB.l -text {A sees B nod:}
set nodAB 1
entry .nAB.e -textvariable nodAB
button .nAB.b -text {submit} -command {go::nod A B $nodAB}
pack .nAB.l .nAB.e .nAB.b -side left

frame .nAC
pack .nAC
label .nAC.l -text {A sees C nod:}
set nodAC 1
entry .nAC.e -textvariable nodAC
button .nAC.b -text {submit} -command {go::nod A C $nodAC}
pack .nAC.l .nAC.e .nAC.b -side left

frame .nBA
pack .nBA
label .nBA.l -text {B sees A nod:}
set nodBA 1
entry .nBA.e -textvariable nodBA
button .nBA.b -text {submit} -command {go::nod B A $nodBA}
pack .nBA.l .nBA.e .nBA.b -side left

frame .nBC
pack .nBC
label .nBC.l -text {B sees C nod:}
set nodBC 1
entry .nBC.e -textvariable nodBC
button .nBC.b -text {submit} -command {go::nod B C $nodBC}
pack .nBC.l .nBC.e .nBC.b -side left

frame .nCA
pack .nCA
label .nCA.l -text {C sees A nod:}
set nodCA 1
entry .nCA.e -textvariable nodCA
button .nCA.b -text {submit} -command {go::nod C A $nodCA}
pack .nCA.l .nCA.e .nCA.b -side left

frame .nCB
pack .nCB
label .nCB.l -text {C sees B nod:}
set nodCB 1
entry .nCB.e -textvariable nodCB
button .nCB.b -text {submit} -command {go::nod C B $nodCB}
pack .nCB.l .nCB.e .nCB.b -side left

frame .r
pack .r
label .r.l -text {global nod enhance:}
set nodvalue 1
entry .r.e -textvariable nodvalue
button .r.b -text {submit} -command {setglobalnod $nodvalue}
pack .r.l .r.e .r.b -side left

proc setglobalnod args {
    global nodAB nodAC nodBA nodBC nodCA nodCB
    set nod [lindex $args 0]
    set nodAB $nod
    set nodAC $nod
    set nodBA $nod
    set nodBC $nod
    set nodCA $nod
    set nodCB $nod
    go::globalnod $nod
}

# Don't use command exit, because that will kill Go as well
button .q -text {exit} -command {destroy .}
pack .q -side left

`))

	x(tk.RegisterCommand("go::runcommand", func(command string) {
		fmt.Println("Command: run", command)
		chLog <- "C run begin: " + command
		args := strings.Fields(command)
		cmd := exec.Command(args[0], args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if w(err) != nil {
			fmt.Println(err)
		}
		chLog <- "C run end: " + command
	}))

	x(tk.RegisterCommand("go::cubesize", func(w, h, d string) {
		chCmd <- "cubesize " + w + " " + h + " " + d
	}))

	x(tk.RegisterCommand("go::recenter", func(s string) {
		chCmd <- "recenter " + s
	}))

	x(tk.RegisterCommand("go::nod", func(sees, seen, value string) {
		chCmd <- "nod " + sees + " " + seen + " " + value
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
