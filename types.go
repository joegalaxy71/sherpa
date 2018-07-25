package main

type historyReq struct {
	Cmd string
}

type historyRes struct {
	List []string
}

// CLEANUP

type cleanupReq struct {
	Cmd string
}

type cleanupRes struct {
	Res string
}
