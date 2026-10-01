package tests

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	corerepository "github.com/e2engine/core/repository"
)

func TestTestCRUD(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Test

	test1 := corerepository.Test(getResource(t, "test", 1))
	test1.Spec = []byte(`{"tags":["tag1","shared"]}`)
	created1, err := repository.Create(t.Context(), &test1)
	if err != nil {
		t.Fatalf("Unexpected test1 create error: %v", err)
	}
	if created1 == nil {
		t.Fatalf("Expected created1 test to be non-nil")
	}
	assertEqualResources(t, corerepository.Resource(test1), corerepository.Resource(*created1))

	test2 := corerepository.Test(getResource(t, "test", 2))
	test2.Spec = []byte(`{"tags":["tag2","shared"]}`)
	created2, err := repository.Create(t.Context(), &test2)
	if err != nil {
		t.Fatalf("Unexpected test2 create error: %v", err)
	}
	if created2 == nil {
		t.Fatalf("Expected created2 test to be non-nil")
	}
	assertEqualResources(t, corerepository.Resource(test2), corerepository.Resource(*created2))

	got1, err := repository.Get(t.Context(), test1.ID)
	if err != nil {
		t.Fatalf("Unexpected test1 get error: %v", err)
	}
	if got1 == nil {
		t.Fatalf("Expected test1 to be non-nil")
	}
	assertEqualResources(t, corerepository.Resource(test1), corerepository.Resource(*got1))

	got2, err := repository.Get(t.Context(), test2.ID)
	if err != nil {
		t.Fatalf("Unexpected test2 get error: %v", err)
	}
	if got2 == nil {
		t.Fatalf("Expected test2 to be non-nil")
	}
	assertEqualResources(t, corerepository.Resource(test2), corerepository.Resource(*got2))

	ids, err := repository.GetIDsByID(t.Context(), idPrefixEnvironment)
	if err != nil {
		t.Fatalf("Unexpected test IDs by ID error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("Expected 2 test IDs by ID, got %d", len(ids))
	}
	if ids[0] != test1.ID || ids[1] != test2.ID {
		t.Fatalf("Expected test IDs %v, got %v", []string{test1.ID, test2.ID}, ids)
	}

	ids, err = repository.GetIDsByName(t.Context(), "Resource test")
	if err != nil {
		t.Fatalf("Unexpected test IDs by name error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("Expected 2 test IDs by name, got %d", len(ids))
	}
	if ids[0] != test1.ID || ids[1] != test2.ID {
		t.Fatalf("Expected test IDs %v, got %v", []string{test1.ID, test2.ID}, ids)
	}

	testsByTag, err := repository.GetByTag(t.Context(), "shared")
	if err != nil {
		t.Fatalf("Unexpected tests by tag error: %v", err)
	}
	if len(testsByTag) != 2 {
		t.Fatalf("Expected 2 tests by tag, got %d", len(testsByTag))
	}
	assertEqualResources(
		t,
		corerepository.Resource(test1),
		corerepository.Resource(testsByTag[0]),
	)
	assertEqualResources(
		t,
		corerepository.Resource(test2),
		corerepository.Resource(testsByTag[1]),
	)

	deleted1, err := repository.Delete(t.Context(), test1.ID)
	if err != nil {
		t.Fatalf("Unexpected test1 delete error: %v", err)
	}
	if deleted1 == nil {
		t.Fatalf("Expected deleted1 test to be non-nil")
	}
	assertEqualResources(t, corerepository.Resource(test1), corerepository.Resource(*deleted1))

	deleted2, err := repository.Delete(t.Context(), test2.ID)
	if err != nil {
		t.Fatalf("Unexpected test2 delete error: %v", err)
	}
	if deleted2 == nil {
		t.Fatalf("Expected deleted2 test to be non-nil")
	}
	assertEqualResources(t, corerepository.Resource(test2), corerepository.Resource(*deleted2))
}

func TestTestCreateErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Test

	testCases := []struct {
		name          string
		mutate        func(*corerepository.Test)
		expectedError error
	}{
		{
			name: "empty ID",
			mutate: func(test *corerepository.Test) {
				test.ID = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "empty version",
			mutate: func(test *corerepository.Test) {
				test.Version = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "empty name",
			mutate: func(test *corerepository.Test) {
				test.Name = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "empty spec",
			mutate: func(test *corerepository.Test) {
				test.Spec = []byte{}
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "zero created at",
			mutate: func(test *corerepository.Test) {
				test.CreatedAt = time.Time{}
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "zero updated at",
			mutate: func(test *corerepository.Test) {
				test.UpdatedAt = time.Time{}
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
	}

	for i, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			test := corerepository.Test(getResource(t, "test", i+1))
			testCase.mutate(&test)

			created, err := repository.Create(t.Context(), &test)
			if err == nil {
				t.Fatalf("Expected test create error")
			}
			if created != nil {
				t.Fatalf("Expected created test to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestTestCreateEmptyDescription(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Test

	test := corerepository.Test(
		getResource(t, "test", 1),
	)
	test.Description = ""

	created, err := repository.Create(t.Context(), &test)
	if err != nil {
		t.Fatalf("Unexpected test create error: %v", err)
	}
	if created == nil {
		t.Fatalf("Expected created test to be non-nil")
	}

	assertEqualResources(
		t,
		corerepository.Resource(test),
		corerepository.Resource(*created),
	)
}

func TestTestCreateDuplicatedName(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Test

	test1 := corerepository.Test(getResource(t, "test", 1))
	created1, err := repository.Create(t.Context(), &test1)
	if err != nil {
		t.Fatalf("Unexpected test1 create error: %v", err)
	}
	if created1 == nil {
		t.Fatalf("Expected created1 test to be non-nil")
	}

	test2 := corerepository.Test(getResource(t, "test", 2))
	test2.Name = test1.Name
	created2, err := repository.Create(t.Context(), &test2)
	if err == nil {
		t.Fatalf("Expected test2 create error")
	}
	if created2 != nil {
		t.Fatalf("Expected created2 test to be nil")
	}
	if !errors.Is(err, corerepository.ErrAlreadyExists) {
		t.Fatalf(
			"Expected error %v, got %v",
			corerepository.ErrAlreadyExists,
			err,
		)
	}
}

func TestTestCreateIgnoresPayloadKind(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Test

	test := corerepository.Test(getResource(t, "invalid", 1))
	created, err := repository.Create(t.Context(), &test)
	if err != nil {
		t.Fatalf("Unexpected test create error: %v", err)
	}
	if created == nil {
		t.Fatalf("Expected created test to be non-nil")
	}
	if created.Kind != "test" {
		t.Fatalf("Expected created test kind to be test, got %q", created.Kind)
	}
}

func TestTestGetErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Test

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
			test, err := repository.Get(t.Context(), testCase.id)
			if err == nil {
				t.Fatalf("Expected test get error")
			}
			if test != nil {
				t.Fatalf("Expected test to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestTestDeleteErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Test

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
			test, err := repository.Delete(t.Context(), testCase.id)
			if err == nil {
				t.Fatalf("Expected test delete error")
			}
			if test != nil {
				t.Fatalf("Expected test to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestTestGetIDsByID(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Test
	resources := getResourcesForList(t, "test")
	createTestResources(t, repository, resources)

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
			ids, err := repository.GetIDsByID(
				testCase.context(),
				testCase.id,
			)

			if testCase.expectedError != nil {
				if err == nil || !errors.Is(err, testCase.expectedError) {
					t.Fatalf(
						"Expected error %v, got %v",
						testCase.expectedError,
						err,
					)
				}
				if ids != nil {
					t.Fatalf("Expected IDs to be nil, got %v", ids)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected GetIDsByID error: %v", err)
			}

			slices.Sort(ids)

			if !slices.Equal(ids, testCase.expectedIDs) {
				t.Fatalf(
					"Expected IDs %v, got %v",
					testCase.expectedIDs,
					ids,
				)
			}
		})
	}
}

func TestTestGetIDsByName(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Test
	resources := getResourcesForList(t, "test")
	createTestResources(t, repository, resources)

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
			ids, err := repository.GetIDsByName(
				testCase.context(),
				testCase.value,
			)

			if testCase.expectedError != nil {
				if err == nil || !errors.Is(err, testCase.expectedError) {
					t.Fatalf(
						"Expected error %v, got %v",
						testCase.expectedError,
						err,
					)
				}
				if ids != nil {
					t.Fatalf("Expected IDs to be nil, got %v", ids)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected GetIDsByName error: %v", err)
			}

			slices.Sort(ids)

			if !slices.Equal(ids, testCase.expectedIDs) {
				t.Fatalf(
					"Expected IDs %v, got %v",
					testCase.expectedIDs,
					ids,
				)
			}
		})
	}
}

func TestTestGetByTag(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Test

	test1 := corerepository.Test(getResource(t, "test", 1))
	test1.Spec = []byte(`{"tags":["tag1","shared"]}`)

	test2 := corerepository.Test(getResource(t, "test", 2))
	test2.Spec = []byte(`{"tags":["tag2","shared"]}`)

	test3 := corerepository.Test(getResource(t, "test", 3))
	test3.Spec = []byte(`{"tags":["tag3"]}`)

	tests := []corerepository.Test{
		test1,
		test2,
		test3,
	}

	for i := range tests {
		created, err := repository.Create(t.Context(), &tests[i])
		if err != nil {
			t.Fatalf("Unexpected test%d create error: %v", i+1, err)
		}
		if created == nil {
			t.Fatalf("Expected created%d test to be non-nil", i+1)
		}
	}

	testCases := []struct {
		name          string
		tag           string
		context       func() context.Context
		expectedTests []corerepository.Test
		expectedError error
	}{
		{
			name: "matching tag",
			tag:  "shared",
			context: func() context.Context {
				return t.Context()
			},
			expectedTests: []corerepository.Test{
				test1,
				test2,
			},
		},
		{
			name: "empty tag",
			tag:  "",
			context: func() context.Context {
				return t.Context()
			},
			expectedTests: []corerepository.Test{},
		},
		{
			name: "non-existing tag",
			tag:  "does-not-exist",
			context: func() context.Context {
				return t.Context()
			},
			expectedTests: []corerepository.Test{},
		},
		{
			name: "query error",
			tag:  "shared",
			context: func() context.Context {
				ctx, cancel := context.WithCancel(t.Context())
				cancel()

				return ctx
			},
			expectedError: context.Canceled,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := repository.GetByTag(
				testCase.context(),
				testCase.tag,
			)

			if testCase.expectedError != nil {
				if err == nil || !errors.Is(err, testCase.expectedError) {
					t.Fatalf(
						"Expected error %v, got %v",
						testCase.expectedError,
						err,
					)
				}
				if actual != nil {
					t.Fatalf(
						"Expected tests to be nil, got %v",
						actual,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected GetByTag error: %v", err)
			}

			if len(actual) != len(testCase.expectedTests) {
				t.Fatalf(
					"Expected %d tests, got %d",
					len(testCase.expectedTests),
					len(actual),
				)
			}

			for i := range actual {
				assertEqualResources(
					t,
					corerepository.Resource(testCase.expectedTests[i]),
					corerepository.Resource(actual[i]),
				)
			}
		})
	}
}

func TestTestCreateInvalidSpecJSON(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Test

	test := corerepository.Test(getResource(t, "test", 1))
	test.Spec = []byte("{")

	created, err := repository.Create(t.Context(), &test)
	if err == nil {
		t.Fatalf("Expected test create error")
	}
	if created != nil {
		t.Fatalf("Expected created test to be nil")
	}

	got, err := repository.Get(t.Context(), test.ID)
	if !errors.Is(err, corerepository.ErrNotFound) {
		t.Fatalf(
			"Expected error %v, got %v",
			corerepository.ErrNotFound,
			err,
		)
	}
	if got != nil {
		t.Fatalf("Expected test to be nil")
	}
}

func TestTestCreateDuplicateTag(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Test

	test := corerepository.Test(getResource(t, "test", 1))
	test.Spec = []byte(`{"tags":["duplicate","duplicate"]}`)

	created, err := repository.Create(t.Context(), &test)
	if err == nil {
		t.Fatalf("Expected test create error")
	}
	if created != nil {
		t.Fatalf("Expected created test to be nil")
	}
	if !errors.Is(err, corerepository.ErrAlreadyExists) {
		t.Fatalf(
			"Expected error %v, got %v",
			corerepository.ErrAlreadyExists,
			err,
		)
	}

	got, err := repository.Get(t.Context(), test.ID)
	if !errors.Is(err, corerepository.ErrNotFound) {
		t.Fatalf(
			"Expected error %v, got %v",
			corerepository.ErrNotFound,
			err,
		)
	}
	if got != nil {
		t.Fatalf("Expected test to be nil")
	}
}
