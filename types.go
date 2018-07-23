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
