// Copyright 2020 Jebbs. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package jsons

// Option is the option for merger
type Option func(o *Options)

// OptionField is the OptionField for rules
type OptionField struct {
	Key    string // field key
	Remove bool   // whether to remove the field after merged
}

// WithOrderBy is the order by field for slice sort rule
func WithOrderBy(key string) Option {
	return func(o *Options) {
		o.OrderBy = append(o.OrderBy, OptionField{
			Key:    key,
			Remove: false,
		})
	}
}

// WithMergeBy is the merge by field for slice sort rule
func WithMergeBy(key string) Option {
	return func(o *Options) {
		o.MergeBy = append(o.MergeBy, OptionField{
			Key:    key,
			Remove: false,
		})
	}
}

// WithOrderByAndRemove is the order by field for slice merge rule
func WithOrderByAndRemove(key string) Option {
	return func(o *Options) {
		o.OrderBy = append(o.OrderBy, OptionField{
			Key:    key,
			Remove: true,
		})
	}
}

// WithMergeByAndRemove is the merge by field for slice merge rule
func WithMergeByAndRemove(key string) Option {
	return func(o *Options) {
		o.MergeBy = append(o.MergeBy, OptionField{
			Key:    key,
			Remove: true,
		})
	}
}
