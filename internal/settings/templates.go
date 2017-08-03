package settings

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/kyokomi/emoji"
)

type TemplateSettings map[string]string

func (s TemplateSettings) Dump(indent string) string {
	maxlen := 0
	ids := []string{}
	for id := range s {
		if maxlen < len(id) {
			maxlen = len(id)
		}
		ids = append(ids, id)
	}
	sort.Strings(ids)
	result := ""
	for _, id := range ids {
		tmpl := s[id]
		str := fmt.Sprintf("%s%-"+strconv.Itoa(maxlen)+"s = %s\n", indent, id, tmpl)
		result += emoji.Sprint(str)
	}
	return result
}
