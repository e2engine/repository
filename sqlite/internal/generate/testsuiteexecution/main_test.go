package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestIsTargetType(t *testing.T) {
	testCases := []struct {
		name     string
		typeName string
		expected bool
	}{
		{
			name:     "list row",
			typeName: "ListTestSuiteExecutionsByIDFirstAscRow",
			expected: true,
		},
		{
			name:     "single row",
			typeName: "GetTestSuiteExecutionRow",
			expected: false,
		},
		{
			name:     "list params",
			typeName: "ListTestSuiteExecutionsByIDFirstAscParams",
			expected: false,
		},
		{
			name:     "unrelated row",
			typeName: "ListTestExecutionsByIDFirstAscRow",
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := isTargetType(testCase.typeName)

			if actual != testCase.expected {
				t.Fatalf(
					"Expected %v for %q, got %v",
					testCase.expected,
					testCase.typeName,
					actual,
				)
			}
		})
	}
}

func TestValidateFields(t *testing.T) {
	structType := &ast.StructType{
		Fields: &ast.FieldList{},
	}

	for _, field := range requiredFields {
		structType.Fields.List = append(
			structType.Fields.List,
			&ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent(field),
				},
			},
		)
	}

	if err := validateFields("ValidRow", structType); err != nil {
		t.Fatalf("Unexpected validation error: %v", err)
	}

	structType.Fields.List = structType.Fields.List[:len(structType.Fields.List)-1]

	err := validateFields("InvalidRow", structType)
	if err == nil {
		t.Fatalf("Expected validation error")
	}

	if !strings.Contains(err.Error(), "Summary") {
		t.Fatalf("Expected missing Summary error, got %v", err)
	}
}

func TestGenerate(t *testing.T) {
	source, err := generate([]string{
		"ListTestSuiteExecutionsByIDFirstAscRow",
	})
	if err != nil {
		t.Fatalf("Unexpected generate error: %v", err)
	}

	generated := string(source)

	expected := []string{
		"func normalizeListTestSuiteExecutionsByIDFirstAsc(",
		"values []db.ListTestSuiteExecutionsByIDFirstAscRow",
		"ID:              values[i].ID",
		"StartedAt:       values[i].StartedAt",
		"FinishedAt:      values[i].FinishedAt",
		"Status:          values[i].Status",
		"EnvironmentID:   values[i].EnvironmentID",
		"EnvironmentName: values[i].EnvironmentName",
		"TestSuiteID:     values[i].TestSuiteID",
		"TestSuiteName:   values[i].TestSuiteName",
		"TestsCount:      values[i].TestsCount",
		"Summary:         values[i].Summary",
	}

	for _, value := range expected {
		if !strings.Contains(generated, value) {
			t.Fatalf(
				"Expected generated source to contain %q:\n%s",
				value,
				generated,
			)
		}
	}

	if _, err := parser.ParseFile(
		token.NewFileSet(),
		"generated.go",
		source,
		0,
	); err != nil {
		t.Fatalf("Generated invalid Go source: %v", err)
	}
}
