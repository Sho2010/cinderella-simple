package claim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateError(t *testing.T) {
	var tests = []struct {
		name     string
		expected error
		given    Claim
	}{
		{"Expected raise ErrorRequireNamespace", ErrorRequireNamespace, &ClaimBase{Namespaces: []string{}, Subject: "subject"}},
		{"Expected raise ErrorRequireSubject", ErrorRequireSubject, &ClaimBase{Namespaces: []string{}, Subject: ""}},
		{"Expected raise Error RFC1123 format",
			&ClaimValidationError{field: "Namespaces", errorType: "RFC1123"},
			&ClaimBase{Namespaces: []string{"invalid@namespace"}, Subject: "subject"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.given.Validate()
			assert.ErrorIs(t, err, tt.expected)
		})
	}
}

func TestNormalizeDNS1123(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
		given    string
	}{
		{"scott-tiger unnecessary normalize", "scott-tiger", "scott-tiger"},
		{"replace normalize target", "scott-tiger", "scott.tiger"},
		{"replace normalize target", "scott-tiger", "scott@tiger"},
		{"replace normalize target", "scott-tiger", "scott/tiger"},
		{"replace normalize target", "scott-tiger", "scott_tiger"},
		{"replace normalize target", "scotttiger", "ScottTiger"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NormalizeDNS1123(tt.given)
			if actual != tt.expected {
				t.Logf("%v", err)
				t.Errorf("given(%s): expected [%s], actual [%s]", tt.given, tt.expected, actual)
			}
		})
	}
}
