package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewRandomString(t *testing.T) {
	testCases := []struct {
		name   string
		length int
		expect int
	}{
		{
			name:   "0 length input",
			length: 0,
			expect: 6,
		},
		{
			name:   "default length input",
			length: 6,
			expect: 6,
		},
		{
			name:   "11 length input",
			length: 11,
			expect: 11,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewRandomString(tc.length)
			require.Equal(t, tc.expect, len(s))
		})
	}
}
