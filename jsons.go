package jsons

// Merge merges input JSONs into a single JSON.
//
// The default merger does only simple merging, which means:
//   - Simple values (string, number, boolean) are overwritten by later ones.
//   - Container values (object, array) are merged recursively.
//
// Accepted Input:
//
//   - string: path to a local file
//   - []string: paths of local files
//   - []byte: content of a file
//   - [][]byte: content list of files
//   - io.Reader: content reader
//   - []io.Reader: content readers
//
// If you need complex merging, create a custom merger with options.
func Merge(inputs ...interface{}) ([]byte, error) {
	return NewMerger().Merge(inputs...)
}
