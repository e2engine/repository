package tests

import (
	"testing"
	"time"

	corerepository "github.com/e2engine/core/repository"
)

type testExecutionOrderByCase struct {
	name             string
	orderBy          corerepository.OrderBy
	ascendingIndexes []int
	setAnchorValue   func(*corerepository.ListParams, corerepository.TestExecution)
}

func TestTestExecutionListFirstPage(t *testing.T) {
	repository, testExecutions := createTestExecutionListFixtures(t)

	testCases := []struct {
		name            string
		orderBy         corerepository.OrderBy
		orderDirection  corerepository.OrderDirection
		expectedIndexes []int
	}{
		{
			name:            "ID ascending",
			orderBy:         corerepository.OrderByID,
			orderDirection:  corerepository.OrderDirectionAsc,
			expectedIndexes: []int{0, 1},
		},
		{
			name:            "ID descending",
			orderBy:         corerepository.OrderByID,
			orderDirection:  corerepository.OrderDirectionDesc,
			expectedIndexes: []int{4, 3},
		},
		{
			name:            "started at ascending",
			orderBy:         corerepository.OrderByStartedAt,
			orderDirection:  corerepository.OrderDirectionAsc,
			expectedIndexes: []int{2, 0},
		},
		{
			name:            "started at descending",
			orderBy:         corerepository.OrderByStartedAt,
			orderDirection:  corerepository.OrderDirectionDesc,
			expectedIndexes: []int{4, 1},
		},
		{
			name:            "finished at ascending",
			orderBy:         corerepository.OrderByFinishedAt,
			orderDirection:  corerepository.OrderDirectionAsc,
			expectedIndexes: []int{1, 3},
		},
		{
			name:            "finished at descending",
			orderBy:         corerepository.OrderByFinishedAt,
			orderDirection:  corerepository.OrderDirectionDesc,
			expectedIndexes: []int{0, 2},
		},
		{
			name:            "status ascending",
			orderBy:         corerepository.OrderByStatus,
			orderDirection:  corerepository.OrderDirectionAsc,
			expectedIndexes: []int{2, 1},
		},
		{
			name:            "status descending",
			orderBy:         corerepository.OrderByStatus,
			orderDirection:  corerepository.OrderDirectionDesc,
			expectedIndexes: []int{4, 0},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := repository.List(
				t.Context(),
				&corerepository.ListParams{
					Limit:          2,
					Position:       corerepository.PositionAfter,
					OrderBy:        testCase.orderBy,
					OrderDirection: testCase.orderDirection,
				},
			)
			if err != nil {
				t.Fatalf("Unexpected test executions list error: %v", err)
			}

			assertEqualTestExecutionLists(
				t,
				testExecutionsAt(testExecutions, testCase.expectedIndexes...),
				actual,
			)
		})
	}
}

func TestTestExecutionListCursorPage(t *testing.T) {
	repository, testExecutions := createTestExecutionListFixtures(t)
	traversalCases := []struct {
		name           string
		position       corerepository.Position
		orderDirection corerepository.OrderDirection
		anchorPosition int
		expectedStart  int
	}{
		{
			name:           "after ascending",
			position:       corerepository.PositionAfter,
			orderDirection: corerepository.OrderDirectionAsc,
			anchorPosition: 1,
			expectedStart:  2,
		},
		{
			name:           "after descending",
			position:       corerepository.PositionAfter,
			orderDirection: corerepository.OrderDirectionDesc,
			anchorPosition: 1,
			expectedStart:  2,
		},
		{
			name:           "before ascending",
			position:       corerepository.PositionBefore,
			orderDirection: corerepository.OrderDirectionAsc,
			anchorPosition: 3,
			expectedStart:  1,
		},
		{
			name:           "before descending",
			position:       corerepository.PositionBefore,
			orderDirection: corerepository.OrderDirectionDesc,
			anchorPosition: 3,
			expectedStart:  1,
		},
	}

	for _, orderByCase := range testExecutionListOrderByCases() {
		for _, traversalCase := range traversalCases {
			t.Run(orderByCase.name+" "+traversalCase.name, func(t *testing.T) {
				orderedIndexes := append([]int(nil), orderByCase.ascendingIndexes...)
				if traversalCase.orderDirection == corerepository.OrderDirectionDesc {
					reverseIndexes(orderedIndexes)
				}

				anchor := testExecutions[orderedIndexes[traversalCase.anchorPosition]]
				params := &corerepository.ListParams{
					Limit:          2,
					ID:             anchor.ID,
					Position:       traversalCase.position,
					OrderBy:        orderByCase.orderBy,
					OrderDirection: traversalCase.orderDirection,
				}
				orderByCase.setAnchorValue(params, anchor)

				actual, err := repository.List(t.Context(), params)
				if err != nil {
					t.Fatalf("Unexpected test executions list error: %v", err)
				}

				expectedIndexes := append(
					[]int(nil),
					orderedIndexes[traversalCase.expectedStart:traversalCase.expectedStart+2]...,
				)
				if traversalCase.position == corerepository.PositionBefore {
					reverseIndexes(expectedIndexes)
				}

				assertEqualTestExecutionLists(
					t,
					testExecutionsAt(testExecutions, expectedIndexes...),
					actual,
				)
			})
		}
	}
}

