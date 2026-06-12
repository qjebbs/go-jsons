package preprocess

// PreprocessorFunc is the type of preprocessor function, which can be used to preprocess values before merging.
type PreprocessorFunc func(key string, value interface{}) interface{}
