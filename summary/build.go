package summary

type SummarisedBuild struct {
	BuiltOn         string `json:"builtOn"`
	Duration        int64  `json:"duration"`
	FullDisplayName string `json:"fullDisplayName"`
	ID              string `json:"id"`
	Result          string `json:"result"`
	Timestamp       int64  `json:"timestamp"`
	URL             string `json:"url"`
}
