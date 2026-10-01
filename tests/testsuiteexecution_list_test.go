package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	corerepository "github.com/e2engine/core/repository"
)

type testSuiteExecutionOrderByCase struct {
	name             string
	orderBy          corerepository.OrderBy
	ascendingIndexes []int
	setAnchorValue   func(*corerepository.ListParams, corerepository.TestSuiteExecution)
}

func TestTestSuiteExecutionListFirstPage(t *testing.T) {
	repository, executions := createTestSuiteExecutionListFixtures(t)
	testCases := []struct {
		name      string
		orderBy   corerepository.OrderBy
		direction corerepository.OrderDirection
		expected  []int
	}{
		{"ID ascending", corerepository.OrderByID, corerepository.OrderDirectionAsc, []int{0, 1}},
		{"ID descending", corerepository.OrderByID, corerepository.OrderDirectionDesc, []int{4, 3}},
		{"started at ascending", corerepository.OrderByStartedAt, corerepository.OrderDirectionAsc, []int{2, 0}},
		{"started at descending", corerepository.OrderByStartedAt, corerepository.OrderDirectionDesc, []int{4, 1}},
		{"finished at ascending", corerepository.OrderByFinishedAt, corerepository.OrderDirectionAsc, []int{1, 3}},
		{"finished at descending", corerepository.OrderByFinishedAt, corerepository.OrderDirectionDesc, []int{0, 2}},
		{"status ascending", corerepository.OrderByStatus, corerepository.OrderDirectionAsc, []int{2, 1}},
		{"status descending", corerepository.OrderByStatus, corerepository.OrderDirectionDesc, []int{4, 0}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := listTestSuiteExecutions(t, repository, &corerepository.ListParams{
				Limit: 2, Position: corerepository.PositionAfter, OrderBy: testCase.orderBy, OrderDirection: testCase.direction,
			})
			assertEqualTestSuiteExecutionLists(t, testSuiteExecutionsAt(executions, testCase.expected...), actual)
		})
	}
}

func TestTestSuiteExecutionListCursorPage(t *testing.T) {
	repository, executions := createTestSuiteExecutionListFixtures(t)
	traversals := []struct {
		name                          string
		position                      corerepository.Position
		direction                     corerepository.OrderDirection
		anchorPosition, expectedStart int
	}{
		{"after ascending", corerepository.PositionAfter, corerepository.OrderDirectionAsc, 1, 2},
		{"after descending", corerepository.PositionAfter, corerepository.OrderDirectionDesc, 1, 2},
		{"before ascending", corerepository.PositionBefore, corerepository.OrderDirectionAsc, 3, 1},
		{"before descending", corerepository.PositionBefore, corerepository.OrderDirectionDesc, 3, 1},
	}

	for _, orderCase := range testSuiteExecutionListOrderByCases() {
		for _, traversal := range traversals {
			t.Run(orderCase.name+" "+traversal.name, func(t *testing.T) {
				ordered := append([]int(nil), orderCase.ascendingIndexes...)
				if traversal.direction == corerepository.OrderDirectionDesc {
					reverseIndexes(ordered)
				}
				anchor := executions[ordered[traversal.anchorPosition]]
				params := &corerepository.ListParams{
					Limit: 2, ID: anchor.ID, Position: traversal.position,
					OrderBy: orderCase.orderBy, OrderDirection: traversal.direction,
				}
				orderCase.setAnchorValue(params, anchor)
				expected := append([]int(nil), ordered[traversal.expectedStart:traversal.expectedStart+2]...)
				if traversal.position == corerepository.PositionBefore {
					reverseIndexes(expected)
				}
				assertEqualTestSuiteExecutionLists(t, testSuiteExecutionsAt(executions, expected...), listTestSuiteExecutions(t, repository, params))
			})
		}
	}
}

