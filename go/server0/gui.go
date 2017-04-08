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

	using := make([]int, 0)
	for i, cube := range cubes {
		if cube.gui {
			using = append(using, i)
		}
	}

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
`))

	for _, idx := range using {
		lbl := cubes[idx].uid
		x(tk.Eval(fmt.Sprintf(`
set f_%s_val %d
frame .f_%s
pack .f_%s
label .f_%s.l -text {face %s:}
entry .f_%s.e  -textvariable f_%s_val
button .f_%s.b -text {submit} -command {go::face %s $f_%s_val}
pack .f_%s.l .f_%s.e .f_%s.b -side left
`,
			lbl, cubes[idx].face, lbl, lbl, lbl, lbl, lbl, lbl, lbl, lbl, lbl, lbl, lbl, lbl)))
	}

	ww := make([]string, 0, len(using))
	for idx := range cubes {
		lbl := cubes[idx].uid
		ww = append(ww, ".c_"+lbl)
		x(tk.Eval(fmt.Sprintf(`
button .c_%s -text {recenter %s} -command {go::recenter %s}
`, lbl, lbl, lbl)))
	}
	x(tk.Eval("pack " + strings.Join(ww, " ") + "\n"))

	x(tk.Eval(`
button .cc -text {recenter all} -command {go::recenter_all}
pack .cc

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
`))

	for _, item := range []string{"nod", "shake", "tilt"} {

		globals := make([]string, 0)
		sets := make([]string, 0)
		for _, i := range using {
			li := cubes[i].uid
			for _, j := range using {
				if i == j {
					continue
				}
				lj := cubes[j].uid
				globals = append(globals, fmt.Sprintf("%s_%s_%s", item, li, lj))
				sets = append(sets, fmt.Sprintf("set %s_%s_%s $v", item, li, lj))
				x(tk.Eval(fmt.Sprintf(`
frame .nst.%s.n_%s_%s
pack .nst.%s.n_%s_%s
label .nst.%s.n_%s_%s.l -text {%s sees %s:}
set %s_%s_%s 1
entry .nst.%s.n_%s_%s.e -textvariable %s_%s_%s
button .nst.%s.n_%s_%s.b -text {submit} -command {go::%s %s %s $%s_%s_%s}
pack .nst.%s.n_%s_%s.l .nst.%s.n_%s_%s.e .nst.%s.n_%s_%s.b -side left
`,
					item, li, lj,
					item, li, lj,
					item, li, lj, li, lj,
					item, li, lj,
					item, li, lj, item, li, lj,
					item, li, lj, item, li, lj, item, li, lj,
					item, li, lj, item, li, lj, item, li, lj)))
			}
		}
		x(tk.Eval(fmt.Sprintf(`
frame .nst.%s.r
pack .nst.%s.r
label .nst.%s.r.l -text {global:}
set %svalue 1
entry .nst.%s.r.e -textvariable %svalue
button .nst.%s.r.b -text {submit} -command {setglobal%s $%svalue}
pack .nst.%s.r.l .nst.%s.r.e .nst.%s.r.b -side left

proc setglobal%s args {
    global %s
    set v [lindex $args 0]
%s
    go::global%s $v
}
`,
			item,
			item,
			item,
			item,
			item, item,
			item, item, item,
			item, item, item,
			item,
			strings.Join(globals, " "),
			strings.Join(sets, "\n"),
			item)))

	}

	x(tk.Eval(`

frame .t
pack .t
`))

	frames := make([]string, 0)
	for _, i := range using {
		uid := cubes[i].uid
		frames = append(frames, ".t.w"+uid)
		x(tk.Eval("frame .t.w" + uid + " -relief groove -borderwidth 3 -padx 4 -pady 4\n"))
	}

	if len(frames) > 0 {
		x(tk.Eval(`

pack ` + strings.Join(frames, " ") + ` -side left -padx 4 -pady 4
`))
	}

	for _, i := range using {
		li := cubes[i].uid
		x(tk.Eval("label .t.w" + li + ".title -text {" + li + " sees...}\n"))
		combis := make([]string, 0)
		for _, j := range using {
			if i == j {
				continue
			}
			lj := cubes[j].uid
			for _, k := range using {
				if i == k || j == k {
					continue
				}
				lk := cubes[k].uid
				win := fmt.Sprintf(".t.w%s.w%s_%s", li, lj, lk)
				combis = append(combis, win)
				x(tk.Eval(fmt.Sprintf(`
set sees_%s_%s_%s off
checkbutton .t.w%s.w%s_%s -variable sees_%s_%s_%s -onvalue on -offvalue off -text {%s looking at %s} -command {go::turn %s %s %s $sees_%s_%s_%s}
`,
					li, lj, lk,
					li, lj, lk,
					li, lj, lk,
					lj, lk,
					li, lj, lk,
					li, lj, lk)))
			}
		}
		x(tk.Eval(fmt.Sprintf("pack .t.w%s.title %s\n", li, strings.Join(combis, " "))))

	}

	x(tk.Eval(`

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

	x(tk.RegisterCommand("go::recenter_all", func() {
		chCmd <- "recenter_all"
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
