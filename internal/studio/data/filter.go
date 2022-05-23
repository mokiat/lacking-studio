package data

import (
	"strings"

	"github.com/mokiat/lacking/util/filter"
)

func FilterWithKind(kind ResourceKind) filter.Func[*Resource] {
	return func(resource *Resource) bool {
		return resource.kind == kind
	}
}

func FilterWithSimilarName(name string) filter.Func[*Resource] {
	if name == "" {
		return filter.Always[*Resource]()
	}
	name = strings.ToLower(name)
	return func(resource *Resource) bool {
		return strings.Contains(strings.ToLower(resource.name), name)
	}
}
