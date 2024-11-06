package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	tests := map[string]struct {
		input       string
		expected    string
		expectederr bool
	}{
		"basic": {
			input:    "https://user1:passw123@example.com",
			expected: "protocol=https\nhost=example.com\nusername=user1\npassword=passw123\n",
		},
		"trailing-whitespace": {
			input:    "https://user1:passw123@example.com\n",
			expected: "protocol=https\nhost=example.com\nusername=user1\npassword=passw123\n",
		},
		"no-username-password": {
			input:       "not-a-url",
			expectederr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			out, err := parseUrl(tt.input)

			if tt.expectederr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if !tt.expectederr {
				assert.Equal(t, tt.expected, out.String())
			}
		})
	}
}
