package tests

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	corerepository "github.com/e2engine/core/repository"
)

func TestTestSuiteExecutionCRUD(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuiteExecution

	executions, err := repository.List(t.Context(), &corerepository.ListParams{
		Limit: 10, Position: corerepository.PositionAfter,
		OrderBy: corerepository.OrderByID, OrderDirection: corerepository.OrderDirectionAsc,
	})
	if err != nil {
		t.Fatalf("Unexpected empty test suite executions list error: %v", err)
	}
	if len(executions) != 0 {
		t.Fatalf("Expected empty test suite executions list, got %d test suite executions", len(executions))
	}

	execution1 := createTestSuiteExecution(t, repository, 1)
	execution2 := createTestSuiteExecution(t, repository, 2)

	childExecution1 := getDefaultTestExecution(t, 1)
	childExecution1.TestSuiteExecutionID = execution1.ID
	childExecution2 := getDefaultTestExecution(t, 2)
	childExecution2.TestSuiteExecutionID = execution1.ID
	neighborChildExecution := getDefaultTestExecution(t, 3)
	neighborChildExecution.TestSuiteExecutionID = execution2.ID
	createTestExecutions(
		t,
		repositories.TestExecution,
		[]corerepository.TestExecution{
			childExecution1,
			childExecution2,
			neighborChildExecution,
		},
	)
	execution1.TestsCount = 2
	execution2.TestsCount = 1

	executions, err = repository.List(t.Context(), &corerepository.ListParams{
		Limit: 10, Position: corerepository.PositionAfter,
		OrderBy: corerepository.OrderByID, OrderDirection: corerepository.OrderDirectionAsc,
	})
	if err != nil {
		t.Fatalf("Unexpected test suite executions list error: %v", err)
	}
	assertEqualTestSuiteExecutionLists(
		t,
		[]corerepository.TestSuiteExecution{execution1, execution2},
		executions,
	)

	for _, expected := range []corerepository.TestSuiteExecution{execution1, execution2} {
		actual, err := repository.Get(t.Context(), expected.ID)
		if err != nil {
			t.Fatalf("Unexpected test suite execution get error: %v", err)
		}
		if actual == nil {
			t.Fatalf("Expected test suite execution to be non-nil")
		}
		assertEqualTestSuiteExecutions(t, expected, *actual)

		statusGot, err := repository.GetStatus(t.Context(), actual.ID)
		if err != nil {
			t.Fatalf("Unexpected testExecution1 get status error: %v", err)
		}
		if statusGot != actual.Status {
			t.Fatalf("Expected testExecution1 status to be %q, got %q", actual.Status, statusGot)
		}
	}

	ids, err := repository.GetIDsByID(t.Context(), idPrefixTestSuiteExecution)
	if err != nil {
		t.Fatalf("Unexpected test suite execution IDs by ID error: %v", err)
	}
	if !slices.Equal(ids, []string{execution1.ID, execution2.ID}) {
		t.Fatalf("Expected test suite execution IDs %v, got %v", []string{execution1.ID, execution2.ID}, ids)
	}

	execution1.StartedAt = execution1.StartedAt.Add(time.Hour)
	execution1.Status = "running"
	err = repository.SetRunning(t.Context(), &corerepository.SetRunningParams{
		ID: execution1.ID, StartedAt: execution1.StartedAt,
	})
	if err != nil {
		t.Fatalf("Unexpected test suite execution set running error: %v", err)
	}
	actual1, err := repository.Get(t.Context(), execution1.ID)
	if err != nil {
		t.Fatalf("Unexpected test suite execution get after set running error: %v", err)
	}
	if actual1 == nil {
		t.Fatalf("Expected running test suite execution to be non-nil")
	}
	assertEqualTestSuiteExecutions(t, execution1, *actual1)

	statusRunning, err := repository.GetStatus(t.Context(), actual1.ID)
	if err != nil {
		t.Fatalf("Unexpected testExecution1 get status error: %v", err)
	}
	if statusRunning != actual1.Status {
		t.Fatalf("Expected testExecution1 status to be %q, got %q", actual1.Status, statusRunning)
	}

	execution2.FinishedAt = execution2.FinishedAt.Add(time.Hour)
	execution2.Status = "failed"
	execution2.Summary = []byte(`{"result":"failed"}`)
	err = repository.SetCompleted(t.Context(), &corerepository.SetCompletedParams{
		ID: execution2.ID, FinishedAt: execution2.FinishedAt,
		Status: execution2.Status, Summary: execution2.Summary,
	})
	if err != nil {
		t.Fatalf("Unexpected test suite execution set completed error: %v", err)
	}
	actual2, err := repository.Get(t.Context(), execution2.ID)
	if err != nil {
		t.Fatalf("Unexpected test suite execution get after set completed error: %v", err)
	}
	if actual2 == nil {
		t.Fatalf("Expected completed test suite execution to be non-nil")
	}
	assertEqualTestSuiteExecutions(t, execution2, *actual2)

	statusCompleted2, err := repository.GetStatus(t.Context(), actual2.ID)
	if err != nil {
		t.Fatalf("Unexpected testExecution2 get status error: %v", err)
	}
	if statusCompleted2 != actual2.Status {
		t.Fatalf("Expected testExecution2 status to be %q, got %q", actual2.Status, statusCompleted2)
	}

	for _, childExecution := range []corerepository.TestExecution{
		childExecution1,
		childExecution2,
		neighborChildExecution,
	} {
		deleted, err := repositories.TestExecution.Delete(t.Context(), childExecution.ID)
		if err != nil {
			t.Fatalf("Unexpected child test execution delete error: %v", err)
		}
		if deleted == nil {
			t.Fatalf("Expected deleted child test execution to be non-nil")
		}
		assertEqualTestExecutions(t, childExecution, *deleted)
	}
	execution1.TestsCount = 0
	execution2.TestsCount = 0

	for _, expected := range []corerepository.TestSuiteExecution{execution1, execution2} {
		deleted, err := repository.Delete(t.Context(), expected.ID)
		if err != nil {
			t.Fatalf("Unexpected test suite execution delete error: %v", err)
		}
		if deleted == nil {
			t.Fatalf("Expected deleted test suite execution to be non-nil")
		}
		assertEqualTestSuiteExecutions(t, expected, *deleted)
	}

	executions, err = repository.List(t.Context(), &corerepository.ListParams{
		Limit: 10, Position: corerepository.PositionAfter,
		OrderBy: corerepository.OrderByID, OrderDirection: corerepository.OrderDirectionAsc,
	})
	if err != nil {
		t.Fatalf("Unexpected empty test suite executions list error: %v", err)
	}
	if len(executions) != 0 {
		t.Fatalf("Expected empty test suite executions list, got %d test suite executions", len(executions))
	}
}

func TestTestSuiteExecutionCreateErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuiteExecution
	testCases := []struct {
		name          string
		mutate        func(*corerepository.TestSuiteExecution)
		expectedError error
	}{
		{
			name:          "empty ID",
			mutate:        func(execution *corerepository.TestSuiteExecution) { execution.ID = "" },
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name:          "invalid status",
			mutate:        func(execution *corerepository.TestSuiteExecution) { execution.Status = "invalid" },
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name:          "empty environment ID",
			mutate:        func(execution *corerepository.TestSuiteExecution) { execution.EnvironmentID = "" },
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name:          "empty environment name",
			mutate:        func(execution *corerepository.TestSuiteExecution) { execution.EnvironmentName = "" },
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name:          "empty test suite ID",
			mutate:        func(execution *corerepository.TestSuiteExecution) { execution.TestSuiteID = "" },
			expectedError: corerepository.ErrInvalidPayload,
		},
		{
			name:          "empty test suite name",
			mutate:        func(execution *corerepository.TestSuiteExecution) { execution.TestSuiteName = "" },
			expectedError: corerepository.ErrInvalidPayload,
		},
	}

	for i, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			execution := getDefaultTestSuiteExecution(t, i+1)
			testCase.mutate(&execution)
			created, err := repository.Create(t.Context(), &execution)
			if err == nil {
				t.Fatalf("Expected test suite execution create error")
			}
			if created != nil {
				t.Fatalf("Expected created test suite execution to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestTestSuiteExecutionGetErrors(t *testing.T) {
	testTestSuiteExecutionGetOrDeleteErrors(t, "get", func(
		repository corerepository.TestSuiteExecutionRepository,
		ctx context.Context,
		id string,
	) (*corerepository.TestSuiteExecution, error) {
		return repository.Get(ctx, id)
	})
}

func TestTestSuiteExecutionDeleteErrors(t *testing.T) {
	testTestSuiteExecutionGetOrDeleteErrors(t, "delete", func(
		repository corerepository.TestSuiteExecutionRepository,
		ctx context.Context,
		id string,
	) (*corerepository.TestSuiteExecution, error) {
		return repository.Delete(ctx, id)
	})
}

func TestTestSuiteExecutionGetStatusErrors(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuiteExecution

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
				t.Fatalf("Expected testSuiteExecution get status error")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestTestSuiteExecutionSetRunning(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuiteExecution
	executions := []corerepository.TestSuiteExecution{
		createTestSuiteExecution(t, repository, 1),
		createTestSuiteExecution(t, repository, 2),
	}

	for i, offset := range []time.Duration{time.Hour, 2 * time.Hour} {
		t.Run("test suite execution "+string(rune('1'+i)), func(t *testing.T) {
			expected := executions[i]
			expected.StartedAt = expected.StartedAt.Add(offset)
			expected.Status = "running"

			err := repository.SetRunning(
				t.Context(),
				&corerepository.SetRunningParams{
					ID:        expected.ID,
					StartedAt: expected.StartedAt,
				},
			)
			if err != nil {
				t.Fatalf(
					"Unexpected test suite execution set running error: %v",
					err,
				)
			}

			actual, err := repository.Get(t.Context(), expected.ID)
			if err != nil || actual == nil {
				t.Fatalf(
					"Unexpected test suite execution get error: %v",
					err,
				)
			}

			assertEqualTestSuiteExecutions(t, expected, *actual)
		})
	}

	t.Run("nonexistent test suite execution", func(t *testing.T) {
		err := repository.SetRunning(
			t.Context(),
			&corerepository.SetRunningParams{
				ID:        "does-not-exist",
				StartedAt: time.Now().UTC(),
			},
		)
		if err == nil {
			t.Fatalf(
				"Expected SetRunning to return error: %v, but got nil",
				corerepository.ErrNotFound,
			)
		}
		if !errors.Is(err, corerepository.ErrNotFound) {
			t.Fatalf(
				"Expected SetRunning to return error: %v but got: %v",
				corerepository.ErrNotFound,
				err,
			)
		}
	})
}

func TestTestSuiteExecutionSetCompleted(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuiteExecution
	executions := []corerepository.TestSuiteExecution{
		createTestSuiteExecution(t, repository, 1),
		createTestSuiteExecution(t, repository, 2),
		createTestSuiteExecution(t, repository, 3),
		createTestSuiteExecution(t, repository, 4),
	}

	testCases := []struct {
		name          string
		index         int
		offset        time.Duration
		status        string
		summary       []byte
		expectedError error
	}{
		{
			name:    "passed test suite execution",
			index:   0,
			offset:  time.Hour,
			status:  "passed",
			summary: []byte(`{"result":"passed"}`),
		},
		{
			name:    "failed test suite execution",
			index:   1,
			offset:  2 * time.Hour,
			status:  "failed",
			summary: []byte(`{"result":"failed"}`),
		},
		{
			name:    "errored test suite execution",
			index:   2,
			offset:  3 * time.Hour,
			status:  "error",
			summary: []byte(`{"result":"error"}`),
		},
		{
			name:          "test suite execution invalid status",
			index:         3,
			offset:        4 * time.Hour,
			status:        "invalid",
			summary:       []byte(`{"result":"invalid"}`),
			expectedError: corerepository.ErrInvalidPayload,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			expected := executions[testCase.index]

			err := repository.SetCompleted(
				t.Context(),
				&corerepository.SetCompletedParams{
					ID:         expected.ID,
					FinishedAt: expected.FinishedAt.Add(testCase.offset),
					Status:     testCase.status,
					Summary:    testCase.summary,
				},
			)

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
			} else if err != nil {
				t.Fatalf(
					"Unexpected test suite execution set completed error: %v",
					err,
				)
			}

			if testCase.expectedError == nil {
				expected.FinishedAt = expected.FinishedAt.Add(testCase.offset)
				expected.Status = testCase.status
				expected.Summary = testCase.summary
			}

			actual, err := repository.Get(t.Context(), expected.ID)
			if err != nil {
				t.Fatalf(
					"Unexpected test suite execution get error: %v",
					err,
				)
			}
			if actual == nil {
				t.Fatalf("Expected test suite execution to be non-nil")
			}

			assertEqualTestSuiteExecutions(t, expected, *actual)
		})
	}

	t.Run("nonexistent test suite execution", func(t *testing.T) {
		err := repository.SetCompleted(
			t.Context(),
			&corerepository.SetCompletedParams{
				ID:         "does-not-exist",
				FinishedAt: time.Now().UTC(),
				Status:     "passed",
				Summary:    []byte(`{"result":"passed"}`),
			},
		)
		if err == nil {
			t.Fatalf(
				"Expected SetCompleted to return error: %v, but got nil",
				corerepository.ErrNotFound,
			)
		}
		if !errors.Is(err, corerepository.ErrNotFound) {
			t.Fatalf(
				"Expected SetCompleted to return error: %v but got: %v",
				corerepository.ErrNotFound,
				err,
			)
		}
	})
}

