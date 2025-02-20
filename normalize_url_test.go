package main

import (
	"reflect"
	"testing"
)

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

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputHTML string
		expected  []string
		wantErr   bool
	}{
		{
			name:     "base case with relative and absolute URLs",
			inputURL: "https://blog.deloz.dev",
			inputHTML: `
							<html>
									<body>
											<a href="/about">About</a>
											<a href="https://deloz.dev/test">Test</a>
									</body>
							</html>`,
			expected: []string{
				"https://blog.deloz.dev/about",
				"https://deloz.dev/test",
			},
			wantErr: false,
		},
		{
			name:     "handle empty href",
			inputURL: "https://blog.deloz.dev",
			inputHTML: `
							<html>
									<body>
											<a href="">Empty</a>
									</body>
							</html>`,
			expected: []string{
				"https://blog.deloz.dev",
			},
			wantErr: false,
		},
		{
			name:     "nested elements",
			inputURL: "https://blog.deloz.dev",
			inputHTML: `
							<html>
									<body>
											<div>
													<a href="/one">One</a>
													<div>
															<a href="/two">Two</a>
													</div>
											</div>
									</body>
							</html>`,
			expected: []string{
				"https://blog.deloz.dev/one",
				"https://blog.deloz.dev/two",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getURLsFromHTML(tt.inputHTML, tt.inputURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("getURLsFromHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("getURLsFromHTML() = %v, want %v", got, tt.expected)
			}
		})
	}
}
