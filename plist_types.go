package plist

import (
	"hash/crc32"
	"sort"
	"time"
	"strconv"
)

// magic value used in the non-binary encoding of UIDs
// (stored as a dictionary mapping CF$UID->integer)
const cfUIDMagic = "CF$UID"

type cfValue interface {
	typeName() string
	hash() interface{}
}

type cfDictionary struct {
	keys   sort.StringSlice
	values []cfValue
}

func (*cfDictionary) typeName() string {
	return "dictionary"
}

func (p *cfDictionary) hash() interface{} {
	return p
}

func (p *cfDictionary) Len() int {
	return len(p.keys)
}

func (p *cfDictionary) Less(i, j int) bool {
	return p.keys.Less(i, j)
}

func (p *cfDictionary) Swap(i, j int) {
	p.keys.Swap(i, j)
	p.values[i], p.values[j] = p.values[j], p.values[i]
}

func (p *cfDictionary) sort() {
	sort.Sort(p)
}

func (c *cfDictionary) filterNodes(nodes []Node) (keys []string, values []cfValue, filtered []Node) {
	m := c.toMap()
	add := map[string]bool{}
	for _, node := range nodes {
		key := node.Value()
		if _, ok := node.(*Annotation); ok {
			keys = append(keys, key)
			values = append(values, nil)
			add[key] = true

		}
		if _, ok := node.(*MetaNode); ok {
			if v, ok := m[key]; ok {
				keys = append(keys, key)
				values = append(values, v)
				add[key] = true
			}
		}
		if add[key] {
			filtered = append(filtered, node)
		}
	}

	// Make sure keys and values without corresponding node get included.
	for k, v := range m {
		if _, ok := add[k]; !ok {
			keys = append(keys, k)
			values = append(values, v)
		}
	}
	return keys, values, filtered
}

func (c *cfDictionary) toMap() map[string]cfValue {
	m := map[string]cfValue{}
	for i, k := range c.keys {
		m[k] = c.values[i]
	}
	return m
}

func (p *cfDictionary) maybeUID(lax bool) cfValue {
	if len(p.keys) == 1 && p.keys[0] == "CF$UID" && len(p.values) == 1 {
		pval := p.values[0]
		if integer, ok := pval.(*cfNumber); ok {
			return cfUID(integer.value)
		}
		// Openstep only has cfString. Act like the unmarshaller a bit.
		if lax {
			if str, ok := pval.(cfString); ok {
				if i, err := strconv.ParseUint(string(str), 10, 64); err == nil {
					return cfUID(i)
				}
			}
		}
	}
	return p
}

type cfArray struct {
	values []cfValue
}

func (*cfArray) typeName() string {
	return "array"
}

func (p *cfArray) hash() interface{} {
	return p
}

type cfString string

func (cfString) typeName() string {
	return "string"
}

func (p cfString) hash() interface{} {
	return string(p)
}

type cfNumber struct {
	signed bool
	value  uint64
}

func (*cfNumber) typeName() string {
	return "integer"
}

func (p *cfNumber) hash() interface{} {
	if p.signed {
		return int64(p.value)
	}
	return p.value
}

type cfReal struct {
	wide  bool
	value float64
}

func (cfReal) typeName() string {
	return "real"
}

func (p *cfReal) hash() interface{} {
	if p.wide {
		return p.value
	}
	return float32(p.value)
}

type cfBoolean bool

func (cfBoolean) typeName() string {
	return "boolean"
}

func (p cfBoolean) hash() interface{} {
	return bool(p)
}

type cfUID UID

func (cfUID) typeName() string {
	return "UID"
}

func (p cfUID) hash() interface{} {
	return p
}

func (p cfUID) toDict() *cfDictionary {
	return &cfDictionary{
		keys: []string{cfUIDMagic},
		values: []cfValue{&cfNumber{
			signed: false,
			value:  uint64(p),
		}},
	}
}

type cfData []byte

func (cfData) typeName() string {
	return "data"
}

func (p cfData) hash() interface{} {
	// Data are uniqued by their checksums.
	// Todo: Look at calculating this only once and storing it somewhere;
	// crc32 is fairly quick, however.
	return crc32.ChecksumIEEE([]byte(p))
}

type cfDate time.Time

func (cfDate) typeName() string {
	return "date"
}

func (p cfDate) hash() interface{} {
	return time.Time(p)
}

type cfAnnotation struct {
	value string
}

func (cfAnnotation) typeName() string {
	return "annotation"
}

func (p cfAnnotation) hash() interface{} {
	return p.value
}
