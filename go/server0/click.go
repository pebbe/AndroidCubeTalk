package main

func clickHandle(from, to int) {
	infoMakeChoice(from, users[from].uid+"-clicked-"+users[to].uid, "Yes", "No", []string{"Do you like this one?"})
}
