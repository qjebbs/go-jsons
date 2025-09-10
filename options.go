// Copyright 2020 Jebbs. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package jsons

// Option is the option for merger
type Option func(m *Merger)

// options is the merge options
type options struct {
	OrderBy       []field
	MergeBy       []field
	TypeOverride  bool
	MarshalPrefix string
	MarshalIndent string
}

// field is the field for rules
type field struct {
	Name   string // field name
	Remove bool   // whether to remove the field after merged
}

// WithOrderBy is the order by field for slice sort rule
func WithOrderBy(name string) Option {
	return func(m *Merger) {
		m.options.OrderBy = append(m.options.OrderBy, field{
			Name:   name,
			Remove: false,
		})
	}
}

// WithMergeBy is the merge by field for slice sort rule
func WithMergeBy(name string) Option {
	return func(m *Merger) {
		m.options.MergeBy = append(m.options.MergeBy, field{
			Name:   name,
			Remove: false,
		})
	}
}

// WithOrderByAndRemove is the order by field for slice merge rule
func WithOrderByAndRemove(name string) Option {
	return func(m *Merger) {
		m.options.OrderBy = append(m.options.OrderBy, field{
			Name:   name,
			Remove: true,
		})
	}
}

// WithMergeByAndRemove is the merge by field for slice merge rule
func WithMergeByAndRemove(name string) Option {
	return func(m *Merger) {
		m.options.MergeBy = append(m.options.MergeBy, field{
			Name:   name,
			Remove: true,
		})
	}
}

// WithTypeOverride sets whether to override the type when merging.
func WithTypeOverride(override bool) Option {
	return func(m *Merger) {
		m.options.TypeOverride = override
	}
}

// WithIndent sets the indent options for merged output.
func WithIndent(prefix, indent string) Option {
	return func(m *Merger) {
		m.options.MarshalPrefix = prefix
		m.options.MarshalIndent = indent
	}
}
