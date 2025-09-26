package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 15",
			size: 15,
		}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stringOne := NewRandomString(test.size)
			stringTwo := NewRandomString(test.size)

			assert.Len(t, stringOne, test.size)
			assert.Len(t, stringTwo, test.size)
			assert.NotEqual(t, stringOne, stringTwo)

		})
	}

}
