package main

import (
	"io/ioutil"
	"sort"
	"strings"
)

const history_file = "/Users/simo/.bash_history"

var he []string //history entries

func historyInit() error {
	input, err := ioutil.ReadFile(history_file)
	inputStrings := strings.Split(string(input), "\n")
	he = unique(inputStrings)
	sort.Sort(sort.StringSlice(he))
	return err
}

func historyServe(req request) response {
	log.Warning("reached historyServer")

	//type assert request/response
	log.Noticef("Searched: %s", req.Req)

	var res response

	// actual work done
	res.List = append(res.List, "zfs list", "zfs list -t snap", "zfs list -t snap -o name")

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
