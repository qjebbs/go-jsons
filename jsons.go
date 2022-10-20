package jsons

// stdMerger is the standard json files merger
var stdMerger = NewMerger()

func init() {
	must(stdMerger.RegisterDefaultLoader())
}

// Merge merges inputs into a single json.
//
// It detects the format by file extension, or try all mergers
// if no extension found
//
// Accepted Input:
//
//  - `[]byte`: content of a file
//  - `string`: path to a file, either local or remote
//  - `[]string`: a list of files, either local or remote
//  - `io.Reader`: a file content reader
func Merge(inputs ...interface{}) ([]byte, error) {
	return stdMerger.Merge(inputs...)
}

// MergeAs loads inputs of the specific format and merges into a single json.
//
// Accepted Input:
//
//  - `[]byte`: content of a file
//  - `string`: path to a file, either local or remote
//  - `[]string`: a list of files, either local or remote
//  - `io.Reader`: a file content reader
func MergeAs(format Format, inputs ...interface{}) ([]byte, error) {
	return stdMerger.MergeAs(format, inputs...)
}