func TestTestExecutionListForwardTraversal(t *testing.T) {
	repository, testExecutions := createTestExecutionListFixtures(t)
	directions := []corerepository.OrderDirection{
		corerepository.OrderDirectionAsc,
		corerepository.OrderDirectionDesc,
	}

	for _, orderByCase := range testExecutionListOrderByCases() {
		for _, direction := range directions {
			t.Run(orderByCase.name+" "+string(direction), func(t *testing.T) {
				expectedIndexes := append([]int(nil), orderByCase.ascendingIndexes...)
				if direction == corerepository.OrderDirectionDesc {
					reverseIndexes(expectedIndexes)
				}

				params := &corerepository.ListParams{
					Limit:          2,
					Position:       corerepository.PositionAfter,
					OrderBy:        orderByCase.orderBy,
					OrderDirection: direction,
				}
				var actual []corerepository.TestExecution

				for pageNumber := 0; ; pageNumber++ {
					if pageNumber > len(testExecutions) {
						t.Fatalf(
							"Expected traversal to finish, got more than %d pages",
							len(testExecutions)+1,
						)
					}

					page, err := repository.List(t.Context(), params)
					if err != nil {
						t.Fatalf("Unexpected test executions list error: %v", err)
					}
					if len(page) == 0 {
						break
					}

					actual = append(actual, page...)
					anchor := page[len(page)-1]
					params.ID = anchor.ID
					orderByCase.setAnchorValue(params, anchor)
				}

				assertEqualTestExecutionLists(
					t,
					testExecutionsAt(testExecutions, expectedIndexes...),
					actual,
				)
			})
		}
	}
}

func TestTestExecutionListBackwardTraversal(t *testing.T) {
	repository, testExecutions := createTestExecutionListFixtures(t)
	directions := []corerepository.OrderDirection{
		corerepository.OrderDirectionAsc,
		corerepository.OrderDirectionDesc,
	}

	for _, orderByCase := range testExecutionListOrderByCases() {
		for _, direction := range directions {
			t.Run(orderByCase.name+" "+string(direction), func(t *testing.T) {
				expectedIndexes := append([]int(nil), orderByCase.ascendingIndexes...)
				if direction == corerepository.OrderDirectionDesc {
					reverseIndexes(expectedIndexes)
				}

				params := &corerepository.ListParams{
					Limit:          2,
					Position:       corerepository.PositionBefore,
					OrderBy:        orderByCase.orderBy,
					OrderDirection: direction,
				}
				last := testExecutions[expectedIndexes[len(expectedIndexes)-1]]
				params.ID = last.ID
				orderByCase.setAnchorValue(params, last)

				actual := []corerepository.TestExecution{last}

				for pageNumber := 0; ; pageNumber++ {
					if pageNumber > len(testExecutions) {
						t.Fatalf(
							"Expected traversal to finish, got more than %d pages",
							len(testExecutions)+1,
						)
					}

					page, err := repository.List(t.Context(), params)
					if err != nil {
						t.Fatalf("Unexpected test executions list error: %v", err)
					}
					if len(page) == 0 {
						break
					}

					reverseTestExecutions(page)
					actual = append(page, actual...)
					anchor := page[0]
					params.ID = anchor.ID
					orderByCase.setAnchorValue(params, anchor)
				}

				assertEqualTestExecutionLists(
					t,
					testExecutionsAt(testExecutions, expectedIndexes...),
					actual,
				)
			})
		}
	}
}

