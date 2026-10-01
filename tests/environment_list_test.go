package tests

import (
	"testing"
	"time"

	corerepository "github.com/e2engine/core/repository"
)

func TestEnvironmentListFirstPage(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment
	resources := getResourcesForList(t, "environment")
	createEnvironmentResources(t, repository, resources)

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
			name:            "name ascending",
			orderBy:         corerepository.OrderByName,
			orderDirection:  corerepository.OrderDirectionAsc,
			expectedIndexes: []int{1, 3},
		},
		{
			name:            "name descending",
			orderBy:         corerepository.OrderByName,
			orderDirection:  corerepository.OrderDirectionDesc,
			expectedIndexes: []int{2, 4},
		},
		{
			name:            "version ascending",
			orderBy:         corerepository.OrderByVersion,
			orderDirection:  corerepository.OrderDirectionAsc,
			expectedIndexes: []int{1, 3},
		},
		{
			name:            "version descending",
			orderBy:         corerepository.OrderByVersion,
			orderDirection:  corerepository.OrderDirectionDesc,
			expectedIndexes: []int{2, 4},
		},
		{
			name:            "created at ascending",
			orderBy:         corerepository.OrderByCreatedAt,
			orderDirection:  corerepository.OrderDirectionAsc,
			expectedIndexes: []int{2, 0},
		},
		{
			name:            "created at descending",
			orderBy:         corerepository.OrderByCreatedAt,
			orderDirection:  corerepository.OrderDirectionDesc,
			expectedIndexes: []int{4, 1},
		},
		{
			name:            "updated at ascending",
			orderBy:         corerepository.OrderByUpdatedAt,
			orderDirection:  corerepository.OrderDirectionAsc,
			expectedIndexes: []int{1, 3},
		},
		{
			name:            "updated at descending",
			orderBy:         corerepository.OrderByUpdatedAt,
			orderDirection:  corerepository.OrderDirectionDesc,
			expectedIndexes: []int{0, 2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			environments, err := repository.List(
				t.Context(),
				&corerepository.ListParams{
					Limit:          2,
					Position:       corerepository.PositionAfter,
					OrderBy:        testCase.orderBy,
					OrderDirection: testCase.orderDirection,
				},
			)
			if err != nil {
				t.Fatalf("Unexpected environments list error: %v", err)
			}

			assertEqualEnvironments(
				t,
				environmentsAt(resources, testCase.expectedIndexes...),
				environments,
			)
		})
	}
}

func TestEnvironmentListCursorPage(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment
	resources := getResourcesForList(t, "environment")
	createEnvironmentResources(t, repository, resources)

	orderByCases := []struct {
		name             string
		orderBy          corerepository.OrderBy
		ascendingIndexes []int
		setAnchorValue   func(*corerepository.ListParams, corerepository.Resource)
	}{
		{
			name:             "ID",
			orderBy:          corerepository.OrderByID,
			ascendingIndexes: []int{0, 1, 2, 3, 4},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueString = resource.ID
			},
		},
		{
			name:             "name",
			orderBy:          corerepository.OrderByName,
			ascendingIndexes: []int{1, 3, 0, 4, 2},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueString = resource.Name
			},
		},
		{
			name:             "version",
			orderBy:          corerepository.OrderByVersion,
			ascendingIndexes: []int{1, 3, 0, 4, 2},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueString = resource.Version
			},
		},
		{
			name:             "created at",
			orderBy:          corerepository.OrderByCreatedAt,
			ascendingIndexes: []int{2, 0, 3, 1, 4},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueTimestamp = resource.CreatedAt
			},
		},
		{
			name:             "updated at",
			orderBy:          corerepository.OrderByUpdatedAt,
			ascendingIndexes: []int{1, 3, 4, 2, 0},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueTimestamp = resource.UpdatedAt
			},
		},
	}
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

	for _, orderByCase := range orderByCases {
		for _, traversalCase := range traversalCases {
			t.Run(orderByCase.name+" "+traversalCase.name, func(t *testing.T) {
				orderedIndexes := append([]int(nil), orderByCase.ascendingIndexes...)
				if traversalCase.orderDirection == corerepository.OrderDirectionDesc {
					for left, right := 0, len(orderedIndexes)-1; left < right; left, right = left+1, right-1 {
						orderedIndexes[left], orderedIndexes[right] = orderedIndexes[right], orderedIndexes[left]
					}
				}

				anchor := resources[orderedIndexes[traversalCase.anchorPosition]]
				params := &corerepository.ListParams{
					Limit:          2,
					ID:             anchor.ID,
					Position:       traversalCase.position,
					OrderBy:        orderByCase.orderBy,
					OrderDirection: traversalCase.orderDirection,
				}
				orderByCase.setAnchorValue(params, anchor)

				environments, err := repository.List(t.Context(), params)
				if err != nil {
					t.Fatalf("Unexpected environments list error: %v", err)
				}

				expectedIndexes := append(
					[]int(nil),
					orderedIndexes[traversalCase.expectedStart:traversalCase.expectedStart+2]...,
				)
				if traversalCase.position == corerepository.PositionBefore {
					for left, right := 0, len(expectedIndexes)-1; left < right; left, right = left+1, right-1 {
						expectedIndexes[left], expectedIndexes[right] = expectedIndexes[right], expectedIndexes[left]
					}
				}

				assertEqualEnvironments(
					t,
					environmentsAt(resources, expectedIndexes...),
					environments,
				)
			})
		}
	}
}