func TestTestSuiteExecutionListForwardTraversal(t *testing.T) {
	repository, executions := createTestSuiteExecutionListFixtures(t)
	for _, orderCase := range testSuiteExecutionListOrderByCases() {
		for _, direction := range []corerepository.OrderDirection{corerepository.OrderDirectionAsc, corerepository.OrderDirectionDesc} {
			t.Run(orderCase.name+" "+string(direction), func(t *testing.T) {
				expected := append([]int(nil), orderCase.ascendingIndexes...)
				if direction == corerepository.OrderDirectionDesc {
					reverseIndexes(expected)
				}
				params := &corerepository.ListParams{Limit: 2, Position: corerepository.PositionAfter, OrderBy: orderCase.orderBy, OrderDirection: direction}
				var actual []corerepository.TestSuiteExecution
				for pageNumber := 0; ; pageNumber++ {
					if pageNumber > len(executions) {
						t.Fatalf("Expected traversal to finish, got more than %d pages", len(executions)+1)
					}
					page := listTestSuiteExecutions(t, repository, params)
					if len(page) == 0 {
						break
					}
					actual = append(actual, page...)
					anchor := page[len(page)-1]
					params.ID = anchor.ID
					orderCase.setAnchorValue(params, anchor)
				}
				assertEqualTestSuiteExecutionLists(t, testSuiteExecutionsAt(executions, expected...), actual)
			})
		}
	}
}

func TestTestSuiteExecutionListBackwardTraversal(t *testing.T) {
	repository, executions := createTestSuiteExecutionListFixtures(t)
	for _, orderCase := range testSuiteExecutionListOrderByCases() {
		for _, direction := range []corerepository.OrderDirection{corerepository.OrderDirectionAsc, corerepository.OrderDirectionDesc} {
			t.Run(orderCase.name+" "+string(direction), func(t *testing.T) {
				expected := append([]int(nil), orderCase.ascendingIndexes...)
				if direction == corerepository.OrderDirectionDesc {
					reverseIndexes(expected)
				}
				last := executions[expected[len(expected)-1]]
				params := &corerepository.ListParams{Limit: 2, ID: last.ID, Position: corerepository.PositionBefore, OrderBy: orderCase.orderBy, OrderDirection: direction}
				orderCase.setAnchorValue(params, last)
				actual := []corerepository.TestSuiteExecution{last}
				for pageNumber := 0; ; pageNumber++ {
					if pageNumber > len(executions) {
						t.Fatalf("Expected traversal to finish, got more than %d pages", len(executions)+1)
					}
					page := listTestSuiteExecutions(t, repository, params)
					if len(page) == 0 {
						break
					}
					reverseTestSuiteExecutions(page)
					actual = append(page, actual...)
					anchor := page[0]
					params.ID = anchor.ID
					orderCase.setAnchorValue(params, anchor)
				}
				assertEqualTestSuiteExecutionLists(t, testSuiteExecutionsAt(executions, expected...), actual)
			})
		}
	}
}

func TestTestSuiteExecutionListCursorEqualSortValue(t *testing.T) {
	testCases := []struct {
		name           string
		orderBy        corerepository.OrderBy
		prepare        func([]corerepository.TestSuiteExecution)
		setAnchorValue func(*corerepository.ListParams, corerepository.TestSuiteExecution)
	}{
		{
			"started at", corerepository.OrderByStartedAt,
			func(executions []corerepository.TestSuiteExecution) {
				base := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
				executions[0].StartedAt, executions[1].StartedAt, executions[2].StartedAt, executions[3].StartedAt, executions[4].StartedAt = base.Add(-time.Hour), base, base, base, base.Add(time.Hour)
			},
			func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
				params.ValueTimestamp = execution.StartedAt
			},
		},
		{
			"finished at", corerepository.OrderByFinishedAt,
			func(executions []corerepository.TestSuiteExecution) {
				base := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
				executions[0].FinishedAt, executions[1].FinishedAt, executions[2].FinishedAt, executions[3].FinishedAt, executions[4].FinishedAt = base.Add(-time.Hour), base, base, base, base.Add(time.Hour)
			},
			func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
				params.ValueTimestamp = execution.FinishedAt
			},
		},
		{
			"status", corerepository.OrderByStatus,
			func(executions []corerepository.TestSuiteExecution) {
				executions[0].Status, executions[1].Status, executions[2].Status, executions[3].Status, executions[4].Status = "error", "failed", "failed", "failed", "running"
			},
			func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
				params.ValueString = execution.Status
			},
		},
	}

	for _, testCase := range testCases {
		for _, direction := range []corerepository.OrderDirection{corerepository.OrderDirectionAsc, corerepository.OrderDirectionDesc} {
			t.Run(testCase.name+" "+string(direction), func(t *testing.T) {
				executions := getTestSuiteExecutionsForList(t)
				testCase.prepare(executions)
				repository := createTestSuiteExecutionListRepository(t, executions)
				anchor := executions[2]
				params := &corerepository.ListParams{Limit: 2, ID: anchor.ID, Position: corerepository.PositionAfter, OrderBy: testCase.orderBy, OrderDirection: direction}
				testCase.setAnchorValue(params, anchor)
				expected := []int{3, 4}
				if direction == corerepository.OrderDirectionDesc {
					expected = []int{1, 0}
				}
				assertEqualTestSuiteExecutionLists(t, testSuiteExecutionsAt(executions, expected...), listTestSuiteExecutions(t, repository, params))
			})
		}
	}
}

