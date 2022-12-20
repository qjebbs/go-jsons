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
//   - `string`: path to a local file
//   - `[]string`: paths of local files
//   - `[]byte`: content of a file
//   - `[][]byte`: content list of files
//   - `io.Reader`: content reader
//   - `[]io.Reader`: content readers
func Merge(inputs ...interface{}) ([]byte, error) {
	return stdMerger.Merge(inputs...)
}

// MergeAs loads inputs of the specific format and merges into a single json.
//
// Accepted Input:
//
//   - `string`: path to a local file
//   - `[]string`: paths of local files
//   - `[]byte`: content of a file
//   - `[][]byte`: content list of files
//   - `io.Reader`: content reader
//   - `[]io.Reader`: content readers
func MergeAs(format Format, inputs ...interface{}) ([]byte, error) {
	return stdMerger.MergeAs(format, inputs...)
}
