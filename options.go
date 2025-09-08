// Copyright 2020 Jebbs. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package jsons

// Option is the option for merger
type Option func(o *Options)

// Options is the merge rules
type Options struct {
	OrderBy      []OptionField
	MergeBy      []OptionField
	TypeOverride bool
}

// OptionField is the OptionField for rules
type OptionField struct {
	Name   string // field name
	Remove bool   // whether to remove the field after merged
}

// WithOrderBy is the order by field for slice sort rule
func WithOrderBy(name string) Option {
	return func(o *Options) {
		o.OrderBy = append(o.OrderBy, OptionField{
			Name:   name,
			Remove: false,
		})
	}
}

// WithMergeBy is the merge by field for slice sort rule
func WithMergeBy(name string) Option {
	return func(o *Options) {
		o.MergeBy = append(o.MergeBy, OptionField{
			Name:   name,
			Remove: false,
		})
	}
}

// WithOrderByAndRemove is the order by field for slice merge rule
func WithOrderByAndRemove(name string) Option {
	return func(o *Options) {
		o.OrderBy = append(o.OrderBy, OptionField{
			Name:   name,
			Remove: true,
		})
	}
}

// WithMergeByAndRemove is the merge by field for slice merge rule
func WithMergeByAndRemove(name string) Option {
	return func(o *Options) {
		o.MergeBy = append(o.MergeBy, OptionField{
			Name:   name,
			Remove: true,
		})
	}
}

// WithTypeOverride sets whether to override the type when merging.
func WithTypeOverride(override bool) Option {
	return func(o *Options) {
		o.TypeOverride = override
	}
}
