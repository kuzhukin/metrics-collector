package endpoint

/*
	kind:
	 - gauge: float metric
	 - counter: integer metric with accumulating
*/

const (
	// GET: returning all metrics in html
	RootEndpoint = "/"
	// POST: write metric on server
	UpdateEndpoint = "/update/{kind}/{name}/{value}"
	// GET: returning metric value
	ValueEndpoint = "/value/{kind}/{name}"
)