func TestTestExecutionListCursorEqualSortValue(t *testing.T) {
	testCases := []struct {
		name            string
		orderBy         corerepository.OrderBy
		prepare         func([]corerepository.TestExecution)
		setAnchorValue  func(*corerepository.ListParams, corerepository.TestExecution)
		expectedIndexes map[corerepository.OrderDirection][]int
	}{
		{
			name:    "started at",
			orderBy: corerepository.OrderByStartedAt,
			prepare: func(testExecutions []corerepository.TestExecution) {
				base := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
				testExecutions[0].StartedAt = base.Add(-time.Hour)
				testExecutions[1].StartedAt = base
				testExecutions[2].StartedAt = base
				testExecutions[3].StartedAt = base
				testExecutions[4].StartedAt = base.Add(time.Hour)
			},
			setAnchorValue: func(params *corerepository.ListParams, testExecution corerepository.TestExecution) {
				params.ValueTimestamp = testExecution.StartedAt
			},
			expectedIndexes: map[corerepository.OrderDirection][]int{
				corerepository.OrderDirectionAsc:  {3, 4},
				corerepository.OrderDirectionDesc: {1, 0},
			},
		},
		{
			name:    "finished at",
			orderBy: corerepository.OrderByFinishedAt,
			prepare: func(testExecutions []corerepository.TestExecution) {
				base := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
				testExecutions[0].FinishedAt = base.Add(-time.Hour)
				testExecutions[1].FinishedAt = base
				testExecutions[2].FinishedAt = base
				testExecutions[3].FinishedAt = base
				testExecutions[4].FinishedAt = base.Add(time.Hour)
			},
			setAnchorValue: func(params *corerepository.ListParams, testExecution corerepository.TestExecution) {
				params.ValueTimestamp = testExecution.FinishedAt
			},
			expectedIndexes: map[corerepository.OrderDirection][]int{
				corerepository.OrderDirectionAsc:  {3, 4},
				corerepository.OrderDirectionDesc: {1, 0},
			},
		},
		{
			name:    "status",
			orderBy: corerepository.OrderByStatus,
			prepare: func(testExecutions []corerepository.TestExecution) {
				testExecutions[0].Status = "error"
				testExecutions[1].Status = "failed"
				testExecutions[2].Status = "failed"
				testExecutions[3].Status = "failed"
				testExecutions[4].Status = "running"
			},
			setAnchorValue: func(params *corerepository.ListParams, testExecution corerepository.TestExecution) {
				params.ValueString = testExecution.Status
			},
			expectedIndexes: map[corerepository.OrderDirection][]int{
				corerepository.OrderDirectionAsc:  {3, 4},
				corerepository.OrderDirectionDesc: {1, 0},
			},
		},
	}
	directions := []corerepository.OrderDirection{
		corerepository.OrderDirectionAsc,
		corerepository.OrderDirectionDesc,
	}

	for _, testCase := range testCases {
		for _, direction := range directions {
			t.Run(testCase.name+" "+string(direction), func(t *testing.T) {
				testExecutions := getTestExecutionsForList(
					t,
					idPrefixTestSuiteExecution+idBody+"1",
				)
				testCase.prepare(testExecutions)
				repository := createTestExecutionListRepository(t, testExecutions)

				anchor := testExecutions[2]
				params := &corerepository.ListParams{
					Limit:          2,
					ID:             anchor.ID,
					Position:       corerepository.PositionAfter,
					OrderBy:        testCase.orderBy,
					OrderDirection: direction,
				}
				testCase.setAnchorValue(params, anchor)

				actual, err := repository.List(t.Context(), params)
				if err != nil {
					t.Fatalf("Unexpected test executions list error: %v", err)
				}

				assertEqualTestExecutionLists(
					t,
					testExecutionsAt(testExecutions, testCase.expectedIndexes[direction]...),
					actual,
				)
			})
		}
	}
}

func TestTestExecutionListBoundaryCursorPositions(t *testing.T) {
	repository, testExecutions := createTestExecutionListFixtures(t)
	boundaryCases := []struct {
		name           string
		position       corerepository.Position
		anchorPosition func([]int) int
	}{
		{
			name:     "after last",
			position: corerepository.PositionAfter,
			anchorPosition: func(indexes []int) int {
				return len(indexes) - 1
			},
		},
		{
			name:     "before first",
			position: corerepository.PositionBefore,
			anchorPosition: func([]int) int {
				return 0
			},
		},
	}
	directions := []corerepository.OrderDirection{
		corerepository.OrderDirectionAsc,
		corerepository.OrderDirectionDesc,
	}

	for _, orderByCase := range testExecutionListOrderByCases() {
		for _, direction := range directions {
			for _, boundaryCase := range boundaryCases {
				t.Run(
					orderByCase.name+" "+string(direction)+" "+boundaryCase.name,
					func(t *testing.T) {
						orderedIndexes := append([]int(nil), orderByCase.ascendingIndexes...)
						if direction == corerepository.OrderDirectionDesc {
							reverseIndexes(orderedIndexes)
						}

						anchor := testExecutions[orderedIndexes[boundaryCase.anchorPosition(orderedIndexes)]]
						params := &corerepository.ListParams{
							Limit:          2,
							ID:             anchor.ID,
							Position:       boundaryCase.position,
							OrderBy:        orderByCase.orderBy,
							OrderDirection: direction,
						}
						orderByCase.setAnchorValue(params, anchor)

						actual, err := repository.List(t.Context(), params)
						if err != nil {
							t.Fatalf("Unexpected test executions list error: %v", err)
						}
						if len(actual) != 0 {
							t.Fatalf(
								"Expected empty test executions list, got %d test executions",
								len(actual),
							)
						}
					},
				)
			}
		}
	}
}

