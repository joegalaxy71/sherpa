package main

// SPECS
// follow the history file
// when a change is detected load the delta (shell only append)
// and insert into sqlite db
// use gorm and gorm defaults, so all records will have automatically "created at"

import (
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"sort"
	"strings"
)

const HISTORY_FILE = "/Users/simo/.bash_history"

var he []string //history entries
var watcher *fsnotify.Watcher

func historyInit() error {

	go watchHistory()

	// read the whole file in a []byte
	input, err := ioutil.ReadFile(HISTORY_FILE)
	// typecats to string, then split by newline into a []string
	inputStrings := strings.Split(string(input), "\n")
	// remove duplicates
	he = unique(inputStrings)
	// sort them
	sort.Sort(sort.StringSlice(he))
	return err
}

func historyServe(req request) response {
	log.Warningf("reached historyServer")

	//type assert request/response
	log.Noticef("Searched: %s", req.Req)

	var res response

	for _, entry := range he {
		if strings.Contains(entry, req.Req) {
			// actual work done
			res.List = append(res.List, entry)
		}
	}

	return res
}

func historyCleanup() error {
	return nil
}

/////

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func watchHistory() {
	// we create a watcher to watch the history file
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// we create and launch a goroutine to persist after init
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				//log.Debugf("event: %s", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Debugf("modified file: %s", event.Name)
					updateEntriesDB()
				}
			case err := <-watcher.Errors:
				log.Debugf("error: %s", err)
			}
		}
	}()

	err = watcher.Add(HISTORY_FILE)
	if err != nil {
		log.Fatal(err)
	}
}

func updateEntriesDB() {
	// determine if we've added entries in the past (check status)

	// if yes we read only the new portion of the file

	// else we read the whole file

	//TODO: what happens if the history file gets destroyed or modified?
	// should we treat it as a new file?

	// then we split what we've read into strings
	// and we insert in the db
}
