package main

type history_req struct {
	Cmd string
}

type history_resp struct {
	List []string
}
