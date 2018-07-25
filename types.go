package main

type command struct {
	From   string
	Action string
	Arg    string
}

type response struct {
	To     string
	Result string
}

type history_req struct {
	Cmd string
}

type history_resp struct {
	List []string
}
