package main

// MICROSERVER

type microServer struct {
	name    string
	init    func() error
	cleanup func() error
}
type UpdateInfo struct {
	BuildNumber string
}
