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
set runcmd {mplayer -quiet ` + filepath.Join(filepath.Dir(os.Args[0]), "cling.mp3") + `}
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
label .fC.l -text {face BOT:}
entry .fC.e  -textvariable fCval
button .fC.b -text {submit} -command {go::face BOT $fCval}
pack .fC.l .fC.e .fC.b -side left

button .cA -text {recenter A} -command {go::recenter A}
button .cB -text {recenter B} -command {go::recenter B}
pack .cA .cB

frame .nst
pack .nst

frame .nst.nod -relief groove -borderwidth 3 -padx 4 -pady 4
frame .nst.shake -relief groove -borderwidth 3 -padx 4 -pady 4
frame .nst.tilt -relief groove -borderwidth 3 -padx 4 -pady 4

pack .nst.nod .nst.shake .nst.tilt -side left -padx 4 -pady 4

label .nst.nod.title -text {Nod}
pack .nst.nod.title

label .nst.shake.title -text {Shake}
pack .nst.shake.title

label .nst.tilt.title -text {Tilt}
pack .nst.tilt.title

frame .nst.nod.nAB
pack .nst.nod.nAB
label .nst.nod.nAB.l -text {A sees B:}
set nodAB 1
entry .nst.nod.nAB.e -textvariable nodAB
button .nst.nod.nAB.b -text {submit} -command {go::nod A B $nodAB}
pack .nst.nod.nAB.l .nst.nod.nAB.e .nst.nod.nAB.b -side left

frame .nst.nod.nAC
pack .nst.nod.nAC
label .nst.nod.nAC.l -text {A sees BOT:}
set nodAC 1
entry .nst.nod.nAC.e -textvariable nodAC
button .nst.nod.nAC.b -text {submit} -command {go::nod A BOT $nodAC}
pack .nst.nod.nAC.l .nst.nod.nAC.e .nst.nod.nAC.b -side left

frame .nst.nod.nBA
pack .nst.nod.nBA
label .nst.nod.nBA.l -text {B sees A:}
set nodBA 1
entry .nst.nod.nBA.e -textvariable nodBA
button .nst.nod.nBA.b -text {submit} -command {go::nod B A $nodBA}
pack .nst.nod.nBA.l .nst.nod.nBA.e .nst.nod.nBA.b -side left

frame .nst.nod.nBC
pack .nst.nod.nBC
label .nst.nod.nBC.l -text {B sees BOT:}
set nodBC 1
entry .nst.nod.nBC.e -textvariable nodBC
button .nst.nod.nBC.b -text {submit} -command {go::nod B BOT $nodBC}
pack .nst.nod.nBC.l .nst.nod.nBC.e .nst.nod.nBC.b -side left

frame .nst.nod.nCA
pack .nst.nod.nCA
label .nst.nod.nCA.l -text {BOT sees A:}
set nodCA 1
entry .nst.nod.nCA.e -textvariable nodCA
button .nst.nod.nCA.b -text {submit} -command {go::nod BOT A $nodCA}
pack .nst.nod.nCA.l .nst.nod.nCA.e .nst.nod.nCA.b -side left

frame .nst.nod.nCB
pack .nst.nod.nCB
label .nst.nod.nCB.l -text {BOT sees B:}
set nodCB 1
entry .nst.nod.nCB.e -textvariable nodCB
button .nst.nod.nCB.b -text {submit} -command {go::nod BOT B $nodCB}
pack .nst.nod.nCB.l .nst.nod.nCB.e .nst.nod.nCB.b -side left

frame .nst.nod.r
pack .nst.nod.r
label .nst.nod.r.l -text {global:}
set nodvalue 1
entry .nst.nod.r.e -textvariable nodvalue
button .nst.nod.r.b -text {submit} -command {setglobalnod $nodvalue}
pack .nst.nod.r.l .nst.nod.r.e .nst.nod.r.b -side left

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



frame .nst.shake.nAB
pack .nst.shake.nAB
label .nst.shake.nAB.l -text {A sees B:}
set shakeAB 1
entry .nst.shake.nAB.e -textvariable shakeAB
button .nst.shake.nAB.b -text {submit} -command {go::shake A B $shakeAB}
pack .nst.shake.nAB.l .nst.shake.nAB.e .nst.shake.nAB.b -side left

