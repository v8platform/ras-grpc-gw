package md

import (
	"context"
	"fmt"
	"strings"
)

type RequestMetadata map[string][]string

type AnnotateContextOption func(ctx context.Context) context.Context

type AnnotationHandler func(ctx context.Context, req interface{}) RequestMetadata

// New creates an RequestMetadata from a given key-value map.
//
// Only the following ASCII characters are allowed in keys:
//  - digits: 0-9
//  - uppercase letters: A-Z (normalized to lower)
//  - lowercase letters: a-z
//  - special characters: -_.
// Uppercase letters are automatically converted to lowercase.
//
// Keys beginning with "grpc-" are reserved for grpc-internal use only and may
// result in errors if set in metadata.
func New(m map[string]string) RequestMetadata {
	md := RequestMetadata{}
	for k, val := range m {
		key := strings.ToLower(k)
		md[key] = append(md[key], val)
	}
	return md
}

// Pairs returns an RequestMetadata formed by the mapping of key, value ...
// Pairs panics if len(kv) is odd.
//
// Only the following ASCII characters are allowed in keys:
//  - digits: 0-9
//  - uppercase letters: A-Z (normalized to lower)
//  - lowercase letters: a-z
//  - special characters: -_.
// Uppercase letters are automatically converted to lowercase.
//
// Keys beginning with "grpc-" are reserved for grpc-internal use only and may
// result in errors if set in metadata.
func Pairs(kv ...string) RequestMetadata {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: Pairs got the odd number of input pairs for metadata: %d", len(kv)))
	}
	md := RequestMetadata{}
	for i := 0; i < len(kv); i += 2 {
		key := strings.ToLower(kv[i])
		md[key] = append(md[key], kv[i+1])
	}
	return md
}

func Join(mds ...RequestMetadata) RequestMetadata {
	out := RequestMetadata{}
	for _, md := range mds {
		for k, v := range md {
			out[k] = append(out[k], v...)
		}
	}
	return out
}

func (m RequestMetadata) Clone(copiedKeys ...string) RequestMetadata {
	newMd := Pairs()
	for k, vv := range m {
		found := false
		if len(copiedKeys) == 0 {
			found = true
		} else {
			for _, allowedKey := range copiedKeys {
				if strings.EqualFold(allowedKey, k) {
					found = true
					break
				}
			}
		}
		if !found {
			continue
		}
		newMd[k] = make([]string, len(vv))
		copy(newMd[k], vv)
	}
	return newMd
}

func (m RequestMetadata) ToContext(ctx context.Context) context.Context {
	return NewRequestMetadataContext(ctx, m)
}

func (m RequestMetadata) Range(fn func(key string, val []string) bool) {
	for key, val := range m {
		if !fn(key, val) {
			return
		}
	}
}

func (m RequestMetadata) RangeKey(key string, fn func(index int, val string) bool) {

	values, ok := m[key]
	if !ok {
		return
	}
	for i, val := range values {
		if !fn(i, val) {
			return
		}
	}
}

func (m RequestMetadata) Get(key string) string {
	k := strings.ToLower(key)
	vv, ok := m[k]
	if !ok {
		return ""
	}
	return vv[0]
}

func (m RequestMetadata) Has(key string) bool {
	k := strings.ToLower(key)
	_, ok := m[k]
	return ok
}

func (m RequestMetadata) Del(key string) RequestMetadata {
	k := strings.ToLower(key)
	delete(m, k)
	return m
}

func (m RequestMetadata) Set(key string, value string) RequestMetadata {
	k := strings.ToLower(key)
	m[k] = []string{value}
	return m
}

func (m RequestMetadata) Add(key string, value string) RequestMetadata {
	k := strings.ToLower(key)
	m[k] = append(m[k], value)
	return m
}
