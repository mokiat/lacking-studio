package observer

import "fmt"

type Change interface {
	Description() string
}

func NewChange(description string) Change {
	return &stringChange{
		description: description,
	}
}

type stringChange struct {
	description string
}

func (c *stringChange) Description() string {
	return c.description
}

func ExtChange(parent Change, description string) Change {
	return &extendedChange{
		parent:      parent,
		description: description,
	}
}

type extendedChange struct {
	parent      Change
	description string
}

func (c *extendedChange) Description() string {
	return fmt.Sprintf("%s: %s", c.parent.Description(), c.description)
}

func (c *extendedChange) Is(target Change) bool {
	return (c == target) || IsChange(c.parent, target)
}

// func ExtendChange(parent, current Change) Change {
// 	return &extendedChange{
// 		parent:  parent,
// 		current: current,
// 	}
// }

// type extendedChange struct {
// 	parent  Change
// 	current Change
// }

// func (c *extendedChange) Description() string {
// 	return fmt.Sprintf("%s: %s", c.parent.Description(), c.current.Description())
// }

// func (c *extendedChange) Is(target Change) bool {
// 	if c == target {
// 		return true
// 	}
// 	return IsChange(c.current, target) || IsChange(c.parent, target)
// }

type MultiChange struct {
	Changes []Change
}

func (c MultiChange) Description() string {
	return fmt.Sprintf("multi-change (%d)", len(c.Changes))
}

func (c MultiChange) Is(target Change) bool {
	for _, candidate := range c.Changes {
		if IsChange(candidate, target) {
			return true
		}
	}
	return false
}

func IsChange(change, target Change) bool {
	if change == target {
		return true
	}
	comparable, ok := change.(interface {
		Is(Change) bool
	})
	if ok && comparable.Is(target) {
		return true
	}

	parentable, ok := change.(interface {
		Parent() Change
	})
	if !ok {
		return false
	}
	parent := parentable.Parent()
	if parent == nil {
		return false
	}
	return IsChange(parent, target)
}
