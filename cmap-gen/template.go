package main

import "text/template"

var tmplSource = `
// {{.MapTypeName}} is sharded concurrent map which key type is string and
// value type is {{.TypeName}}
type {{.MapTypeName}} struct {
	shards []*{{.ShardTypeName}}
}

// {{.NewMethodName}} creates new {{.MapTypeName}} with specified shards count
func {{.NewMethodName}}(nShards int) *{{.MapTypeName}} {
	shards := make([]*{{.ShardTypeName}}, nShards)
	for i := 0; i < nShards; i++ {
		shards[i] = New{{.ShardTypeName}}()
	}
	return &{{.MapTypeName}}{shards: shards}
}

func (c {{.MapTypeName}}) hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32() % uint32(len(c.shards))
}

// Get returns the value stored by specified key
func (c {{.MapTypeName}}) Get(key {{.KeyTypeName}}) {{.TypeName}} {
	return c.shards[c.hash(key)].Get(key)
}

// Set stores the specified value under the specified key
func (c {{.MapTypeName}}) Set(key {{.KeyTypeName}}, value {{.TypeName}}) {
	c.shards[c.hash(key)].Set(key, value)
}

// {{.ShardTypeName}} is concurrent map which key type is string and
// value type is {{.TypeName}}
type {{.ShardTypeName}} struct {
	mu   sync.RWMutex
	data map[string]{{.TypeName}}
}

// New{{.ShardTypeName}} creates new {{.ShardTypeName}}
func New{{.ShardTypeName}}() *{{.ShardTypeName}} {
	return &{{.ShardTypeName}}{
		data: make(map[{{.KeyTypeName}}]{{.TypeName}}),
	}
}

// Get returns the value stored by specified key
func (c *{{.ShardTypeName}}) Get(key {{.KeyTypeName}}) {{.TypeName}} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

// Set stores the specified value under the specified key
func (c *{{.ShardTypeName}}) Set(key {{.KeyTypeName}}, value {{.TypeName}}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}
`

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("name").Parse(tmplSource))
}
