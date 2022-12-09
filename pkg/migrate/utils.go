package migrate

import (
	"fmt"
	"github.com/cespare/xxhash/v2"
	lru "github.com/hashicorp/golang-lru"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/prompb"
)

const (
	labelsCacheSize = 50 * 1000
)

var (
	labelsCache *lru.Cache
	seps        = []byte{'\xff'}
)

func init() {
	var err error
	labelsCache, err = lru.New(labelsCacheSize)
	if err != nil {
		// Error only when size < 0.
		lg.Errorf("error creating label cache: %s", err.Error())
	}
}

func CreatePrombQuery(mint, maxt int64, matchers []*labels.Matcher) (*prompb.Query, error) {
	ms, err := toLabelMatchers(matchers)
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}

	return &prompb.Query{
		StartTimestampMs: mint,
		EndTimestampMs:   maxt - 1, // maxt is exclusive while EndTimestampMs is inclusive.
		Matchers:         ms,
	}, nil
}

// HashLabels returns the hash of the provided promb labels-set. It fetches the value of the hash from the labelsCache if
// exists, else creates the hash and returns it after storing the new hash value.
func HashLabels(lset prompb.Labels) uint64 {
	if hash, found := labelsCache.Get(lset.String()); found {
		return hash.(uint64)
	}

	b := make([]byte, 0, 1024)

	for i, v := range lset.Labels {
		if len(b)+len(v.Name)+len(v.Value)+2 >= cap(b) {
			// If labels entry is 1KB+ do not allocate whole entry.
			h := xxhash.New()
			_, _ = h.Write(b)
			for _, v := range lset.Labels[i:] {
				_, _ = h.WriteString(v.Name)
				_, _ = h.Write(seps)
				_, _ = h.WriteString(v.Value)
				_, _ = h.Write(seps)
			}
			return h.Sum64()
		}

		b = append(b, v.Name...)
		b = append(b, seps[0])
		b = append(b, v.Value...)
		b = append(b, seps[0])
	}

	sum := xxhash.Sum64(b)
	str := lset.String()
	labelsCache.Add(str, sum)

	return sum
}

// LabelSet creates a new label_set for the provided metric name and job name.
//func LabelSet(metricName, migrationJobName string) []prompb.Label {
//	return []prompb.Label{
//		{Name: labels.MetricName, Value: metricName},
//		{Name: LabelJob, Value: migrationJobName},
//	}
//}

func toLabelMatchers(matchers []*labels.Matcher) ([]*prompb.LabelMatcher, error) {
	pbMatchers := make([]*prompb.LabelMatcher, 0, len(matchers))

	for _, m := range matchers {
		var mType prompb.LabelMatcher_Type

		switch m.Type {
		case labels.MatchEqual:
			mType = prompb.LabelMatcher_EQ
		case labels.MatchNotEqual:
			mType = prompb.LabelMatcher_NEQ
		case labels.MatchRegexp:
			mType = prompb.LabelMatcher_RE
		case labels.MatchNotRegexp:
			mType = prompb.LabelMatcher_NRE
		default:
			return nil, fmt.Errorf("invalid matcher type")
		}

		pbMatchers = append(pbMatchers, &prompb.LabelMatcher{
			Type:  mType,
			Name:  m.Name,
			Value: m.Value,
		})
	}

	return pbMatchers, nil
}

func GetLabelsStr(lbls []prompb.Label) (ls string) {
	for i := 0; i < len(lbls); i++ {
		ls += fmt.Sprintf("%s|%s|", lbls[i].Name, lbls[i].Value)
	}

	if len(lbls) > 1 {
		ls = ls[:len(ls)-1] // Ignore the ending '|'.
	}

	return
}