package migrate

import "github.com/prometheus/prometheus/prompb"

type PrompbResponse struct {
	ID                   int
	Result               *prompb.QueryResult
	NumBytesCompressed   int
	NumBytesUncompressed int
}