func TestEnvironmentListForwardTraversal(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment
	resources := getResourcesForList(t, "environment")
	createEnvironmentResources(t, repository, resources)

	orderByCases := []struct {
		name             string
		orderBy          corerepository.OrderBy
		ascendingIndexes []int
		setAnchorValue   func(*corerepository.ListParams, corerepository.Resource)
	}{
		{
			name:             "ID",
			orderBy:          corerepository.OrderByID,
			ascendingIndexes: []int{0, 1, 2, 3, 4},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueString = resource.ID
			},
		},
		{
			name:             "name",
			orderBy:          corerepository.OrderByName,
			ascendingIndexes: []int{1, 3, 0, 4, 2},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueString = resource.Name
			},
		},
		{
			name:             "version",
			orderBy:          corerepository.OrderByVersion,
			ascendingIndexes: []int{1, 3, 0, 4, 2},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueString = resource.Version
			},
		},
		{
			name:             "created at",
			orderBy:          corerepository.OrderByCreatedAt,
			ascendingIndexes: []int{2, 0, 3, 1, 4},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueTimestamp = resource.CreatedAt
			},
		},
		{
			name:             "updated at",
			orderBy:          corerepository.OrderByUpdatedAt,
			ascendingIndexes: []int{1, 3, 4, 2, 0},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueTimestamp = resource.UpdatedAt
			},
		},
	}

	directions := []corerepository.OrderDirection{
		corerepository.OrderDirectionAsc,
		corerepository.OrderDirectionDesc,
	}

	for _, orderByCase := range orderByCases {
		for _, direction := range directions {
			t.Run(orderByCase.name+" "+string(direction), func(t *testing.T) {
				expectedIndexes := append([]int(nil), orderByCase.ascendingIndexes...)
				if direction == corerepository.OrderDirectionDesc {
					for left, right := 0, len(expectedIndexes)-1; left < right; left, right = left+1, right-1 {
						expectedIndexes[left], expectedIndexes[right] = expectedIndexes[right], expectedIndexes[left]
					}
				}

				params := &corerepository.ListParams{
					Limit:          2,
					Position:       corerepository.PositionAfter,
					OrderBy:        orderByCase.orderBy,
					OrderDirection: direction,
				}

				var actual []corerepository.Environment

				for pageNumber := 0; ; pageNumber++ {
					if pageNumber > len(resources) {
						t.Fatalf("Expected traversal to finish, got more than %d pages", len(resources)+1)
					}

					page, err := repository.List(t.Context(), params)
					if err != nil {
						t.Fatalf("Unexpected environments list error: %v", err)
					}
					if len(page) == 0 {
						break
					}

					actual = append(actual, page...)

					anchor := corerepository.Resource(page[len(page)-1])
					params.ID = anchor.ID
					orderByCase.setAnchorValue(params, anchor)
				}

				assertEqualEnvironments(
					t,
					environmentsAt(resources, expectedIndexes...),
					actual,
				)
			})
		}
	}
}