frame .nst.shake.nAC
pack .nst.shake.nAC
label .nst.shake.nAC.l -text {A sees BOT:}
set shakeAC 1
entry .nst.shake.nAC.e -textvariable shakeAC
button .nst.shake.nAC.b -text {submit} -command {go::shake A BOT $shakeAC}
pack .nst.shake.nAC.l .nst.shake.nAC.e .nst.shake.nAC.b -side left

frame .nst.shake.nBA
pack .nst.shake.nBA
label .nst.shake.nBA.l -text {B sees A:}
set shakeBA 1
entry .nst.shake.nBA.e -textvariable shakeBA
button .nst.shake.nBA.b -text {submit} -command {go::shake B A $shakeBA}
pack .nst.shake.nBA.l .nst.shake.nBA.e .nst.shake.nBA.b -side left

frame .nst.shake.nBC
pack .nst.shake.nBC
label .nst.shake.nBC.l -text {B sees BOT:}
set shakeBC 1
entry .nst.shake.nBC.e -textvariable shakeBC
button .nst.shake.nBC.b -text {submit} -command {go::shake B BOT $shakeBC}
pack .nst.shake.nBC.l .nst.shake.nBC.e .nst.shake.nBC.b -side left

frame .nst.shake.nCA
pack .nst.shake.nCA
label .nst.shake.nCA.l -text {BOT sees A:}
set shakeCA 1
entry .nst.shake.nCA.e -textvariable shakeCA
button .nst.shake.nCA.b -text {submit} -command {go::shake BOT A $shakeCA}
pack .nst.shake.nCA.l .nst.shake.nCA.e .nst.shake.nCA.b -side left

frame .nst.shake.nCB
pack .nst.shake.nCB
label .nst.shake.nCB.l -text {BOT sees B:}
set shakeCB 1
entry .nst.shake.nCB.e -textvariable shakeCB
button .nst.shake.nCB.b -text {submit} -command {go::shake BOT B $shakeCB}
pack .nst.shake.nCB.l .nst.shake.nCB.e .nst.shake.nCB.b -side left

frame .nst.shake.r
pack .nst.shake.r
label .nst.shake.r.l -text {global:}
set shakevalue 1
entry .nst.shake.r.e -textvariable shakevalue
button .nst.shake.r.b -text {submit} -command {setglobalshake $shakevalue}
pack .nst.shake.r.l .nst.shake.r.e .nst.shake.r.b -side left

proc setglobalshake args {
    global shakeAB shakeAC shakeBA shakeBC shakeCA shakeCB
    set shake [lindex $args 0]
    set shakeAB $shake
    set shakeAC $shake
    set shakeBA $shake
    set shakeBC $shake
    set shakeCA $shake
    set shakeCB $shake
    go::globalshake $shake
}



frame .nst.tilt.nAB
pack .nst.tilt.nAB
label .nst.tilt.nAB.l -text {A sees B:}
set tiltAB 1
entry .nst.tilt.nAB.e -textvariable tiltAB
button .nst.tilt.nAB.b -text {submit} -command {go::tilt A B $tiltAB}
pack .nst.tilt.nAB.l .nst.tilt.nAB.e .nst.tilt.nAB.b -side left

frame .nst.tilt.nAC
pack .nst.tilt.nAC
label .nst.tilt.nAC.l -text {A sees BOT:}
set tiltAC 1
entry .nst.tilt.nAC.e -textvariable tiltAC
button .nst.tilt.nAC.b -text {submit} -command {go::tilt A BOT $tiltAC}
pack .nst.tilt.nAC.l .nst.tilt.nAC.e .nst.tilt.nAC.b -side left

frame .nst.tilt.nBA
pack .nst.tilt.nBA
label .nst.tilt.nBA.l -text {B sees A:}
set tiltBA 1
entry .nst.tilt.nBA.e -textvariable tiltBA
button .nst.tilt.nBA.b -text {submit} -command {go::tilt B A $tiltBA}
pack .nst.tilt.nBA.l .nst.tilt.nBA.e .nst.tilt.nBA.b -side left

frame .nst.tilt.nBC
pack .nst.tilt.nBC
label .nst.tilt.nBC.l -text {B sees BOT:}
set tiltBC 1
entry .nst.tilt.nBC.e -textvariable tiltBC
button .nst.tilt.nBC.b -text {submit} -command {go::tilt B BOT $tiltBC}
pack .nst.tilt.nBC.l .nst.tilt.nBC.e .nst.tilt.nBC.b -side left

