package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	name string
	size int
}{
	{
		name: "size = 1",
		size: 1,
	},
	{
		name: "size = 10",
		size: 10,
	},
	{
		name: "size = 30",
		size: 30,
	},
}

func TestRandom(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			str1 := GenerateRandomString(tc.size)
			str2 := GenerateRandomString(tc.size)
			assert.Len(t, str1, tc.size)
			assert.Len(t, str2, tc.size)
			assert.NotEqual(t, str1, str2)
		})
	}
}
