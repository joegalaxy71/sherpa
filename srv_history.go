package main

// SPECS
// follow the history file
// when a change is detected load the delta (shell only append)
// and insert into sqlite db
// use gorm and gorm defaults, so all records will have automatically "created at"

import (
	_ "errors"
	"os"
	"sort"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const HISTORY_FILE = "/Users/simo/.bash_history"

var he []string //history entries
var watcher *fsnotify.Watcher

func historyInit() error {

	go watchHistory()

	return nil
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
					err := updateEntriesDB()
					if err != nil {
						log.Errorf("Error:%s", err)
					}
				}
			case err := <-watcher.Errors:
				log.Debugf("error: %s", err)
			}
		}
	}()

	err = watcher.Add(HISTORY_FILE)
	if err != nil {
		log.Error(err)
	}
}

func updateEntriesDB() error {
	// read the history file from Histfrom to end of file

	// we open the file every time so we don't leave any locks around
	file, err := os.Open(HISTORY_FILE)
	if err != nil {
		return err
	}

	// get the size
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	size := fileInfo.Size()

	// Whence is the point of reference for offset
	// 0 = Beginning of file
	// 1 = Current position
	// 2 = End of file
	var whence = 0

	var chunk uint32 = 0
	var readFrom uint32 = 0

	//check if size is bigger that the read part or
	if status.HistFrom == 0 {
		log.Noticef("history file never read before")
		// the chunk we read == file size (real all file)
		chunk = uint32(size)
		readFrom = 0
	} else if uint32(size) < status.HistFrom {
		//if read part = 0 (never read before)
		log.Noticef("history file size < of part already read, reading whole file")
		// the chunk we read == file size (real all file)
		chunk = uint32(size)
		readFrom = 0
	} else {
		//calculate chunk to read
		log.Noticef("history file size > of part already read, reading delta")

		chunk = uint32(size) - status.HistFrom
		readFrom = status.HistFrom

	}

	// we'll close it when done
	defer file.Close()

	// seek to HistFrom

	newPosition, err := file.Seek(int64(readFrom), whence)
	log.Debugf("New position is: %v", newPosition)
	if err != nil {
		return err
	}

	// The file.Read() function will happily read a tiny file in to a large
	// byte slice, but io.ReadFull() will return an
	// error if the file is smaller than the byte slice.
	byteSlice := make([]byte, chunk)
	numBytesRead, err := file.Read(byteSlice)
	log.Debugf("Read %v bytes from %s", numBytesRead, fileInfo.Name())
	if err != nil {
		return err
	}

	// update status whit the info about the new read part
	status.HistFrom = uint32(size)

	//always keep the status in a global var

	db.First(&status, "one = ?", "one")

	db.Model(&status).Update("HistFrom", uint32(size))

	//db.Save(&status)

	log.Noticef("%+v", db)

	// typecats to string, then split by newline into a []string
	inputStrings := strings.Split(string(byteSlice), "\n")
	// remove duplicates
	he = unique(inputStrings)
	// sort them
	sort.Sort(sort.StringSlice(he))

	//TODO: what happens if the history file gets destroyed or modified?
	// should we treat it as a new file?

	// then we split what we've read into strings
	// and we insert in the db

	return nil
}
