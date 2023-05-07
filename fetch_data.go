package main

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
)

const DATA_URL = "https://ec.europa.eu/eurostat/estat-navtree-portlet-prod/BulkDownloadListing?file=data/demo_r_mwk_05.tsv.gz"

// ReadData makes a HTTP requests to fetch gzipped TSV data from Eurostat website.
// Returns the TSV data as a string.
func ReadData() (string, error) {
	var data string

	resp, err := http.Get(DATA_URL)
	if err != nil {
		return data, err
	}

	gzipBody, err := gzip.NewReader(resp.Body)
	if err != nil {
		return data, err
	}

	tsvData, err := ioutil.ReadAll(gzipBody)
	if err != nil {
		return data, err
	}

	return string(tsvData), nil
}