func TestTestSuiteExecutionListBoundaryCursorPositions(t *testing.T) {
	repository, executions := createTestSuiteExecutionListFixtures(t)
	for _, orderCase := range testSuiteExecutionListOrderByCases() {
		for _, direction := range []corerepository.OrderDirection{corerepository.OrderDirectionAsc, corerepository.OrderDirectionDesc} {
			for _, boundary := range []struct {
				name     string
				position corerepository.Position
				index    func([]int) int
			}{
				{"after last", corerepository.PositionAfter, func(indexes []int) int { return len(indexes) - 1 }},
				{"before first", corerepository.PositionBefore, func([]int) int { return 0 }},
			} {
				t.Run(orderCase.name+" "+string(direction)+" "+boundary.name, func(t *testing.T) {
					ordered := append([]int(nil), orderCase.ascendingIndexes...)
					if direction == corerepository.OrderDirectionDesc {
						reverseIndexes(ordered)
					}
					anchor := executions[ordered[boundary.index(ordered)]]
					params := &corerepository.ListParams{Limit: 2, ID: anchor.ID, Position: boundary.position, OrderBy: orderCase.orderBy, OrderDirection: direction}
					orderCase.setAnchorValue(params, anchor)
					actual := listTestSuiteExecutions(t, repository, params)
					if len(actual) != 0 {
						t.Fatalf("Expected empty test suite executions list, got %d test suite executions", len(actual))
					}
				})
			}
		}
	}
}

func TestTestSuiteExecutionListLimitBoundaries(t *testing.T) {
	repository, executions := createTestSuiteExecutionListFixtures(t)
	for _, testCase := range []struct {
		name                 string
		limit, expectedCount int
		expectedError        string
	}{
		{"limit one", 1, 1, ""},
		{"limit equal to test suite executions count", len(executions), len(executions), ""},
		{"limit greater than test suite executions count", len(executions) + 1, len(executions), ""},
		{"zero limit", 0, 0, "invalid page limit"},
		{"negative limit", -1, 0, "invalid page limit"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := repository.List(t.Context(), &corerepository.ListParams{Limit: testCase.limit, Position: corerepository.PositionAfter, OrderBy: corerepository.OrderByID, OrderDirection: corerepository.OrderDirectionAsc})
			if testCase.expectedError != "" {
				if err == nil || err.Error() != testCase.expectedError {
					t.Fatalf("Expected test suite executions list error %q, got %v", testCase.expectedError, err)
				}
				if actual != nil {
					t.Fatalf("Expected test suite executions to be nil, got %d test suite executions", len(actual))
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected test suite executions list error: %v", err)
			}
			assertEqualTestSuiteExecutionLists(t, testSuiteExecutionsAt(executions, makeSequentialIndexes(testCase.expectedCount)...), actual)
		})
	}
}

