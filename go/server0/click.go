package main

var (
	clickHandlers = map[string]func(int, int){
		"":     clickNone,
		"none": clickNone,
		"demo": clickDemo,
	}
	clickHandle = clickHandlers[""]
)

func clickNone(from, to int) {
}

func clickDemo(from, to int) {
	if to < 0 {
		return
	}
	infoMakeChoice(from, users[from].uid+"-clicked-"+users[to].uid, "Yes", "No", []string{"Do you like this one?"})
}