func TestEnvironmentListBackwardTraversal(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment
	resources := getResourcesForList(t, "environment")
	createEnvironmentResources(t, repository, resources)

	orderByCases := []struct {
		name             string
		orderBy          corerepository.OrderBy
		ascendingIndexes []int
		setAnchorValue   func(*corerepository.ListParams, corerepository.Resource)
	}{
		{
			name:             "ID",
			orderBy:          corerepository.OrderByID,
			ascendingIndexes: []int{0, 1, 2, 3, 4},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueString = resource.ID
			},
		},
		{
			name:             "name",
			orderBy:          corerepository.OrderByName,
			ascendingIndexes: []int{1, 3, 0, 4, 2},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueString = resource.Name
			},
		},
		{
			name:             "version",
			orderBy:          corerepository.OrderByVersion,
			ascendingIndexes: []int{1, 3, 0, 4, 2},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueString = resource.Version
			},
		},
		{
			name:             "created at",
			orderBy:          corerepository.OrderByCreatedAt,
			ascendingIndexes: []int{2, 0, 3, 1, 4},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueTimestamp = resource.CreatedAt
			},
		},
		{
			name:             "updated at",
			orderBy:          corerepository.OrderByUpdatedAt,
			ascendingIndexes: []int{1, 3, 4, 2, 0},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueTimestamp = resource.UpdatedAt
			},
		},
	}

	directions := []corerepository.OrderDirection{
		corerepository.OrderDirectionAsc,
		corerepository.OrderDirectionDesc,
	}

	for _, orderByCase := range orderByCases {
		for _, direction := range directions {
			t.Run(orderByCase.name+" "+string(direction), func(t *testing.T) {
				expectedIndexes := append([]int(nil), orderByCase.ascendingIndexes...)
				if direction == corerepository.OrderDirectionDesc {
					for left, right := 0, len(expectedIndexes)-1; left < right; left, right = left+1, right-1 {
						expectedIndexes[left], expectedIndexes[right] = expectedIndexes[right], expectedIndexes[left]
					}
				}

				params := &corerepository.ListParams{
					Limit:          2,
					Position:       corerepository.PositionBefore,
					OrderBy:        orderByCase.orderBy,
					OrderDirection: direction,
				}

				last := resources[expectedIndexes[len(expectedIndexes)-1]]
				params.ID = last.ID
				orderByCase.setAnchorValue(params, last)

				actual := []corerepository.Environment{
					corerepository.Environment(last),
				}

				for pageNumber := 0; ; pageNumber++ {
					if pageNumber > len(resources) {
						t.Fatalf(
							"Expected traversal to finish, got more than %d pages",
							len(resources)+1,
						)
					}

					page, err := repository.List(t.Context(), params)
					if err != nil {
						t.Fatalf("Unexpected environments list error: %v", err)
					}
					if len(page) == 0 {
						break
					}

					for left, right := 0, len(page)-1; left < right; left, right = left+1, right-1 {
						page[left], page[right] = page[right], page[left]
					}

					actual = append(page, actual...)

					anchor := corerepository.Resource(page[0])
					params.ID = anchor.ID
					orderByCase.setAnchorValue(params, anchor)
				}

				assertEqualEnvironments(
					t,
					environmentsAt(resources, expectedIndexes...),
					actual,
				)
			})
		}
	}
}