func TestTestExecutionListLimitBoundaries(t *testing.T) {
	repository, testExecutions := createTestExecutionListFixtures(t)

	testCases := []struct {
		name          string
		limit         int
		expectedCount int
		expectError   string
	}{
		{
			name:          "limit one",
			limit:         1,
			expectedCount: 1,
		},
		{
			name:          "limit equal to test executions count",
			limit:         len(testExecutions),
			expectedCount: len(testExecutions),
		},
		{
			name:          "limit greater than test executions count",
			limit:         len(testExecutions) + 1,
			expectedCount: len(testExecutions),
		},
		{
			name:        "zero limit",
			limit:       0,
			expectError: "invalid page limit",
		},
		{
			name:        "negative limit",
			limit:       -1,
			expectError: "invalid page limit",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := repository.List(
				t.Context(),
				&corerepository.ListParams{
					Limit:          testCase.limit,
					Position:       corerepository.PositionAfter,
					OrderBy:        corerepository.OrderByID,
					OrderDirection: corerepository.OrderDirectionAsc,
				},
			)

			if testCase.expectError != "" {
				if err == nil {
					t.Fatalf("Expected test executions list error for limit %d", testCase.limit)
				}
				if actual != nil {
					t.Fatalf(
						"Expected test executions to be nil, got %d test executions",
						len(actual),
					)
				}
				if err.Error() != testCase.expectError {
					t.Fatalf(
						"Expected test executions list error %q, got %q",
						testCase.expectError,
						err.Error(),
					)
				}
				return
			}

			assertEqualTestExecutionLists(
				t,
				testExecutionsAt(
					testExecutions,
					makeSequentialIndexes(testCase.expectedCount)...,
				),
				actual,
			)
		})
	}
}

func TestTestExecutionListInvalidTraversal(t *testing.T) {
	repository, testExecutions := createTestExecutionListFixtures(t)

	testCases := []struct {
		name   string
		mutate func(*corerepository.ListParams)
	}{
		{
			name: "invalid order by",
			mutate: func(params *corerepository.ListParams) {
				params.OrderBy = corerepository.OrderBy("invalid")
			},
		},
		{
			name: "invalid position",
			mutate: func(params *corerepository.ListParams) {
				params.Position = corerepository.Position("invalid")
			},
		},
		{
			name: "invalid order direction",
			mutate: func(params *corerepository.ListParams) {
				params.OrderDirection = corerepository.OrderDirection("invalid")
			},
		},
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			params := &corerepository.ListParams{
				Limit:          2,
				ID:             testExecutions[2].ID,
				ValueString:    testExecutions[2].ID,
				Position:       corerepository.PositionAfter,
				OrderBy:        corerepository.OrderByID,
				OrderDirection: corerepository.OrderDirectionAsc,
			}
			testCase.mutate(params)

			actual, err := repository.List(t.Context(), params)
			if err == nil {
				t.Fatalf("Expected test executions list error")
			}
			if actual != nil {
				t.Fatalf(
					"Expected test executions to be nil, got %d test executions",
					len(actual),
				)
			}
			if err.Error() != "invalid page traversal" {
				t.Fatalf(
					"Expected test executions list error %q, got %q",
					"invalid page traversal",
					err.Error(),
				)
			}
		})
	}
}

