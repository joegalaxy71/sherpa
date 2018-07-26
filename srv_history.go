package main

func historyServer(req request) response {
	log.Warning("reached historyServer")

	//type assert request/response
	log.Notice("Searched: %s", req.Req)

	var res response

	// actual work done
	res.List = append(res.List, "zfs list", "zfs list -t snap", "zfs list -t snap -o name")

	return res
}