func TestTestSuiteExecutionListInvalidTraversal(t *testing.T) {
	repository, executions := createTestSuiteExecutionListFixtures(t)
	for _, testCase := range []struct {
		name   string
		mutate func(*corerepository.ListParams)
	}{
		{"invalid order by", func(params *corerepository.ListParams) { params.OrderBy = corerepository.OrderBy("invalid") }},
		{"invalid position", func(params *corerepository.ListParams) { params.Position = corerepository.Position("invalid") }},
		{"invalid order direction", func(params *corerepository.ListParams) {
			params.OrderDirection = corerepository.OrderDirection("invalid")
		}},
		{
			name: "first page invalid order by",
			mutate: func(params *corerepository.ListParams) {
				params.ID = ""
				params.OrderBy = corerepository.OrderBy("invalid")
			},
		},
		{
			name: "first page invalid order direction",
			mutate: func(params *corerepository.ListParams) {
				params.ID = ""
				params.OrderDirection = corerepository.OrderDirection("invalid")
			},
		},
		{
			name: "invalid position order by id",
			mutate: func(params *corerepository.ListParams) {
				params.OrderBy = corerepository.OrderByID
				params.Position = corerepository.Position("invalid")
			},
		},
		{
			name: "invalid position order by status",
			mutate: func(params *corerepository.ListParams) {
				params.OrderBy = corerepository.OrderByStatus
				params.Position = corerepository.Position("invalid")
			},
		},
		{
			name: "invalid position order by started at",
			mutate: func(params *corerepository.ListParams) {
				params.OrderBy = corerepository.OrderByStartedAt
				params.Position = corerepository.Position("invalid")
			},
		},
		{
			name: "invalid position order by finished at",
			mutate: func(params *corerepository.ListParams) {
				params.OrderBy = corerepository.OrderByFinishedAt
				params.Position = corerepository.Position("invalid")
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			params := &corerepository.ListParams{Limit: 2, ID: executions[2].ID, ValueString: executions[2].ID, Position: corerepository.PositionAfter, OrderBy: corerepository.OrderByID, OrderDirection: corerepository.OrderDirectionAsc}
			testCase.mutate(params)
			actual, err := repository.List(t.Context(), params)
			if err == nil || err.Error() != "invalid page traversal" {
				t.Fatalf("Expected test suite executions list error %q, got %v", "invalid page traversal", err)
			}
			if actual != nil {
				t.Fatalf("Expected test suite executions to be nil, got %d test suite executions", len(actual))
			}
		})
	}
}

func TestTestSuiteExecutionListNonexistentCursorAnchor(t *testing.T) {
	repository, executions := createTestSuiteExecutionListFixtures(t)
	firstPage := listTestSuiteExecutions(t, repository, &corerepository.ListParams{Limit: 2, Position: corerepository.PositionAfter, OrderBy: corerepository.OrderByID, OrderDirection: corerepository.OrderDirectionAsc})
	assertEqualTestSuiteExecutionLists(t, testSuiteExecutionsAt(executions, 0, 1), firstPage)
	anchor := firstPage[len(firstPage)-1]
	deleted, err := repository.Delete(t.Context(), anchor.ID)
	if err != nil {
		t.Fatalf("Unexpected cursor anchor delete error: %v", err)
	}
	if deleted == nil {
		t.Fatalf("Expected deleted cursor anchor to be non-nil")
	}
	assertEqualTestSuiteExecutions(t, anchor, *deleted)
	secondPage := listTestSuiteExecutions(t, repository, &corerepository.ListParams{Limit: 2, ID: anchor.ID, Position: corerepository.PositionAfter, OrderBy: corerepository.OrderByID, OrderDirection: corerepository.OrderDirectionAsc})
	assertEqualTestSuiteExecutionLists(t, testSuiteExecutionsAt(executions, 2, 3), secondPage)
}

