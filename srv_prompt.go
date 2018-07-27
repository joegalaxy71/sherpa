package main

func promptInit() error {
	return nil
}

func promptServe(req request) response {
	log.Warningf("reached promptServer")

	//type assert request/response
	log.Noticef("Searched: %s", req.Req)

	var res response

	// actual work done
	res.Prompts = append(res.Prompts, prompt{"base", "edrfrgt"}, prompt{"other", "koko"})

	return res
}

func promptCleanup() error {
	return nil
}
