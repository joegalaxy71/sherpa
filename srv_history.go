package main

// SPECS
// follow the history file
// when a change is detected load the delta (shell only append)
// and insert into sqlite db
// use gorm and gorm defaults, so all records will have automatically "created at"

import (
	"os"
	"strings"
	"time"

	//"runtime"
	// "time"

	"github.com/fsnotify/fsnotify"
)

//const historyFile = "/root/.bash_history"
var historyFile string
var watcher *fsnotify.Watcher

func historyInit() error {

	historyFile = homedir + "/.bash_history"

	watchHistory()

	return nil
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

	//we launch manually the first db update
	updateEntriesDB()

	err = watcher.Add(historyFile)
	if err != nil {
		log.Error(err)
	}
}

func updateEntriesDB() error {
	// read the history file from Histfrom to end of file

	var err error

	// we open the file every time so we don't leave any locks around
	file, err := os.Open(historyFile)
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

	var chunk int64
	var readFrom int64

	//check if size is bigger that the read part or
	if status.HistFrom == 0 {
		log.Debugf("history file never read before")
		// the chunk we read == file size (real all file)
		chunk = size
		readFrom = 0
	} else if size < status.HistFrom {
		//if read part = 0 (never read before)
		log.Debugf("history file size < of part already read, reading whole file")
		// the chunk we read == file size (real all file)
		chunk = size
		readFrom = 0
	} else {
		//calculate chunk to read
		log.Debugf("history file size > of part already read, reading delta")

		chunk = size - status.HistFrom
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

	//s := string(byteSlice[:])

	//log.Warningf("Read: %s", s)
	// update status whit the info about the new read part
	status.HistFrom = size

	//always keep the status in a global var

	//db.First(&status, "one = ?", "one")

	//db.Model(&status).Update("HistFrom", uint32(size))

	//status.HistFrom = size

	//log.Debugf("%+v", db)

	//start := time.Now()

	// typecats to string, then split by newline into a []string
	inputStrings := strings.Split(string(byteSlice), "\n")
	log.Debugf("inputstsrings len()=%v, cap()=%v", len(inputStrings), cap(inputStrings))

	// no need to sort (it's a db work) or remove duplicates (they're needed)
	// write them to the db

	var cloudUpdated = true

	for _, entry := range inputStrings {
		if entry != "" {
			//log.Debugf("entry=%s, host=%s", entry, "retina")

			var he historyNew
			he.Account = "account"
			he.Entry = entry
			he.UserAtHost = currentUser.Username + "@" + hostName

			var res response

			err = cec.Request("history-new", he, &res, 1000*time.Millisecond)
			if err != nil {
				cloudUpdated = false
			}
		}
	}

	if cloudUpdated == true {
		db.Save(&status)
	}

	test()

	return nil
}

func test() {
}