func TestTestSuiteExecutionListQueryErrors(t *testing.T) {
	repository, executions := createTestSuiteExecutionListFixtures(t)

	t.Run("first page", func(t *testing.T) {
		testCases := []struct {
			name      string
			orderBy   corerepository.OrderBy
			direction corerepository.OrderDirection
		}{
			{"ID ascending", corerepository.OrderByID, corerepository.OrderDirectionAsc},
			{"ID descending", corerepository.OrderByID, corerepository.OrderDirectionDesc},
			{"started at ascending", corerepository.OrderByStartedAt, corerepository.OrderDirectionAsc},
			{"started at descending", corerepository.OrderByStartedAt, corerepository.OrderDirectionDesc},
			{"finished at ascending", corerepository.OrderByFinishedAt, corerepository.OrderDirectionAsc},
			{"finished at descending", corerepository.OrderByFinishedAt, corerepository.OrderDirectionDesc},
			{"status ascending", corerepository.OrderByStatus, corerepository.OrderDirectionAsc},
			{"status descending", corerepository.OrderByStatus, corerepository.OrderDirectionDesc},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				ctx, cancel := context.WithCancel(t.Context())
				cancel()

				actual, err := repository.List(ctx, &corerepository.ListParams{
					Limit:          2,
					Position:       corerepository.PositionAfter,
					OrderBy:        testCase.orderBy,
					OrderDirection: testCase.direction,
				})
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("Expected context canceled error, got %v", err)
				}
				if actual != nil {
					t.Fatalf(
						"Expected test suite executions to be nil, got %d test suite executions",
						len(actual),
					)
				}
			})
		}
	})

	t.Run("cursor page", func(t *testing.T) {
		testCases := []struct {
			name           string
			orderBy        corerepository.OrderBy
			position       corerepository.Position
			direction      corerepository.OrderDirection
			setAnchorValue func(*corerepository.ListParams, corerepository.TestSuiteExecution)
		}{
			{
				"ID after ascending",
				corerepository.OrderByID,
				corerepository.PositionAfter,
				corerepository.OrderDirectionAsc,
				func(*corerepository.ListParams, corerepository.TestSuiteExecution) {},
			},
			{
				"ID after descending",
				corerepository.OrderByID,
				corerepository.PositionAfter,
				corerepository.OrderDirectionDesc,
				func(*corerepository.ListParams, corerepository.TestSuiteExecution) {},
			},
			{
				"ID before ascending",
				corerepository.OrderByID,
				corerepository.PositionBefore,
				corerepository.OrderDirectionAsc,
				func(*corerepository.ListParams, corerepository.TestSuiteExecution) {},
			},
			{
				"ID before descending",
				corerepository.OrderByID,
				corerepository.PositionBefore,
				corerepository.OrderDirectionDesc,
				func(*corerepository.ListParams, corerepository.TestSuiteExecution) {},
			},

			{
				"started at after ascending",
				corerepository.OrderByStartedAt,
				corerepository.PositionAfter,
				corerepository.OrderDirectionAsc,
				func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
					params.ValueTimestamp = execution.StartedAt
				},
			},
			{
				"started at after descending",
				corerepository.OrderByStartedAt,
				corerepository.PositionAfter,
				corerepository.OrderDirectionDesc,
				func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
					params.ValueTimestamp = execution.StartedAt
				},
			},
			{
				"started at before ascending",
				corerepository.OrderByStartedAt,
				corerepository.PositionBefore,
				corerepository.OrderDirectionAsc,
				func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
					params.ValueTimestamp = execution.StartedAt
				},
			},
			{
				"started at before descending",
				corerepository.OrderByStartedAt,
				corerepository.PositionBefore,
				corerepository.OrderDirectionDesc,
				func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
					params.ValueTimestamp = execution.StartedAt
				},
			},

			{
				"finished at after ascending",
				corerepository.OrderByFinishedAt,
				corerepository.PositionAfter,
				corerepository.OrderDirectionAsc,
				func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
					params.ValueTimestamp = execution.FinishedAt
				},
			},
			{
				"finished at after descending",
				corerepository.OrderByFinishedAt,
				corerepository.PositionAfter,
				corerepository.OrderDirectionDesc,
				func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
					params.ValueTimestamp = execution.FinishedAt
				},
			},
			{
				"finished at before ascending",
				corerepository.OrderByFinishedAt,
				corerepository.PositionBefore,
				corerepository.OrderDirectionAsc,
				func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
					params.ValueTimestamp = execution.FinishedAt
				},
			},
			{
				"finished at before descending",
				corerepository.OrderByFinishedAt,
				corerepository.PositionBefore,
				corerepository.OrderDirectionDesc,
				func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
					params.ValueTimestamp = execution.FinishedAt
				},
			},

			{
				"status after ascending",
				corerepository.OrderByStatus,
				corerepository.PositionAfter,
				corerepository.OrderDirectionAsc,
				func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
					params.ValueString = execution.Status
				},
			},
			{
				"status after descending",
				corerepository.OrderByStatus,
				corerepository.PositionAfter,
				corerepository.OrderDirectionDesc,
				func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
					params.ValueString = execution.Status
				},
			},
			{
				"status before ascending",
				corerepository.OrderByStatus,
				corerepository.PositionBefore,
				corerepository.OrderDirectionAsc,
				func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
					params.ValueString = execution.Status
				},
			},
			{
				"status before descending",
				corerepository.OrderByStatus,
				corerepository.PositionBefore,
				corerepository.OrderDirectionDesc,
				func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
					params.ValueString = execution.Status
				},
			},
		}

		anchor := executions[2]

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				params := &corerepository.ListParams{
					Limit:          2,
					ID:             anchor.ID,
					Position:       testCase.position,
					OrderBy:        testCase.orderBy,
					OrderDirection: testCase.direction,
				}
				testCase.setAnchorValue(params, anchor)

				ctx, cancel := context.WithCancel(t.Context())
				cancel()

				actual, err := repository.List(ctx, params)
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("Expected context canceled error, got %v", err)
				}
				if actual != nil {
					t.Fatalf(
						"Expected test suite executions to be nil, got %d test suite executions",
						len(actual),
					)
				}
			})
		}
	})
}

