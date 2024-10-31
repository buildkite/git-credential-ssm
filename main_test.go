package main

import (
	"testing"
)

func TestOutput(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"basic": {
			input:    "https://user1:passw123@example.com",
			expected: "protocol=https\nhost=example.com\nusername=user1\npassword=passw123\n",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			out, err := parseUrl(tt.input)

			if err != nil {
				t.Fatal(err)
			}

			if out.String() != tt.expected {
				t.Errorf("got %s\nwant %s", out.String(), tt.expected)
			}

		})
	}
}
