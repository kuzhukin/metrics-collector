package endpoint

/*
	kind:
	 - gauge: float metric
	 - counter: integer metric with accumulating
*/

const (
	// GET: returning all metrics in html
	RootEndpoint = "/"

	// POST: write metric on server in json format
	// example body: {"id": "metric", "type": "gauge",   "value": 10} - for gauge metric
	// example body: {"id": "metric", "type": "counter", "delta": 1}  - for counter metric
	UpdateEndpointJSON = "/update/"
	// POST: write metric on server
	UpdateEndpoint = "/update/{kind}/{name}/{value}"

	// POST: request metric in json format
	// example body: {"id": "metric", "type": "gauge"}
	ValueEndpointJSON = "/value/"
	// GET: returning metric value
	ValueEndpoint = "/value/{kind}/{name}"

	// GET: check database connection
	PingEndpoint = "/ping"
)
