package lion

import "strings"

type Tag struct {
	Name string
	Opts map[string]string
}

// Tag
// returns the tag information of the field.
func (field *Field) Tag(name string) *Tag {
	if field.ref == nil {
		if field.tags == nil {
			field.tags = map[string]*Tag{}
		}
		v, ok := field.tags[name]
		if ok {
			return v
		}
		tag := parseTag(field.field.Tag.Get(name))
		if tag.Name == "" {
			tag.Name = field.field.Name
		}
		field.tags[name] = tag
		return tag
	}
	return field.ref.Tag(name)
}

func parseTag(v string) *Tag {
	tag := &Tag{}
	parts := strings.Split(v, ",")
	tag.Name = strings.TrimSpace(parts[0])
	for i := 1; i < len(parts); i++ {
		if tag.Opts == nil {
			tag.Opts = map[string]string{}
		}
		v := strings.TrimSpace(parts[i])
		idx := strings.IndexByte(v, '=')
		if idx > -1 {
			tag.Opts[strings.TrimSpace(v[:idx])] = strings.TrimSpace(v[idx+1:])
		} else {
			tag.Opts[v] = "true"
		}
	}
	return tag
}
