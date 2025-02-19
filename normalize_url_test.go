package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
		expectedError bool
	}{
		{
			name:          "remove scheme",
			inputURL:      "https://test.deloz.dev/path",
			expected:      "test.deloz.dev/path",
			expectedError: false,
		},
		{
			name:          "empty url",
			inputURL:      "",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "no subdomain",
			inputURL:      "https://deloz.dev/path",
			expected:      "deloz.dev/path",
			expectedError: false,
		},
		{
			name:          "longer path",
			inputURL:      "https://deloz.dev/path/to/path/",
			expected:      "deloz.dev/path/to/path",
			expectedError: false,
		},
		{
			name:          "broken link",
			inputURL:      "htttps:///thisisabrokenlink.com//path///",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "no path",
			inputURL:      "https://deloz.dev",
			expected:      "deloz.dev",
			expectedError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if tc.expectedError && err == nil {
				t.Errorf("expected error but got none")
				return
			}
			if !tc.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !tc.expectedError && actual != tc.expected {
				t.Errorf("expected URL: %v, actual: %v", tc.expected, actual)
			}
		})
	}
}
