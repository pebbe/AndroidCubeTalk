package main

import (
	"github.com/nsf/gothic"

	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
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
		x(tk.Eval(format(`
set f_[[0]]_val [[1]]
frame .f_[[0]]
pack .f_[[0]]
label .f_[[0]].l -text {face [[0]]:}
entry .f_[[0]].e  -textvariable f_[[0]]_val
button .f_[[0]].b -text {submit} -command {go::face [[0]] $f_[[0]]_val}
pack .f_[[0]].l .f_[[0]].e .f_[[0]].b -side left
`,
			lbl, cubes[idx].face)))
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
				x(tk.Eval(format(`
frame .nst.[[0]].n_[[1]]_[[2]]
pack .nst.[[0]].n_[[1]]_[[2]]
label .nst.[[0]].n_[[1]]_[[2]].l -text {[[1]] sees [[2]]:}
set [[0]]_[[1]]_[[2]] 1
entry .nst.[[0]].n_[[1]]_[[2]].e -textvariable [[0]]_[[1]]_[[2]]
button .nst.[[0]].n_[[1]]_[[2]].b -text {submit} -command {go::[[0]] [[1]] [[2]] $[[0]]_[[1]]_[[2]]}
pack .nst.[[0]].n_[[1]]_[[2]].l .nst.[[0]].n_[[1]]_[[2]].e .nst.[[0]].n_[[1]]_[[2]].b -side left
`, item, li, lj)))
			}
		}
		x(tk.Eval(format(`
frame .nst.[[0]].r
pack .nst.[[0]].r
label .nst.[[0]].r.l -text {global:}
set [[0]]value 1
entry .nst.[[0]].r.e -textvariable [[0]]value
button .nst.[[0]].r.b -text {submit} -command {setglobal[[0]] $[[0]]value}
pack .nst.[[0]].r.l .nst.[[0]].r.e .nst.[[0]].r.b -side left

proc setglobal[[0]] args {
    global [[1]]
    set v [lindex $args 0]
[[2]]
    go::global[[0]] $v
}
`,
			item,
			strings.Join(globals, " "),
			strings.Join(sets, "\n"))))

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
				if j == k {
					continue
				}
				lk := cubes[k].uid
				win := fmt.Sprintf(".t.w%s.w%s_%s", li, lj, lk)
				combis = append(combis, win)
				x(tk.Eval(format(`
set sees_[[0]]_[[1]]_[[0]] off
checkbutton .t.w[[0]].w[[1]]_[[2]] -variable sees_[[0]]_[[1]]_[[2]] -onvalue on -offvalue off -text {[[1]] looking at [[2]]} -command {go::turn [[0]] [[1]] [[2]] $sees_[[0]]_[[1]]_[[2]]}
`, li, lj, lk)))
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

func format(t string, args ...interface{}) string {

	t = strings.Replace(t, "{[[", "{ {{- index . ", -1)
	t = strings.Replace(t, "[[", "{{index . ", -1)
	t = strings.Replace(t, "]]}", " -}} }", -1)
	t = strings.Replace(t, "]]", "}}", -1)

	tmp := template.Must(template.New("tmp").Parse(t))

	var buf bytes.Buffer
	x(tmp.Execute(&buf, args))

	return buf.String()
}