func createTestSuiteExecutionListFixtures(t *testing.T) (corerepository.TestSuiteExecutionRepository, []corerepository.TestSuiteExecution) {
	t.Helper()
	executions := getTestSuiteExecutionsForList(t)
	return createTestSuiteExecutionListRepository(t, executions), executions
}

func createTestSuiteExecutionListRepository(t *testing.T, executions []corerepository.TestSuiteExecution) corerepository.TestSuiteExecutionRepository {
	t.Helper()
	repositories := getRepositories(t)
	for i := range executions {
		execution := executions[i]
		created, err := repositories.TestSuiteExecution.Create(t.Context(), &execution)
		if err != nil {
			t.Fatalf("Unexpected test suite execution%d create error: %v", i+1, err)
		}
		if created == nil {
			t.Fatalf("Expected created%d test suite execution to be non-nil", i+1)
		}
		assertEqualTestSuiteExecutions(t, execution, *created)
	}
	return repositories.TestSuiteExecution
}

func getTestSuiteExecutionsForList(t *testing.T) []corerepository.TestSuiteExecution {
	t.Helper()
	startedAt := []time.Duration{2 * time.Second, 3 * time.Second, time.Second, 2 * time.Second, 3 * time.Second}
	finishedAt := []time.Duration{4 * time.Second, time.Second, 3 * time.Second, time.Second, 2 * time.Second}
	base := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	executions := make([]corerepository.TestSuiteExecution, len(startedAt))
	statuses := []string{"passed", "failed", "error", "failed", "running"}
	for i := range executions {
		executions[i] = getTestSuiteExecution(t, i+1, base.Add(startedAt[i]), base.Add(finishedAt[i]), statuses[i], nil)
	}
	return executions
}

func testSuiteExecutionListOrderByCases() []testSuiteExecutionOrderByCase {
	return []testSuiteExecutionOrderByCase{
		{"ID", corerepository.OrderByID, []int{0, 1, 2, 3, 4}, func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
			params.ValueString = execution.ID
		}},
		{"started at", corerepository.OrderByStartedAt, []int{2, 0, 3, 1, 4}, func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
			params.ValueTimestamp = execution.StartedAt
		}},
		{"finished at", corerepository.OrderByFinishedAt, []int{1, 3, 4, 2, 0}, func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
			params.ValueTimestamp = execution.FinishedAt
		}},
		{"status", corerepository.OrderByStatus, []int{2, 1, 3, 0, 4}, func(params *corerepository.ListParams, execution corerepository.TestSuiteExecution) {
			params.ValueString = execution.Status
		}},
	}
}

func listTestSuiteExecutions(t *testing.T, repository corerepository.TestSuiteExecutionRepository, params *corerepository.ListParams) []corerepository.TestSuiteExecution {
	t.Helper()
	executions, err := repository.List(t.Context(), params)
	if err != nil {
		t.Fatalf("Unexpected test suite executions list error: %v", err)
	}
	return executions
}

func testSuiteExecutionsAt(executions []corerepository.TestSuiteExecution, indexes ...int) []corerepository.TestSuiteExecution {
	result := make([]corerepository.TestSuiteExecution, len(indexes))
	for i, index := range indexes {
		result[i] = executions[index]
	}
	return result
}

func assertEqualTestSuiteExecutionLists(t *testing.T, expected, actual []corerepository.TestSuiteExecution) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Fatalf("Expected %d test suite executions, got %d", len(expected), len(actual))
	}
	for i := range expected {
		assertEqualTestSuiteExecutions(t, expected[i], actual[i])
	}
}

func reverseTestSuiteExecutions(executions []corerepository.TestSuiteExecution) {
	for left, right := 0, len(executions)-1; left < right; left, right = left+1, right-1 {
		executions[left], executions[right] = executions[right], executions[left]
	}
}