frame .nst.tilt.nCA
pack .nst.tilt.nCA
label .nst.tilt.nCA.l -text {BOT sees A:}
set tiltCA 1
entry .nst.tilt.nCA.e -textvariable tiltCA
button .nst.tilt.nCA.b -text {submit} -command {go::tilt BOT A $tiltCA}
pack .nst.tilt.nCA.l .nst.tilt.nCA.e .nst.tilt.nCA.b -side left

frame .nst.tilt.nCB
pack .nst.tilt.nCB
label .nst.tilt.nCB.l -text {BOT sees B:}
set tiltCB 1
entry .nst.tilt.nCB.e -textvariable tiltCB
button .nst.tilt.nCB.b -text {submit} -command {go::tilt BOT B $tiltCB}
pack .nst.tilt.nCB.l .nst.tilt.nCB.e .nst.tilt.nCB.b -side left

frame .nst.tilt.r
pack .nst.tilt.r
label .nst.tilt.r.l -text {global:}
set tiltvalue 1
entry .nst.tilt.r.e -textvariable tiltvalue
button .nst.tilt.r.b -text {submit} -command {setglobaltilt $tiltvalue}
pack .nst.tilt.r.l .nst.tilt.r.e .nst.tilt.r.b -side left

proc setglobaltilt args {
    global tiltAB tiltAC tiltBA tiltBC tiltCA tiltCB
    set tilt [lindex $args 0]
    set tiltAB $tilt
    set tiltAC $tilt
    set tiltBA $tilt
    set tiltBC $tilt
    set tiltCA $tilt
    set tiltCB $tilt
    go::globaltilt $tilt
}


frame .t
pack .t

frame .t.a -relief groove -borderwidth 3 -padx 4 -pady 4
frame .t.b -relief groove -borderwidth 3 -padx 4 -pady 4
frame .t.c -relief groove -borderwidth 3 -padx 4 -pady 4

pack .t.a .t.b .t.c -side left -padx 4 -pady 4

set abc off
set acb off
label .t.a.title -text {A sees...}
checkbutton .t.a.bc -variable abc -onvalue on -offvalue off -text {B looking at BOT} -command {go::turn A B BOT $abc}
checkbutton .t.a.cb -variable acb -onvalue on -offvalue off -text {BOT looking at B} -command {go::turn A BOT B $acb}
pack .t.a.title .t.a.bc .t.a.cb

set bac off
set bca off
label .t.b.title -text {B sees...}
checkbutton .t.b.ac -variable bac -onvalue on -offvalue off -text {A looking at BOT} -command {go::turn B A BOT $bac}
checkbutton .t.b.ca -variable bca -onvalue on -offvalue off -text {BOT looking at A} -command {go::turn B BOT A $bca}
pack .t.b.title .t.b.ac .t.b.ca

set cab off
set cba off
label .t.c.title -text {BOT sees...}
checkbutton .t.c.ab -variable cab -onvalue on -offvalue off -text {A looking at B} -command {go::turn BOT A B $cab}
checkbutton .t.c.ba -variable cba -onvalue on -offvalue off -text {B looking at A} -command {go::turn BOT B A $cba}
pack .t.c.title .t.c.ab .t.c.ba

button .start -text { START } -command { go::start ; destroy .start } -background {#00a000} -foreground white
pack .start -expand yes -fill x -padx 32 -pady 8

# Don't use command exit, because that will kill Go as well
button .q -text {exit} -command {destroy .}
pack .q -side left

`))

	x(tk.RegisterCommand("go::runcommand", func(command string) {
		fmt.Println("Command: run", command)
		chLog <- "C run begin: " + command
		args := strings.Fields(command)
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if w(err) != nil {
			fmt.Println(err)
		}
		chLog <- "C run end: " + command
	}))

	x(tk.RegisterCommand("go::start", func() {
		chCmd <- "start"
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

	x(tk.RegisterCommand("go::shake", func(sees, seen, value string) {
		chCmd <- "shake " + sees + " " + seen + " " + value
	}))

	x(tk.RegisterCommand("go::globalshake", func(s string) {
		chCmd <- "globalshake " + s
	}))

	x(tk.RegisterCommand("go::tilt", func(sees, seen, value string) {
		chCmd <- "tilt " + sees + " " + seen + " " + value
	}))

	x(tk.RegisterCommand("go::globaltilt", func(s string) {
		chCmd <- "globaltilt " + s
	}))

	x(tk.RegisterCommand("go::turn", func(sees, seen, seeing, val string) {
		chCmd <- "turn " + sees + " " + seen + " " + seeing + " " + val
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
