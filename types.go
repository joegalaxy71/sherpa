package main

// MICROSERVER

type microServer struct {
	name string
	run  func(request) response
	res  response
	req  request
}

// HISTORY

type historyReq struct {
	Req string
}

type historyRes struct {
	Res  string
	List []string
}

// CLEANUP

type cleanupReq struct {
	Cmd string
}

type cleanupRes struct {
	Res string
}

// INTERFACES

type request interface {
	req() string
}

type response interface {
	res() string
}

// IMPLEMENTATIONS

func (hreq historyReq) req() string {
	return hreq.Req
}

func (hres historyRes) res() string {
	return hres.Res
}
