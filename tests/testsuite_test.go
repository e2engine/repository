package tests

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	corerepository "github.com/e2engine/core/repository"
)

func TestTestSuiteCRUD(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuite

	testSuite1 := corerepository.TestSuite(getResource(t, "test_suite", 1))
	created1, err := repository.Create(t.Context(), &testSuite1)
	if err != nil {
		t.Fatalf("Unexpected test suite1 create error: %v", err)
	}
	if created1 == nil {
		t.Fatalf("Expected created1 test suite to be non-nil")
	}
	assertEqualResources(t, corerepository.Resource(testSuite1), corerepository.Resource(*created1))

	testSuite2 := corerepository.TestSuite(getResource(t, "test_suite", 2))
	created2, err := repository.Create(t.Context(), &testSuite2)
	if err != nil {
		t.Fatalf("Unexpected test suite2 create error: %v", err)
	}
	if created2 == nil {
		t.Fatalf("Expected created2 test suite to be non-nil")
	}
	assertEqualResources(t, corerepository.Resource(testSuite2), corerepository.Resource(*created2))

	got1, err := repository.Get(t.Context(), testSuite1.ID)
	if err != nil {
		t.Fatalf("Unexpected test suite1 get error: %v", err)
	}
	if got1 == nil {
		t.Fatalf("Expected test suite1 to be non-nil")
	}
	assertEqualResources(t, corerepository.Resource(testSuite1), corerepository.Resource(*got1))

	got2, err := repository.Get(t.Context(), testSuite2.ID)
	if err != nil {
		t.Fatalf("Unexpected test suite2 get error: %v", err)
	}
	if got2 == nil {
		t.Fatalf("Expected test suite2 to be non-nil")
	}
	assertEqualResources(t, corerepository.Resource(testSuite2), corerepository.Resource(*got2))

	ids, err := repository.GetIDsByID(t.Context(), idPrefixEnvironment)
	if err != nil {
		t.Fatalf("Unexpected test suite IDs by ID error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("Expected 2 test suite IDs by ID, got %d", len(ids))
	}
	if ids[0] != testSuite1.ID || ids[1] != testSuite2.ID {
		t.Fatalf("Expected test suite IDs %v, got %v", []string{testSuite1.ID, testSuite2.ID}, ids)
	}

	ids, err = repository.GetIDsByName(t.Context(), "Resource test_suite")
	if err != nil {
		t.Fatalf("Unexpected test suite IDs by name error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("Expected 2 test suite IDs by name, got %d", len(ids))
	}
	if ids[0] != testSuite1.ID || ids[1] != testSuite2.ID {
		t.Fatalf("Expected test suite IDs %v, got %v", []string{testSuite1.ID, testSuite2.ID}, ids)
	}

	deleted1, err := repository.Delete(t.Context(), testSuite1.ID)
	if err != nil {
		t.Fatalf("Unexpected test suite1 delete error: %v", err)
	}
	if deleted1 == nil {
		t.Fatalf("Expected deleted1 test suite to be non-nil")
	}
	assertEqualResources(t, corerepository.Resource(testSuite1), corerepository.Resource(*deleted1))

	deleted2, err := repository.Delete(t.Context(), testSuite2.ID)
	if err != nil {
		t.Fatalf("Unexpected test suite2 delete error: %v", err)
	}
	if deleted2 == nil {
		t.Fatalf("Expected deleted2 test suite to be non-nil")
	}
	assertEqualResources(t, corerepository.Resource(testSuite2), corerepository.Resource(*deleted2))
}

func TestTestSuiteCreateErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuite

	testCases := []struct {
		name          string
		mutate        func(*corerepository.TestSuite)
		expectedError error
	}{
		{
			name: "empty ID",
			mutate: func(testSuite *corerepository.TestSuite) {
				testSuite.ID = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "empty version",
			mutate: func(testSuite *corerepository.TestSuite) {
				testSuite.Version = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "empty name",
			mutate: func(testSuite *corerepository.TestSuite) {
				testSuite.Name = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "empty spec",
			mutate: func(testSuite *corerepository.TestSuite) {
				testSuite.Spec = []byte{}
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "zero created at",
			mutate: func(testSuite *corerepository.TestSuite) {
				testSuite.CreatedAt = time.Time{}
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "zero updated at",
			mutate: func(testSuite *corerepository.TestSuite) {
				testSuite.UpdatedAt = time.Time{}
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
	}

	for i, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testSuite := corerepository.TestSuite(getResource(t, "test_suite", i+1))
			testCase.mutate(&testSuite)

			created, err := repository.Create(t.Context(), &testSuite)
			if err == nil {
				t.Fatalf("Expected test suite create error")
			}
			if created != nil {
				t.Fatalf("Expected created test suite to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestTestSuiteCreateEmptyDescription(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuite

	testsuite := corerepository.TestSuite(
		getResource(t, "test_suite", 1),
	)
	testsuite.Description = ""

	created, err := repository.Create(t.Context(), &testsuite)
	if err != nil {
		t.Fatalf("Unexpected testsuite create error: %v", err)
	}
	if created == nil {
		t.Fatalf("Expected created testsuite to be non-nil")
	}

	assertEqualResources(
		t,
		corerepository.Resource(testsuite),
		corerepository.Resource(*created),
	)
}

func TestTestSuiteCreateDuplicatedName(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuite

	testSuite1 := corerepository.TestSuite(getResource(t, "test_suite", 1))
	created1, err := repository.Create(t.Context(), &testSuite1)
	if err != nil {
		t.Fatalf("Unexpected test suite1 create error: %v", err)
	}
	if created1 == nil {
		t.Fatalf("Expected created1 test suite to be non-nil")
	}

	testSuite2 := corerepository.TestSuite(getResource(t, "test_suite", 2))
	testSuite2.Name = testSuite1.Name
	created2, err := repository.Create(t.Context(), &testSuite2)
	if err == nil {
		t.Fatalf("Expected test suite2 create error")
	}
	if created2 != nil {
		t.Fatalf("Expected created2 test suite to be nil")
	}

	if !errors.Is(err, corerepository.ErrAlreadyExists) {
		t.Fatalf(
			"Expected error %v, got %v",
			corerepository.ErrAlreadyExists,
			err,
		)
	}
}

func TestTestSuiteCreateIgnoresPayloadKind(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuite

	testSuite := corerepository.TestSuite(getResource(t, "invalid", 1))
	created, err := repository.Create(t.Context(), &testSuite)
	if err != nil {
		t.Fatalf("Unexpected test suite create error: %v", err)
	}
	if created == nil {
		t.Fatalf("Expected created test suite to be non-nil")
	}
	if created.Kind != "test_suite" {
		t.Fatalf("Expected created test suite kind to be test_suite, got %q", created.Kind)
	}
}

func TestTestSuiteGetErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuite

	testCases := []struct {
		name          string
		id            string
		expectedError error
	}{
		{
			name:          "empty ID",
			id:            "",
			expectedError: corerepository.ErrNotFound,
		},
		{
			name:          "non-existing ID",
			id:            "non-existing-id",
			expectedError: corerepository.ErrNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testSuite, err := repository.Get(t.Context(), testCase.id)
			if err == nil {
				t.Fatalf("Expected test suite get error")
			}
			if testSuite != nil {
				t.Fatalf("Expected test suite to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestTestSuiteDeleteErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuite

	testCases := []struct {
		name          string
		id            string
		expectedError error
	}{
		{
			name:          "empty ID",
			id:            "",
			expectedError: corerepository.ErrNotFound,
		},
		{
			name:          "non-existing ID",
			id:            "non-existing-id",
			expectedError: corerepository.ErrNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testSuite, err := repository.Delete(t.Context(), testCase.id)
			if err == nil {
				t.Fatalf("Expected test suite delete error")
			}
			if testSuite != nil {
				t.Fatalf("Expected test suite to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestTestSuiteGetIDsByID(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuite
	resources := getResourcesForList(t, "test_suite")
	createTestSuiteResources(t, repository, resources)

	allIDs := make([]string, len(resources))
	for i := range resources {
		allIDs[i] = resources[i].ID
	}
	slices.Sort(allIDs)

	testCases := []struct {
		name          string
		id            string
		context       func() context.Context
		expectedIDs   []string
		expectedError error
	}{
		{
			name: "empty ID",
			id:   "",
			context: func() context.Context {
				return t.Context()
			},
			expectedIDs: allIDs,
		},
		{
			name: "query error",
			id:   resources[0].ID,
			context: func() context.Context {
				ctx, cancel := context.WithCancel(t.Context())
				cancel()

				return ctx
			},
			expectedError: context.Canceled,
		},
		{
			name: "non-existing ID",
			id:   "does-not-exist",
			context: func() context.Context {
				return t.Context()
			},
			expectedIDs: []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ids, err := repository.GetIDsByID(testCase.context(), testCase.id)
			assertTestSuiteIDs(t, ids, err, testCase.expectedIDs, testCase.expectedError)
		})
	}
}

func TestTestSuiteGetIDsByName(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuite
	resources := getResourcesForList(t, "test_suite")
	createTestSuiteResources(t, repository, resources)

	allIDs := make([]string, len(resources))
	for i := range resources {
		allIDs[i] = resources[i].ID
	}
	slices.Sort(allIDs)

	testCases := []struct {
		name          string
		value         string
		context       func() context.Context
		expectedIDs   []string
		expectedError error
	}{
		{
			name:  "empty name",
			value: "",
			context: func() context.Context {
				return t.Context()
			},
			expectedIDs: allIDs,
		},
		{
			name:  "query error",
			value: resources[0].Name,
			context: func() context.Context {
				ctx, cancel := context.WithCancel(t.Context())
				cancel()

				return ctx
			},
			expectedError: context.Canceled,
		},
		{
			name:  "non-existing name",
			value: "does-not-exist",
			context: func() context.Context {
				return t.Context()
			},
			expectedIDs: []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ids, err := repository.GetIDsByName(testCase.context(), testCase.value)
			assertTestSuiteIDs(t, ids, err, testCase.expectedIDs, testCase.expectedError)
		})
	}
}

func assertTestSuiteIDs(
	t *testing.T,
	ids []string,
	err error,
	expectedIDs []string,
	expectedError error,
) {
	t.Helper()

	if expectedError != nil {
		if err == nil || !errors.Is(err, expectedError) {
			t.Fatalf("Expected error %v, got %v", expectedError, err)
		}
		if ids != nil {
			t.Fatalf("Expected IDs to be nil, got %v", ids)
		}

		return
	}

	if err != nil {
		t.Fatalf("Unexpected IDs lookup error: %v", err)
	}

	slices.Sort(ids)

	if !slices.Equal(ids, expectedIDs) {
		t.Fatalf("Expected IDs %v, got %v", expectedIDs, ids)
	}
}
