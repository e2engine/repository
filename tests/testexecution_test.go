package tests

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	corerepository "github.com/e2engine/core/repository"
)

func TestTestExecutionCRUD(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestExecution
	testSuiteExecutionRepository := repositories.TestSuiteExecution

	// create a testSuiteExecution to use for the testExecutions
	testSuiteExecution := getDefaultTestSuiteExecution(t, 1)
	createdSuite, err := testSuiteExecutionRepository.Create(t.Context(), &testSuiteExecution)
	if err != nil {
		t.Fatalf("Unexpected testSuiteExecution create error: %v", err)
	}
	if createdSuite == nil {
		t.Fatalf("Expected created testSuiteExecution to be non-nil")
	}
	assertEqualTestSuiteExecutions(t, testSuiteExecution, *createdSuite)

	// check that the list is empty
	testExecutions, err := repository.List(
		t.Context(),
		&corerepository.ListParams{
			Limit:          10,
			Position:       corerepository.PositionAfter,
			OrderBy:        corerepository.OrderByID,
			OrderDirection: corerepository.OrderDirectionAsc,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected empty testExecutions list error: %v", err)
	}
	if len(testExecutions) != 0 {
		t.Fatalf("Expected empty testExecutions list, got %d testExecutions", len(testExecutions))
	}

	// create two testexecutions
	testExecution1 := getDefaultTestExecution(t, 1)
	created1, err := repository.Create(t.Context(), &testExecution1)
	if err != nil {
		t.Fatalf("Unexpected testExecution1 create error: %v", err)
	}
	if created1 == nil {
		t.Fatalf("Expected created1 testExecution to be non-nil")
	}
	assertEqualTestExecutions(t, testExecution1, *created1)

	testExecution2 := getDefaultTestExecution(t, 2)
	created2, err := repository.Create(t.Context(), &testExecution2)
	if err != nil {
		t.Fatalf("Unexpected testExecution2 create error: %v", err)
	}
	if created2 == nil {
		t.Fatalf("Expected created2 testExecution to be non-nil")
	}
	assertEqualTestExecutions(t, testExecution2, *created2)

	// list the testExecutions and check that both are present
	testExecutions2, err := repository.List(
		t.Context(),
		&corerepository.ListParams{
			Limit:          10,
			Position:       corerepository.PositionAfter,
			OrderBy:        corerepository.OrderByID,
			OrderDirection: corerepository.OrderDirectionAsc,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected testExecutions list error: %v", err)
	}
	if len(testExecutions2) != 2 {
		t.Fatalf("Expected 2 testExecutions, got %d", len(testExecutions2))
	}
	assertEqualTestExecutions(t, testExecution1, testExecutions2[0])
	assertEqualTestExecutions(t, testExecution2, testExecutions2[1])

	// get the testExecutions by ID and check that they match
	got1, err := repository.Get(t.Context(), testExecution1.ID)
	if err != nil {
		t.Fatalf("Unexpected testExecution1 get error: %v", err)
	}
	if got1 == nil {
		t.Fatalf("Expected testExecution1 to be non-nil")
	}
	assertEqualTestExecutions(t, testExecution1, *got1)

	statusGot1, err := repository.GetStatus(t.Context(), testExecution1.ID)
	if err != nil {
		t.Fatalf("Unexpected testExecution1 get status error: %v", err)
	}
	if statusGot1 != testExecution1.Status {
		t.Fatalf("Expected testExecution1 status to be %q, got %q", testExecution1.Status, statusGot1)
	}

	got2, err := repository.Get(t.Context(), testExecution2.ID)
	if err != nil {
		t.Fatalf("Unexpected testExecution2 get error: %v", err)
	}
	if got2 == nil {
		t.Fatalf("Expected testExecution2 to be non-nil")
	}
	assertEqualTestExecutions(t, testExecution2, *got2)

	statusGot2, err := repository.GetStatus(t.Context(), testExecution2.ID)
	if err != nil {
		t.Fatalf("Unexpected testExecution2 get status error: %v", err)
	}
	if statusGot2 != testExecution2.Status {
		t.Fatalf("Expected testExecution2 status to be %q, got %q", testExecution2.Status, statusGot2)
	}

	// get the testExecutions IDs by idPrefix and check that they match
	ids, err := repository.GetIDsByID(t.Context(), idPrefixTestExecution)
	if err != nil {
		t.Fatalf("Unexpected testExecution IDs by ID error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("Expected 2 testExecution IDs by ID, got %d", len(ids))
	}
	if ids[0] != testExecution1.ID || ids[1] != testExecution2.ID {
		t.Fatalf("Expected testExecution IDs %v, got %v", []string{testExecution1.ID, testExecution2.ID}, ids)
	}

	// get the testExecutions by testSuiteExecutionID prefix and check that they match
	testExecutions3, err := repository.GetByTestSuiteExecutionID(t.Context(), idPrefixTestSuiteExecution+idBody+"1")
	if err != nil {
		t.Fatalf("Unexpected testExecution IDs by testSuiteExecutionID error: %v", err)
	}
	if len(testExecutions3) != 2 {
		t.Fatalf("Expected 2 testExecutions by testSuiteExecutionID, got %d", len(testExecutions3))
	}
	assertEqualTestExecutions(t, testExecution1, testExecutions3[0])
	assertEqualTestExecutions(t, testExecution2, testExecutions3[1])

	// set the first testExecution to running and check that it was updated
	testExecution1.StartedAt = testExecution1.StartedAt.Add(time.Hour)
	err = repository.SetRunning(
		t.Context(),
		&corerepository.SetRunningParams{
			ID:        testExecution1.ID,
			StartedAt: testExecution1.StartedAt,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected testExecution1 set running error: %v", err)
	}
	testExecution1.Status = "running"

	running1, err := repository.Get(t.Context(), testExecution1.ID)
	if err != nil {
		t.Fatalf("Unexpected testExecution1 get after set running error: %v", err)
	}
	if running1 == nil {
		t.Fatalf("Expected running1 testExecution to be non-nil")
	}
	assertEqualTestExecutions(t, testExecution1, *running1)

	statusRunning1, err := repository.GetStatus(t.Context(), testExecution1.ID)
	if err != nil {
		t.Fatalf("Unexpected testExecution1 get status error: %v", err)
	}
	if statusRunning1 != testExecution1.Status {
		t.Fatalf("Expected testExecution1 status to be %q, got %q", testExecution1.Status, statusRunning1)
	}

	// set the second testExecution to completed and check that it was updated
	testExecution2.FinishedAt = testExecution2.FinishedAt.Add(time.Hour)
	testExecution2.Status = "failed"
	testExecution2.Summary = []byte(`{"result":"failed"}`)
	err = repository.SetCompleted(
		t.Context(),
		&corerepository.SetCompletedParams{
			ID:         testExecution2.ID,
			FinishedAt: testExecution2.FinishedAt,
			Status:     testExecution2.Status,
			Summary:    testExecution2.Summary,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected testExecution2 set completed error: %v", err)
	}

	completed2, err := repository.Get(t.Context(), testExecution2.ID)
	if err != nil {
		t.Fatalf("Unexpected testExecution2 get after set completed error: %v", err)
	}
	if completed2 == nil {
		t.Fatalf("Expected completed2 testExecution to be non-nil")
	}
	assertEqualTestExecutions(t, testExecution2, *completed2)

	statusCompleted2, err := repository.GetStatus(t.Context(), testExecution2.ID)
	if err != nil {
		t.Fatalf("Unexpected testExecution2 get status error: %v", err)
	}
	if statusCompleted2 != testExecution2.Status {
		t.Fatalf("Expected testExecution2 status to be %q, got %q", testExecution2.Status, statusCompleted2)
	}

	// delete the testExecutions by their IDs
	deleted1, err := repository.Delete(t.Context(), testExecution1.ID)
	if err != nil {
		t.Fatalf("Unexpected testExecution1 delete error: %v", err)
	}
	if deleted1 == nil {
		t.Fatalf("Expected deleted1 testExecution to be non-nil")
	}
	assertEqualTestExecutions(t, testExecution1, *deleted1)

	deleted2, err := repository.Delete(t.Context(), testExecution2.ID)
	if err != nil {
		t.Fatalf("Unexpected testExecution2 delete error: %v", err)
	}
	if deleted2 == nil {
		t.Fatalf("Expected deleted2 testExecution to be non-nil")
	}
	assertEqualTestExecutions(t, testExecution2, *deleted2)

	// check that the list is empty again
	testExecutions4, err := repository.List(
		t.Context(),
		&corerepository.ListParams{
			Limit:          10,
			Position:       corerepository.PositionAfter,
			OrderBy:        corerepository.OrderByID,
			OrderDirection: corerepository.OrderDirectionAsc,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected empty testExecutions list error: %v", err)
	}
	if len(testExecutions4) != 0 {
		t.Fatalf("Expected empty testExecutions list, got %d testExecutions", len(testExecutions4))
	}
}

func TestTestExecutionStandaloneRoundTrip(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestExecution

	expected := getDefaultTestExecution(t, 1)
	expected.TestSuiteExecutionID = ""

	created, err := repository.Create(t.Context(), &expected)
	if err != nil {
		t.Fatalf("Unexpected standalone testExecution create error: %v", err)
	}
	if created == nil {
		t.Fatalf("Expected created testExecution to be non-nil")
	}
	assertEqualTestExecutions(t, expected, *created)

	got, err := repository.Get(t.Context(), expected.ID)
	if err != nil {
		t.Fatalf("Unexpected standalone testExecution get error: %v", err)
	}
	if got == nil {
		t.Fatalf("Expected standalone testExecution to be non-nil")
	}
	assertEqualTestExecutions(t, expected, *got)

	if got.TestSuiteExecutionID != "" {
		t.Fatalf(
			"Expected standalone testExecution TestSuiteExecutionID to be empty, got %q",
			got.TestSuiteExecutionID,
		)
	}
}

func TestTestExecutionCreateDuplicateID(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestExecution

	testExecution := getDefaultTestExecution(t, 1)
	testExecution.TestSuiteExecutionID = ""

	created, err := repository.Create(t.Context(), &testExecution)
	if err != nil {
		t.Fatalf("Unexpected testExecution create error: %v", err)
	}
	if created == nil {
		t.Fatalf("Expected created testExecution to be non-nil")
	}

	duplicate := getDefaultTestExecution(t, 2)
	duplicate.ID = testExecution.ID
	duplicate.TestSuiteExecutionID = ""

	created, err = repository.Create(t.Context(), &duplicate)
	if err == nil {
		t.Fatalf("Expected duplicate testExecution create error")
	}
	if created != nil {
		t.Fatalf("Expected duplicate testExecution to be nil")
	}

	if !errors.Is(err, corerepository.ErrAlreadyExists) {
		t.Fatalf(
			"Expected error %v, got %v",
			corerepository.ErrAlreadyExists,
			err,
		)
	}
}

func TestTestExecutionCreateErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestExecution

	testCases := []struct {
		name          string
		mutate        func(*corerepository.TestExecution)
		expectedError error
	}{
		{
			name: "empty ID",
			mutate: func(testExecution *corerepository.TestExecution) {
				testExecution.ID = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "invalid status",
			mutate: func(testExecution *corerepository.TestExecution) {
				testExecution.Status = "invalid"
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "empty environment_id",
			mutate: func(testExecution *corerepository.TestExecution) {
				testExecution.EnvironmentID = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "empty environment_name",
			mutate: func(testExecution *corerepository.TestExecution) {
				testExecution.EnvironmentName = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "empty test_id",
			mutate: func(testExecution *corerepository.TestExecution) {
				testExecution.TestID = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "empty test_name",
			mutate: func(testExecution *corerepository.TestExecution) {
				testExecution.TestName = ""
			},
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "invalid test_suite_execution_id",
			mutate: func(testExecution *corerepository.TestExecution) {
				testExecution.TestSuiteExecutionID = "invalid"
			},
			expectedError: corerepository.ErrConflict,
		},
	}

	for i, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testExecution := getDefaultTestExecution(t, i+1)
			testCase.mutate(&testExecution)

			created, err := repository.Create(t.Context(), &testExecution)
			if err == nil {
				t.Fatalf("Expected testExecution create error")
			}
			if created != nil {
				t.Fatalf("Expected created testExecution to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestTestExecutionGetErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestExecution

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
			testExecution, err := repository.Get(t.Context(), testCase.id)
			if err == nil {
				t.Fatalf("Expected testExecution get error")
			}
			if testExecution != nil {
				t.Fatalf("Expected testExecution to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestTestExecutionGetStatusErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestExecution

	testCases := []struct {
		name          string
		context       func() context.Context
		id            string
		expectedError error
	}{
		{
			name: "empty ID",
			context: func() context.Context {
				return t.Context()
			},
			id:            "",
			expectedError: corerepository.ErrNotFound,
		},
		{
			name: "non-existing ID",
			context: func() context.Context {
				return t.Context()
			},
			id:            "non-existing-id",
			expectedError: corerepository.ErrNotFound,
		},
		{
			name: "query error",
			context: func() context.Context {
				ctx, cancel := context.WithCancel(t.Context())
				cancel()

				return ctx
			},
			id:            "f5497dfbbf8392b4828f5323c07e4b5cc8026790caca7b3065350413bebbd4e3",
			expectedError: context.Canceled,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repository.GetStatus(testCase.context(), testCase.id)
			if err == nil {
				t.Fatalf("Expected testExecution get status error")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestTestExecutionDeleteErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestExecution

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
			testExecution, err := repository.Delete(t.Context(), testCase.id)
			if err == nil {
				t.Fatalf("Expected testExecution delete error")
			}
			if testExecution != nil {
				t.Fatalf("Expected testExecution to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestTestExecutionGetByTestSuiteExecutionID(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestExecution
	testSuiteExecutionRepository := repositories.TestSuiteExecution

	testSuiteExecution1 := getDefaultTestSuiteExecution(t, 1)
	createdSuite1, err := testSuiteExecutionRepository.Create(
		t.Context(),
		&testSuiteExecution1,
	)
	if err != nil {
		t.Fatalf("Unexpected testSuiteExecution1 create error: %v", err)
	}
	if createdSuite1 == nil {
		t.Fatalf("Expected created1 testSuiteExecution to be non-nil")
	}

	testSuiteExecution2 := getDefaultTestSuiteExecution(t, 2)
	createdSuite2, err := testSuiteExecutionRepository.Create(
		t.Context(),
		&testSuiteExecution2,
	)
	if err != nil {
		t.Fatalf("Unexpected testSuiteExecution2 create error: %v", err)
	}
	if createdSuite2 == nil {
		t.Fatalf("Expected created2 testSuiteExecution to be non-nil")
	}

	testExecution1 := getDefaultTestExecution(t, 1)
	testExecution1.TestSuiteExecutionID = testSuiteExecution1.ID
	testExecution2 := getDefaultTestExecution(t, 2)
	testExecution2.TestSuiteExecutionID = testSuiteExecution1.ID
	testExecution3 := getDefaultTestExecution(t, 3)
	testExecution3.TestSuiteExecutionID = testSuiteExecution2.ID
	createTestExecutions(
		t,
		repository,
		[]corerepository.TestExecution{
			testExecution1,
			testExecution2,
			testExecution3,
		},
	)

	testCases := []struct {
		name            string
		context         func() context.Context
		testSuiteExecID string
		expected        []corerepository.TestExecution
		expectedError   error
	}{
		{
			name: "matching test suite execution",
			context: func() context.Context {
				return t.Context()
			},
			testSuiteExecID: testSuiteExecution1.ID,
			expected:        []corerepository.TestExecution{testExecution1, testExecution2},
		},
		{
			name: "test suite execution with one test execution",
			context: func() context.Context {
				return t.Context()
			},
			testSuiteExecID: testSuiteExecution2.ID,
			expected:        []corerepository.TestExecution{testExecution3},
		},
		{
			name: "empty test suite execution ID",
			context: func() context.Context {
				return t.Context()
			},
			testSuiteExecID: "",
			expected:        []corerepository.TestExecution{},
		},
		{
			name: "non-existing test suite execution ID",
			context: func() context.Context {
				return t.Context()
			},
			testSuiteExecID: "does-not-exist",
			expected:        []corerepository.TestExecution{},
		},
		{
			name: "query error",
			context: func() context.Context {
				ctx, cancel := context.WithCancel(t.Context())
				cancel()

				return ctx
			},
			testSuiteExecID: testSuiteExecution1.ID,
			expectedError:   context.Canceled,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testExecutions, err := repository.GetByTestSuiteExecutionID(
				testCase.context(),
				testCase.testSuiteExecID,
			)

			if testCase.expectedError != nil {
				if err == nil || !errors.Is(err, testCase.expectedError) {
					t.Fatalf(
						"Expected error %v, got %v",
						testCase.expectedError,
						err,
					)
				}
				if testExecutions != nil {
					t.Fatalf("Expected testExecutions to be nil, got %v", testExecutions)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected test executions lookup error: %v", err)
			}
			if len(testExecutions) != len(testCase.expected) {
				t.Fatalf(
					"Expected %d testExecutions, got %d",
					len(testCase.expected),
					len(testExecutions),
				)
			}
			for i := range testCase.expected {
				assertEqualTestExecutions(t, testCase.expected[i], testExecutions[i])
			}
		})
	}
}

func TestTestExecutionSetRunning(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestExecution
	testSuiteExecution := getDefaultTestSuiteExecution(t, 1)
	createdSuite, err := repositories.TestSuiteExecution.Create(t.Context(), &testSuiteExecution)
	if err != nil {
		t.Fatalf("Unexpected testSuiteExecution create error: %v", err)
	}
	if createdSuite == nil {
		t.Fatalf("Expected created testSuiteExecution to be non-nil")
	}

	testExecutions := getTestExecutionsForList(t, testSuiteExecution.ID)
	createTestExecutions(t, repository, testExecutions[:2])

	testCases := []struct {
		name          string
		params        *corerepository.SetRunningParams
		expected      corerepository.TestExecution
		expectedError error
	}{
		{
			name: "first test execution",
			params: &corerepository.SetRunningParams{
				ID:        testExecutions[0].ID,
				StartedAt: testExecutions[0].StartedAt.Add(time.Hour),
			},
			expected: func() corerepository.TestExecution {
				expected := testExecutions[0]
				expected.StartedAt = expected.StartedAt.Add(time.Hour)
				expected.Status = "running"

				return expected
			}(),
		},
		{
			name: "second test execution",
			params: &corerepository.SetRunningParams{
				ID:        testExecutions[1].ID,
				StartedAt: testExecutions[1].StartedAt.Add(2 * time.Hour),
			},
			expected: func() corerepository.TestExecution {
				expected := testExecutions[1]
				expected.StartedAt = expected.StartedAt.Add(2 * time.Hour)
				expected.Status = "running"

				return expected
			}(),
		},
		{
			name: "nonexistent test execution",
			params: &corerepository.SetRunningParams{
				ID:        "does-not-exist",
				StartedAt: time.Now().UTC(),
			},
			expectedError: corerepository.ErrNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := repository.SetRunning(t.Context(), testCase.params)
			if testCase.expectedError != nil {
				if err == nil {
					t.Fatalf(
						"Expected SetRunning to return error: %v, but got nil",
						testCase.expectedError,
					)
				}
				if !errors.Is(err, testCase.expectedError) {
					t.Fatalf(
						"Expected SetRunning to return error: %v but got: %v",
						testCase.expectedError,
						err,
					)
				}

				return
			}
			if err != nil {
				t.Fatalf("Unexpected test execution SetRunning error: %v", err)
			}

			testExecution, err := repository.Get(t.Context(), testCase.params.ID)
			if err != nil {
				t.Fatalf("Unexpected test execution get error: %v", err)
			}
			if testExecution == nil {
				t.Fatalf("Expected test execution to be non-nil")
			}

			assertEqualTestExecutions(t, testCase.expected, *testExecution)
		})
	}
}

func TestTestExecutionSetCompleted(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestExecution
	testSuiteExecution := getDefaultTestSuiteExecution(t, 1)
	createdSuite, err := repositories.TestSuiteExecution.Create(t.Context(), &testSuiteExecution)
	if err != nil {
		t.Fatalf("Unexpected testSuiteExecution create error: %v", err)
	}
	if createdSuite == nil {
		t.Fatalf("Expected created testSuiteExecution to be non-nil")
	}

	testExecutions := getTestExecutionsForList(t, testSuiteExecution.ID)
	createTestExecutions(t, repository, testExecutions[:4])

	testCases := []struct {
		name          string
		params        *corerepository.SetCompletedParams
		expected      corerepository.TestExecution
		expectedError error
	}{
		{
			name: "passed test execution",
			params: &corerepository.SetCompletedParams{
				ID:         testExecutions[0].ID,
				FinishedAt: testExecutions[0].FinishedAt.Add(time.Hour),
				Status:     "passed",
				Summary:    []byte(`{"result":"passed"}`),
			},
			expected: func() corerepository.TestExecution {
				expected := testExecutions[0]
				expected.FinishedAt = expected.FinishedAt.Add(time.Hour)
				expected.Status = "passed"
				expected.Summary = []byte(`{"result":"passed"}`)

				return expected
			}(),
		},
		{
			name: "failed test execution",
			params: &corerepository.SetCompletedParams{
				ID:         testExecutions[1].ID,
				FinishedAt: testExecutions[1].FinishedAt.Add(2 * time.Hour),
				Status:     "failed",
				Summary:    []byte(`{"result":"failed"}`),
			},
			expected: func() corerepository.TestExecution {
				expected := testExecutions[1]
				expected.FinishedAt = expected.FinishedAt.Add(2 * time.Hour)
				expected.Status = "failed"
				expected.Summary = []byte(`{"result":"failed"}`)

				return expected
			}(),
		},
		{
			name: "errored test execution",
			params: &corerepository.SetCompletedParams{
				ID:         testExecutions[2].ID,
				FinishedAt: testExecutions[2].FinishedAt.Add(3 * time.Hour),
				Status:     "error",
				Summary:    []byte(`{"result":"error"}`),
			},
			expected: func() corerepository.TestExecution {
				expected := testExecutions[2]
				expected.FinishedAt = expected.FinishedAt.Add(3 * time.Hour)
				expected.Status = "error"
				expected.Summary = []byte(`{"result":"error"}`)

				return expected
			}(),
		},
		{
			name: "test execution invalid status",
			params: &corerepository.SetCompletedParams{
				ID:         testExecutions[3].ID,
				FinishedAt: testExecutions[3].FinishedAt.Add(4 * time.Hour),
				Status:     "invalid",
				Summary:    []byte(`{"result":"error"}`),
			},
			expected:      testExecutions[3],
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name: "nonexistent test execution",
			params: &corerepository.SetCompletedParams{
				ID:         "does-not-exist",
				FinishedAt: time.Now().UTC(),
				Status:     "passed",
				Summary:    []byte(`{"result":"passed"}`),
			},
			expectedError: corerepository.ErrNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := repository.SetCompleted(t.Context(), testCase.params)

			if testCase.expectedError != nil {
				if err == nil {
					t.Fatalf(
						"Expected SetCompleted to return error: %v, but got nil",
						testCase.expectedError,
					)
				}
				if !errors.Is(err, testCase.expectedError) {
					t.Fatalf(
						"Expected SetCompleted to return error: %v but got: %v",
						testCase.expectedError,
						err,
					)
				}

				if errors.Is(testCase.expectedError, corerepository.ErrNotFound) {
					return
				}
			} else if err != nil {
				t.Fatalf("Unexpected test execution SetCompleted error: %v", err)
			}

			testExecution, err := repository.Get(t.Context(), testCase.params.ID)
			if err != nil {
				t.Fatalf("Unexpected test execution get error: %v", err)
			}
			if testExecution == nil {
				t.Fatalf("Expected test execution to be non-nil")
			}

			assertEqualTestExecutions(t, testCase.expected, *testExecution)
		})
	}
}

func TestTestExecutionGetIDsByID(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestExecution

	testSuiteExecution := getDefaultTestSuiteExecution(t, 1)
	createdSuite, err := repositories.TestSuiteExecution.Create(t.Context(), &testSuiteExecution)
	if err != nil {
		t.Fatalf("Unexpected testSuiteExecution create error: %v", err)
	}
	if createdSuite == nil {
		t.Fatalf("Expected created testSuiteExecution to be non-nil")
	}

	tes := getTestExecutionsForList(t, testSuiteExecution.ID)
	createTestExecutions(t, repository, tes)

	allIDs := make([]string, len(tes))
	for i := range tes {
		allIDs[i] = tes[i].ID
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
			id:   tes[0].ID,
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

func TestTestExecutionCascadeDeleteWithTestSuiteExecution(t *testing.T) {
	repositories := getRepositories(t)
	testExecutionRepository := repositories.TestExecution
	testSuiteExecutionRepository := repositories.TestSuiteExecution

	testSuiteExecution := getDefaultTestSuiteExecution(t, 1)
	createdSuite, err := testSuiteExecutionRepository.Create(
		t.Context(),
		&testSuiteExecution,
	)
	if err != nil {
		t.Fatalf("Unexpected testSuiteExecution create error: %v", err)
	}
	if createdSuite == nil {
		t.Fatalf("Expected created testSuiteExecution to be non-nil")
	}

	testExecutions := getTestExecutionsForList(t, testSuiteExecution.ID)
	createTestExecutions(t, testExecutionRepository, testExecutions[:3])

	beforeDelete, err := testExecutionRepository.GetByTestSuiteExecutionID(
		t.Context(),
		testSuiteExecution.ID,
	)
	if err != nil {
		t.Fatalf(
			"Unexpected GetByTestSuiteExecutionID error before parent deletion: %v",
			err,
		)
	}
	if len(beforeDelete) != 3 {
		t.Fatalf(
			"Expected 3 test executions before parent deletion, got %d",
			len(beforeDelete),
		)
	}

	deletedSuite, err := testSuiteExecutionRepository.Delete(
		t.Context(),
		testSuiteExecution.ID,
	)
	if err != nil {
		t.Fatalf("Unexpected testSuiteExecution delete error: %v", err)
	}
	if deletedSuite == nil {
		t.Fatalf("Expected deleted testSuiteExecution to be non-nil")
	}

	afterDelete, err := testExecutionRepository.GetByTestSuiteExecutionID(
		t.Context(),
		testSuiteExecution.ID,
	)
	if err != nil {
		t.Fatalf(
			"Unexpected GetByTestSuiteExecutionID error after parent deletion: %v",
			err,
		)
	}
	if len(afterDelete) != 0 {
		t.Fatalf(
			"Expected test executions to be deleted with parent, got %d",
			len(afterDelete),
		)
	}

	for _, expected := range testExecutions[:3] {
		actual, err := testExecutionRepository.Get(t.Context(), expected.ID)
		if err == nil {
			t.Fatalf(
				"Expected test execution %q to be deleted, but Get returned no error",
				expected.ID,
			)
		}
		if !errors.Is(err, corerepository.ErrNotFound) {
			t.Fatalf(
				"Expected test execution %q Get to return %v, got: %v",
				expected.ID,
				corerepository.ErrNotFound,
				err,
			)
		}
		if actual != nil {
			t.Fatalf(
				"Expected deleted test execution %q to be nil",
				expected.ID,
			)
		}
	}
}
