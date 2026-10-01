package tests

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	corerepository "github.com/e2engine/core/repository"
)

func TestEnvironmentCRUD(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment

	// check that the list is empty
	environments, err := repository.List(
		t.Context(),
		&corerepository.ListParams{
			Limit:          10,
			Position:       corerepository.PositionAfter,
			OrderBy:        corerepository.OrderByID,
			OrderDirection: corerepository.OrderDirectionAsc,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected empty environments list error: %v", err)
	}
	if len(environments) != 0 {
		t.Fatalf("Expected empty environments list, got %d environments", len(environments))
	}

	// create two environments
	environment1 := getResource(t, "environment", 1)
	created1, err := repository.Create(t.Context(), (*corerepository.Environment)(&environment1))
	if err != nil {
		t.Fatalf("Unexpected environment1 create error: %v", err)
	}
	if created1 == nil {
		t.Fatalf("Expected created1 environment to be non-nil")
	}
	assertEqualResources(t, environment1, corerepository.Resource(*created1))

	environment2 := getResource(t, "environment", 2)
	created2, err := repository.Create(t.Context(), (*corerepository.Environment)(&environment2))
	if err != nil {
		t.Fatalf("Unexpected environment2 create error: %v", err)
	}
	if created2 == nil {
		t.Fatalf("Expected created2 environment to be non-nil")
	}
	assertEqualResources(t, environment2, corerepository.Resource(*created2))

	// list the environments and check that both are present
	environments, err = repository.List(
		t.Context(),
		&corerepository.ListParams{
			Limit:          10,
			Position:       corerepository.PositionAfter,
			OrderBy:        corerepository.OrderByID,
			OrderDirection: corerepository.OrderDirectionAsc,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected environments list error: %v", err)
	}
	if len(environments) != 2 {
		t.Fatalf("Expected 2 environments, got %d", len(environments))
	}
	assertEqualResources(t, environment1, corerepository.Resource(environments[0]))
	assertEqualResources(t, environment2, corerepository.Resource(environments[1]))

	// get the environments by ID and check that they match
	got1, err := repository.Get(t.Context(), environment1.ID)
	if err != nil {
		t.Fatalf("Unexpected environment1 get error: %v", err)
	}
	if got1 == nil {
		t.Fatalf("Expected environment1 to be non-nil")
	}
	assertEqualResources(t, environment1, corerepository.Resource(*got1))

	got2, err := repository.Get(t.Context(), environment2.ID)
	if err != nil {
		t.Fatalf("Unexpected environment2 get error: %v", err)
	}
	if got2 == nil {
		t.Fatalf("Expected environment2 to be non-nil")
	}
	assertEqualResources(t, environment2, corerepository.Resource(*got2))

	// get the environments IDs by idPrefix and check that they match
	ids, err := repository.GetIDsByID(t.Context(), idPrefixEnvironment)
	if err != nil {
		t.Fatalf("Unexpected environment IDs by ID error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("Expected 2 environment IDs by ID, got %d", len(ids))
	}
	if ids[0] != environment1.ID || ids[1] != environment2.ID {
		t.Fatalf("Expected environment IDs %v, got %v", []string{environment1.ID, environment2.ID}, ids)
	}

	// get the environments IDs by name prefix and check that they match
	ids, err = repository.GetIDsByName(t.Context(), "Resource environment")
	if err != nil {
		t.Fatalf("Unexpected environment IDs by name error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("Expected 2 environment IDs by name, got %d", len(ids))
	}
	if ids[0] != environment1.ID || ids[1] != environment2.ID {
		t.Fatalf("Expected environment IDs %v, got %v", []string{environment1.ID, environment2.ID}, ids)
	}

	// delete the environments by their IDs
	deleted1, err := repository.Delete(t.Context(), environment1.ID)
	if err != nil {
		t.Fatalf("Unexpected environment1 delete error: %v", err)
	}
	if deleted1 == nil {
		t.Fatalf("Expected deleted1 environment to be non-nil")
	}
	assertEqualResources(t, environment1, corerepository.Resource(*deleted1))

	deleted2, err := repository.Delete(t.Context(), environment2.ID)
	if err != nil {
		t.Fatalf("Unexpected environment2 delete error: %v", err)
	}
	if deleted2 == nil {
		t.Fatalf("Expected deleted2 environment to be non-nil")
	}
	assertEqualResources(t, environment2, corerepository.Resource(*deleted2))

	// check that the list is empty again
	environments, err = repository.List(
		t.Context(),
		&corerepository.ListParams{
			Limit:          10,
			Position:       corerepository.PositionAfter,
			OrderBy:        corerepository.OrderByID,
			OrderDirection: corerepository.OrderDirectionAsc,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected empty environments list error: %v", err)
	}
	if len(environments) != 0 {
		t.Fatalf("Expected empty environments list, got %d environments", len(environments))
	}
}

func TestEnvironmentCreateErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment

	testCases := []struct {
		name          string
		mutate        func(*corerepository.Environment)
		expectedError error
	}{
		{
			name: "empty ID",
			mutate: func(environment *corerepository.Environment) {
				environment.ID = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "empty version",
			mutate: func(environment *corerepository.Environment) {
				environment.Version = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "empty name",
			mutate: func(environment *corerepository.Environment) {
				environment.Name = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "empty spec",
			mutate: func(environment *corerepository.Environment) {
				environment.Spec = []byte{}
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "zero created at",
			mutate: func(environment *corerepository.Environment) {
				environment.CreatedAt = time.Time{}
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "zero updated at",
			mutate: func(environment *corerepository.Environment) {
				environment.UpdatedAt = time.Time{}
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
	}

	for i, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			environment := corerepository.Environment(
				getResource(t, "environment", i+1),
			)
			testCase.mutate(&environment)

			created, err := repository.Create(t.Context(), &environment)
			if err == nil {
				t.Fatalf("Expected environment create error")
			}
			if created != nil {
				t.Fatalf("Expected created environment to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestEnvironmentCreateEmptyDescription(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment

	environment := corerepository.Environment(
		getResource(t, "environment", 1),
	)
	environment.Description = ""

	created, err := repository.Create(t.Context(), &environment)
	if err != nil {
		t.Fatalf("Unexpected environment create error: %v", err)
	}
	if created == nil {
		t.Fatalf("Expected created environment to be non-nil")
	}

	assertEqualResources(
		t,
		corerepository.Resource(environment),
		corerepository.Resource(*created),
	)
}

func TestEnvironmentCreateDuplicatedName(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment

	environment1 := corerepository.Environment(getResource(t, "environment", 1))
	created1, err := repository.Create(t.Context(), &environment1)
	if err != nil {
		t.Fatalf("Unexpected environment1 create error: %v", err)
	}
	if created1 == nil {
		t.Fatalf("Expected created1 environment to be non-nil")
	}

	environment2 := corerepository.Environment(getResource(t, "environment", 2))
	environment2.Name = environment1.Name
	created2, err := repository.Create(t.Context(), &environment2)
	if err == nil {
		t.Fatalf("Expected environment2 create error")
	}
	if created2 != nil {
		t.Fatalf("Expected created2 environment to be nil")
	}
	if !errors.Is(err, corerepository.ErrAlreadyExists) {
		t.Fatalf(
			"Expected error %v, got %v",
			corerepository.ErrAlreadyExists,
			err,
		)
	}
}

func TestEnvironmentCreateIgnoresPayloadKind(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment

	environment := corerepository.Environment(getResource(t, "invalid", 1))
	created, err := repository.Create(t.Context(), &environment)
	if err != nil {
		t.Fatalf("Unexpected environment create error: %v", err)
	}
	if created == nil {
		t.Fatalf("Expected created environment to be non-nil")
	}
	if created.Kind != "environment" {
		t.Fatalf("Expected created environment kind to be environment, got %q", created.Kind)
	}
}

func TestEnvironmentGetErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment

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
			environment, err := repository.Get(t.Context(), testCase.id)
			if err == nil {
				t.Fatalf("Expected environment get error")
			}
			if environment != nil {
				t.Fatalf("Expected environment to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestEnvironmentDeleteErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment

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
			environment, err := repository.Delete(t.Context(), testCase.id)
			if err == nil {
				t.Fatalf("Expected environment delete error")
			}
			if environment != nil {
				t.Fatalf("Expected environment to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestEnvironmentGetIDsByID(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment
	resources := getResourcesForList(t, "environment")
	createEnvironmentResources(t, repository, resources)

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

func TestEnvironmentGetIDsByName(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment
	resources := getResourcesForList(t, "environment")
	createEnvironmentResources(t, repository, resources)

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