func TestEnvironmentListCursorEqualSortValue(t *testing.T) {
	orderByCases := []struct {
		name           string
		orderBy        corerepository.OrderBy
		prepare        func([]corerepository.Resource)
		setAnchorValue func(*corerepository.ListParams, corerepository.Resource)
	}{
		{
			name:    "version",
			orderBy: corerepository.OrderByVersion,
			prepare: func(resources []corerepository.Resource) {
				resources[0].Version = "0.9.0"
				resources[1].Version = "1.0.0"
				resources[2].Version = "1.0.0"
				resources[3].Version = "1.0.0"
				resources[4].Version = "2.0.0"
			},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueString = resource.Version
			},
		},
		{
			name:    "created at",
			orderBy: corerepository.OrderByCreatedAt,
			prepare: func(resources []corerepository.Resource) {
				base := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
				resources[0].CreatedAt = base.Add(-time.Hour)
				resources[1].CreatedAt = base
				resources[2].CreatedAt = base
				resources[3].CreatedAt = base
				resources[4].CreatedAt = base.Add(time.Hour)
			},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueTimestamp = resource.CreatedAt
			},
		},
	}

	directionCases := []struct {
		name            string
		direction       corerepository.OrderDirection
		expectedIndexes []int
	}{
		{
			name:            "ascending",
			direction:       corerepository.OrderDirectionAsc,
			expectedIndexes: []int{3, 4},
		},
		{
			name:            "descending",
			direction:       corerepository.OrderDirectionDesc,
			expectedIndexes: []int{1, 0},
		},
	}

	for _, orderByCase := range orderByCases {
		for _, directionCase := range directionCases {
			t.Run(orderByCase.name+" "+directionCase.name, func(t *testing.T) {
				repositories := getRepositories(t)
				repository := repositories.Environment
				resources := getResourcesForList(t, "environment")
				orderByCase.prepare(resources)
				createEnvironmentResources(t, repository, resources)

				anchor := resources[2]
				params := &corerepository.ListParams{
					Limit:          2,
					ID:             anchor.ID,
					Position:       corerepository.PositionAfter,
					OrderBy:        orderByCase.orderBy,
					OrderDirection: directionCase.direction,
				}
				orderByCase.setAnchorValue(params, anchor)

				environments, err := repository.List(t.Context(), params)
				if err != nil {
					t.Fatalf("Unexpected environments list error: %v", err)
				}

				assertEqualEnvironments(
					t,
					environmentsAt(resources, directionCase.expectedIndexes...),
					environments,
				)
			})
		}
	}
}

func TestEnvironmentListBoundaryCursorPositions(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment
	resources := getResourcesForList(t, "environment")
	createEnvironmentResources(t, repository, resources)

	orderByCases := []struct {
		name             string
		orderBy          corerepository.OrderBy
		ascendingIndexes []int
		setAnchorValue   func(*corerepository.ListParams, corerepository.Resource)
	}{
		{
			name:             "ID",
			orderBy:          corerepository.OrderByID,
			ascendingIndexes: []int{0, 1, 2, 3, 4},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueString = resource.ID
			},
		},
		{
			name:             "name",
			orderBy:          corerepository.OrderByName,
			ascendingIndexes: []int{1, 3, 0, 4, 2},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueString = resource.Name
			},
		},
		{
			name:             "version",
			orderBy:          corerepository.OrderByVersion,
			ascendingIndexes: []int{1, 3, 0, 4, 2},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueString = resource.Version
			},
		},
		{
			name:             "created at",
			orderBy:          corerepository.OrderByCreatedAt,
			ascendingIndexes: []int{2, 0, 3, 1, 4},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueTimestamp = resource.CreatedAt
			},
		},
		{
			name:             "updated at",
			orderBy:          corerepository.OrderByUpdatedAt,
			ascendingIndexes: []int{1, 3, 4, 2, 0},
			setAnchorValue: func(params *corerepository.ListParams, resource corerepository.Resource) {
				params.ValueTimestamp = resource.UpdatedAt
			},
		},
	}

	directionCases := []corerepository.OrderDirection{
		corerepository.OrderDirectionAsc,
		corerepository.OrderDirectionDesc,
	}

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

	for _, orderByCase := range orderByCases {
		for _, direction := range directionCases {
			for _, boundaryCase := range boundaryCases {
				t.Run(
					orderByCase.name+" "+string(direction)+" "+boundaryCase.name,
					func(t *testing.T) {
						orderedIndexes := append([]int(nil), orderByCase.ascendingIndexes...)
						if direction == corerepository.OrderDirectionDesc {
							for left, right := 0, len(orderedIndexes)-1; left < right; left, right = left+1, right-1 {
								orderedIndexes[left], orderedIndexes[right] = orderedIndexes[right], orderedIndexes[left]
							}
						}

						anchorIndex := boundaryCase.anchorPosition(orderedIndexes)
						anchor := resources[orderedIndexes[anchorIndex]]

						params := &corerepository.ListParams{
							Limit:          2,
							ID:             anchor.ID,
							Position:       boundaryCase.position,
							OrderBy:        orderByCase.orderBy,
							OrderDirection: direction,
						}
						orderByCase.setAnchorValue(params, anchor)

						environments, err := repository.List(t.Context(), params)
						if err != nil {
							t.Fatalf("Unexpected environments list error: %v", err)
						}
						if len(environments) != 0 {
							t.Fatalf(
								"Expected empty environments list, got %d environments",
								len(environments),
							)
						}
					},
				)
			}
		}
	}
}