func TestTestSuiteExecutionGetIDsByID(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuiteExecution
	execution1 := createTestSuiteExecution(t, repository, 1)
	execution2 := createTestSuiteExecution(t, repository, 2)
	testCases := []struct {
		name          string
		id            string
		context       func() context.Context
		expected      []string
		expectedError error
	}{
		{
			name:     "empty ID",
			id:       "",
			context:  func() context.Context { return t.Context() },
			expected: []string{execution1.ID, execution2.ID},
		},
		{
			name: "query error",
			id:   execution1.ID,
			context: func() context.Context {
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				return ctx
			},
			expectedError: context.Canceled,
		},
		{
			name:     "non-existing ID",
			id:       "does-not-exist",
			context:  func() context.Context { return t.Context() },
			expected: []string{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ids, err := repository.GetIDsByID(testCase.context(), testCase.id)
			if testCase.expectedError != nil {
				if err == nil || !errors.Is(err, testCase.expectedError) {
					t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
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
			if !slices.Equal(ids, testCase.expected) {
				t.Fatalf("Expected IDs %v, got %v", testCase.expected, ids)
			}
		})
	}
}

func createTestSuiteExecution(t *testing.T, repository corerepository.TestSuiteExecutionRepository, index int) corerepository.TestSuiteExecution {
	t.Helper()
	execution := getDefaultTestSuiteExecution(t, index)
	created, err := repository.Create(t.Context(), &execution)
	if err != nil {
		t.Fatalf("Unexpected test suite execution%d create error: %v", index, err)
	}
	if created == nil {
		t.Fatalf("Expected created%d test suite execution to be non-nil", index)
	}
	assertEqualTestSuiteExecutions(t, execution, *created)
	return execution
}

func testTestSuiteExecutionGetOrDeleteErrors(
	t *testing.T,
	operation string,
	call func(
		corerepository.TestSuiteExecutionRepository,
		context.Context, string,
	) (*corerepository.TestSuiteExecution, error),
) {
	t.Helper()
	repository := getRepositories(t).TestSuiteExecution
	for _, testCase := range []struct {
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
	} {
		t.Run(testCase.name, func(t *testing.T) {
			execution, err := call(repository, t.Context(), testCase.id)
			if err == nil {
				t.Fatalf("Expected test suite execution %s error", operation)
			}
			if execution != nil {
				t.Fatalf("Expected test suite execution to be nil")
			}
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("Expected error %v, got %v", testCase.expectedError, err)
			}
		})
	}
}

func TestTestSuiteExecutionCreateDuplicateID(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.TestSuiteExecution

	execution := getDefaultTestSuiteExecution(t, 1)

	created, err := repository.Create(t.Context(), &execution)
	if err != nil {
		t.Fatalf("Unexpected test suite execution create error: %v", err)
	}
	if created == nil {
		t.Fatalf("Expected created test suite execution to be non-nil")
	}

	duplicate := getDefaultTestSuiteExecution(t, 2)
	duplicate.ID = execution.ID

	created, err = repository.Create(t.Context(), &duplicate)
	if err == nil {
		t.Fatalf("Expected duplicate test suite execution create error")
	}
	if created != nil {
		t.Fatalf("Expected duplicate test suite execution to be nil")
	}
	if !errors.Is(err, corerepository.ErrAlreadyExists) {
		t.Fatalf(
			"Expected error %v, got %v",
			corerepository.ErrAlreadyExists,
			err,
		)
	}
}
