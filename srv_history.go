package main

import (
	"io/ioutil"
	"sort"
	"strings"
)

const history_file = "/Users/simo/.bash_history"

var he []string //history entries

func historyInit() error {
	// read the whole file in a []byte
	input, err := ioutil.ReadFile(history_file)
	// typecats to string, then split by newline into a []string
	inputStrings := strings.Split(string(input), "\n")
	// remove duplicates
	he = unique(inputStrings)
	// sort them
	sort.Sort(sort.StringSlice(he))
	return err
}

func historyServe(req request) response {
	log.Warning("reached historyServer")

	//type assert request/response
	log.Notice("Searched: %s", req.Req)

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
