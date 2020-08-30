package summary

// SummarisedBuild structure for summarising what happened with a build without
// exporting all the log info.
type SummarisedBuild struct {
	BuiltOn         string   `json:"builtOn"`
	Duration        int64    `json:"duration"`
	FullDisplayName string   `json:"fullDisplayName"`
	ID              string   `json:"id"`
	Result          string   `json:"result"`
	Timestamp       int64    `json:"timestamp"`
	URL             string   `json:"url"`
	FailureSummary  string   `json:"failureSummary,omitempty"`
	Problems        []string `json:"problems,omitempty"`
}