func TestTestExecutionListNonexistentCursorAnchor(t *testing.T) {
	repository, testExecutions := createTestExecutionListFixtures(t)

	firstPage, err := repository.List(
		t.Context(),
		&corerepository.ListParams{
			Limit:          2,
			Position:       corerepository.PositionAfter,
			OrderBy:        corerepository.OrderByID,
			OrderDirection: corerepository.OrderDirectionAsc,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected first test executions page error: %v", err)
	}
	assertEqualTestExecutionLists(t, testExecutionsAt(testExecutions, 0, 1), firstPage)

	anchor := firstPage[len(firstPage)-1]
	deleted, err := repository.Delete(t.Context(), anchor.ID)
	if err != nil {
		t.Fatalf("Unexpected cursor anchor delete error: %v", err)
	}
	if deleted == nil {
		t.Fatalf("Expected deleted cursor anchor to be non-nil")
	}
	assertEqualTestExecutions(t, anchor, *deleted)

	secondPage, err := repository.List(
		t.Context(),
		&corerepository.ListParams{
			Limit:          2,
			ID:             anchor.ID,
			Position:       corerepository.PositionAfter,
			OrderBy:        corerepository.OrderByID,
			OrderDirection: corerepository.OrderDirectionAsc,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected second test executions page error: %v", err)
	}
	assertEqualTestExecutionLists(t, testExecutionsAt(testExecutions, 2, 3), secondPage)
}

func createTestExecutionListFixtures(
	t *testing.T,
) (corerepository.TestExecutionRepository, []corerepository.TestExecution) {
	t.Helper()

	testExecutions := getTestExecutionsForList(
		t,
		idPrefixTestSuiteExecution+idBody+"1",
	)
	testExecutions[0].Status = "passed"
	testExecutions[1].Status = "failed"
	testExecutions[2].Status = "error"
	testExecutions[3].Status = "failed"
	testExecutions[4].Status = "running"

	return createTestExecutionListRepository(t, testExecutions), testExecutions
}

func createTestExecutionListRepository(
	t *testing.T,
	testExecutions []corerepository.TestExecution,
) corerepository.TestExecutionRepository {
	t.Helper()

	repositories := getRepositories(t)
	testSuiteExecution := getDefaultTestSuiteExecution(t, 1)
	createdSuite, err := repositories.TestSuiteExecution.Create(t.Context(), &testSuiteExecution)
	if err != nil {
		t.Fatalf("Unexpected test suite execution create error: %v", err)
	}
	if createdSuite == nil {
		t.Fatalf("Expected created test suite execution to be non-nil")
	}

	createTestExecutions(t, repositories.TestExecution, testExecutions)

	return repositories.TestExecution
}

func testExecutionListOrderByCases() []testExecutionOrderByCase {
	return []testExecutionOrderByCase{
		{
			name:             "ID",
			orderBy:          corerepository.OrderByID,
			ascendingIndexes: []int{0, 1, 2, 3, 4},
			setAnchorValue: func(params *corerepository.ListParams, testExecution corerepository.TestExecution) {
				params.ValueString = testExecution.ID
			},
		},
		{
			name:             "started at",
			orderBy:          corerepository.OrderByStartedAt,
			ascendingIndexes: []int{2, 0, 3, 1, 4},
			setAnchorValue: func(params *corerepository.ListParams, testExecution corerepository.TestExecution) {
				params.ValueTimestamp = testExecution.StartedAt
			},
		},
		{
			name:             "finished at",
			orderBy:          corerepository.OrderByFinishedAt,
			ascendingIndexes: []int{1, 3, 4, 2, 0},
			setAnchorValue: func(params *corerepository.ListParams, testExecution corerepository.TestExecution) {
				params.ValueTimestamp = testExecution.FinishedAt
			},
		},
		{
			name:             "status",
			orderBy:          corerepository.OrderByStatus,
			ascendingIndexes: []int{2, 1, 3, 0, 4},
			setAnchorValue: func(params *corerepository.ListParams, testExecution corerepository.TestExecution) {
				params.ValueString = testExecution.Status
			},
		},
	}
}

func testExecutionsAt(
	testExecutions []corerepository.TestExecution,
	indexes ...int,
) []corerepository.TestExecution {
	result := make([]corerepository.TestExecution, len(indexes))

	for i, index := range indexes {
		result[i] = testExecutions[index]
	}

	return result
}

func assertEqualTestExecutionLists(
	t *testing.T,
	expected, actual []corerepository.TestExecution,
) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Fatalf("Expected %d test executions, got %d", len(expected), len(actual))
	}

	for i := range expected {
		assertEqualTestExecutions(t, expected[i], actual[i])
	}
}

func reverseIndexes(indexes []int) {
	for left, right := 0, len(indexes)-1; left < right; left, right = left+1, right-1 {
		indexes[left], indexes[right] = indexes[right], indexes[left]
	}
}

func reverseTestExecutions(testExecutions []corerepository.TestExecution) {
	for left, right := 0, len(testExecutions)-1; left < right; left, right = left+1, right-1 {
		testExecutions[left], testExecutions[right] = testExecutions[right], testExecutions[left]
	}
}
