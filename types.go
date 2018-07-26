package main

// MICROSERVER

type microServer struct {
	name string
	run  func(request) response
}

// REQUEST/RESPPONSE

type request struct {
	Req string
}

type response struct {
	Res  string
	List []string
}
