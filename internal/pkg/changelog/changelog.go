package changelog

import (
	"fmt"
)

func NewChangelog() (*ChangeLog, error) {
	builder, err := NewChangeLogBuilder()
	if err != nil {
		return nil, err
	}

	changeLog, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("❌ %s", err)
	}

	return changeLog, nil
}