func TestEnvironmentListLimitBoundaries(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment
	resources := getResourcesForList(t, "environment")
	createEnvironmentResources(t, repository, resources)

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
			name:          "limit equal to resources count",
			limit:         len(resources),
			expectedCount: len(resources),
		},
		{
			name:          "limit greater than resources count",
			limit:         len(resources) + 1,
			expectedCount: len(resources),
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
			environments, err := repository.List(
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
					t.Fatalf("Expected environments list error for limit %d", testCase.limit)
				}
				if environments != nil {
					t.Fatalf(
						"Expected environments to be nil, got %d environments",
						len(environments),
					)
				}
				if err.Error() != testCase.expectError {
					t.Fatalf(
						"Expected environments list error %q, got %q",
						testCase.expectError,
						err.Error(),
					)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected environments list error: %v", err)
			}
			if len(environments) != testCase.expectedCount {
				t.Fatalf(
					"Expected %d environments, got %d",
					testCase.expectedCount,
					len(environments),
				)
			}

			assertEqualEnvironments(
				t,
				environmentsAt(
					resources,
					makeSequentialIndexes(testCase.expectedCount)...,
				),
				environments,
			)
		})
	}
}

func TestEnvironmentListInvalidTraversal(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment
	resources := getResourcesForList(t, "environment")
	createEnvironmentResources(t, repository, resources)

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
			name: "invalid position order by name",
			mutate: func(params *corerepository.ListParams) {
				params.OrderBy = corerepository.OrderByName
				params.Position = corerepository.Position("invalid")
			},
		},
		{
			name: "invalid position order by version",
			mutate: func(params *corerepository.ListParams) {
				params.OrderBy = corerepository.OrderByVersion
				params.Position = corerepository.Position("invalid")
			},
		},
		{
			name: "invalid position order by created at",
			mutate: func(params *corerepository.ListParams) {
				params.OrderBy = corerepository.OrderByCreatedAt
				params.Position = corerepository.Position("invalid")
			},
		},
		{
			name: "invalid position order by updated at",
			mutate: func(params *corerepository.ListParams) {
				params.OrderBy = corerepository.OrderByUpdatedAt
				params.Position = corerepository.Position("invalid")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			params := &corerepository.ListParams{
				Limit:          2,
				ID:             resources[2].ID,
				ValueString:    resources[2].ID,
				Position:       corerepository.PositionAfter,
				OrderBy:        corerepository.OrderByID,
				OrderDirection: corerepository.OrderDirectionAsc,
			}
			testCase.mutate(params)

			environments, err := repository.List(t.Context(), params)
			if err == nil {
				t.Fatalf("Expected environments list error")
			}
			if environments != nil {
				t.Fatalf(
					"Expected environments to be nil, got %d environments",
					len(environments),
				)
			}
			if err.Error() != "invalid page traversal" {
				t.Fatalf(
					"Expected environments list error %q, got %q",
					"invalid page traversal",
					err.Error(),
				)
			}
		})
	}
}

func TestEnvironmentListNonexistentCursorAnchor(t *testing.T) {
	repositories := getRepositories(t)
	repository := repositories.Environment
	resources := getResourcesForList(t, "environment")
	createEnvironmentResources(t, repository, resources)

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
		t.Fatalf("Unexpected first environments page error: %v", err)
	}

	assertEqualEnvironments(
		t,
		environmentsAt(resources, 0, 1),
		firstPage,
	)

	anchor := corerepository.Resource(firstPage[len(firstPage)-1])

	deleted, err := repository.Delete(t.Context(), anchor.ID)
	if err != nil {
		t.Fatalf("Unexpected cursor anchor delete error: %v", err)
	}
	if deleted == nil {
		t.Fatalf("Expected deleted cursor anchor to be non-nil")
	}
	assertEqualResources(t, anchor, corerepository.Resource(*deleted))

	secondPage, err := repository.List(
		t.Context(),
		&corerepository.ListParams{
			Limit:          2,
			ID:             anchor.ID,
			ValueString:    anchor.ID,
			Position:       corerepository.PositionAfter,
			OrderBy:        corerepository.OrderByID,
			OrderDirection: corerepository.OrderDirectionAsc,
		},
	)
	if err != nil {
		t.Fatalf("Unexpected second environments page error: %v", err)
	}

	assertEqualEnvironments(
		t,
		environmentsAt(resources, 2, 3),
		secondPage,
	)
}
