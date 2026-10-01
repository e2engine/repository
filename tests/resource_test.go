package tests

import (
	"errors"
	"testing"

	corerepository "github.com/e2engine/core/repository"
)

func TestResourceKindIsolation(t *testing.T) {
	testCases := []struct {
		name  string
		kind  string
		index int
	}{
		{
			name:  "environment",
			kind:  "environment",
			index: 0,
		},
		{
			name:  "test",
			kind:  "test",
			index: 1,
		},
		{
			name:  "test suite",
			kind:  "test_suite",
			index: 2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repositories := getRepositories(t)

			resources := getResourcesForList(t, "shared")[:3]

			resources[0].Kind = "environment"
			resources[0].ID = "shared-environment"
			resources[0].Name = "shared-environment"

			resources[1].Kind = "test"
			resources[1].ID = "shared-test"
			resources[1].Name = "shared-test"

			resources[2].Kind = "test_suite"
			resources[2].ID = "shared-test-suite"
			resources[2].Name = "shared-test-suite"

			environment0 := corerepository.Environment(resources[0])
			if _, err := repositories.Environment.Create(
				t.Context(),
				&environment0,
			); err != nil {
				t.Fatalf("Unexpected environment create error: %v", err)
			}

			test1 := corerepository.Test(resources[1])
			if _, err := repositories.Test.Create(
				t.Context(),
				&test1,
			); err != nil {
				t.Fatalf("Unexpected test create error: %v", err)
			}

			testSuite2 := corerepository.TestSuite(resources[2])
			if _, err := repositories.TestSuite.Create(
				t.Context(),
				&testSuite2,
			); err != nil {
				t.Fatalf("Unexpected test suite create error: %v", err)
			}

			var (
				get func(string) (corerepository.Resource, error)

				list func() ([]corerepository.Resource, error)

				getIDsByID   func(string) ([]string, error)
				getIDsByName func(string) ([]string, error)

				deleteResource func(string) (corerepository.Resource, error)
			)

			switch testCase.kind {
			case "environment":
				get = func(id string) (corerepository.Resource, error) {
					resource, err := repositories.Environment.Get(
						t.Context(),
						id,
					)
					if err != nil {
						return corerepository.Resource{}, err
					}

					return corerepository.Resource(*resource), nil
				}

				list = func() ([]corerepository.Resource, error) {
					values, err := repositories.Environment.List(
						t.Context(),
						&corerepository.ListParams{
							Limit:          10,
							Position:       corerepository.PositionAfter,
							OrderBy:        corerepository.OrderByID,
							OrderDirection: corerepository.OrderDirectionAsc,
						},
					)
					if err != nil {
						return nil, err
					}

					result := make([]corerepository.Resource, len(values))
					for i := range values {
						result[i] = corerepository.Resource(values[i])
					}

					return result, nil
				}

				getIDsByID = func(prefix string) ([]string, error) {
					return repositories.Environment.GetIDsByID(
						t.Context(),
						prefix,
					)
				}

				getIDsByName = func(prefix string) ([]string, error) {
					return repositories.Environment.GetIDsByName(
						t.Context(),
						prefix,
					)
				}

				deleteResource = func(id string) (corerepository.Resource, error) {
					resource, err := repositories.Environment.Delete(
						t.Context(),
						id,
					)
					if err != nil {
						return corerepository.Resource{}, err
					}

					return corerepository.Resource(*resource), nil
				}

			case "test":
				get = func(id string) (corerepository.Resource, error) {
					resource, err := repositories.Test.Get(
						t.Context(),
						id,
					)
					if err != nil {
						return corerepository.Resource{}, err
					}

					return corerepository.Resource(*resource), nil
				}

				list = func() ([]corerepository.Resource, error) {
					values, err := repositories.Test.List(
						t.Context(),
						&corerepository.ListParams{
							Limit:          10,
							Position:       corerepository.PositionAfter,
							OrderBy:        corerepository.OrderByID,
							OrderDirection: corerepository.OrderDirectionAsc,
						},
					)
					if err != nil {
						return nil, err
					}

					result := make([]corerepository.Resource, len(values))
					for i := range values {
						result[i] = corerepository.Resource(values[i])
					}

					return result, nil
				}

				getIDsByID = func(prefix string) ([]string, error) {
					return repositories.Test.GetIDsByID(
						t.Context(),
						prefix,
					)
				}

				getIDsByName = func(prefix string) ([]string, error) {
					return repositories.Test.GetIDsByName(
						t.Context(),
						prefix,
					)
				}

				deleteResource = func(id string) (corerepository.Resource, error) {
					resource, err := repositories.Test.Delete(
						t.Context(),
						id,
					)
					if err != nil {
						return corerepository.Resource{}, err
					}

					return corerepository.Resource(*resource), nil
				}

			case "test_suite":
				get = func(id string) (corerepository.Resource, error) {
					resource, err := repositories.TestSuite.Get(
						t.Context(),
						id,
					)
					if err != nil {
						return corerepository.Resource{}, err
					}

					return corerepository.Resource(*resource), nil
				}

				list = func() ([]corerepository.Resource, error) {
					values, err := repositories.TestSuite.List(
						t.Context(),
						&corerepository.ListParams{
							Limit:          10,
							Position:       corerepository.PositionAfter,
							OrderBy:        corerepository.OrderByID,
							OrderDirection: corerepository.OrderDirectionAsc,
						},
					)
					if err != nil {
						return nil, err
					}

					result := make([]corerepository.Resource, len(values))
					for i := range values {
						result[i] = corerepository.Resource(values[i])
					}

					return result, nil
				}

				getIDsByID = func(prefix string) ([]string, error) {
					return repositories.TestSuite.GetIDsByID(
						t.Context(),
						prefix,
					)
				}

				getIDsByName = func(prefix string) ([]string, error) {
					return repositories.TestSuite.GetIDsByName(
						t.Context(),
						prefix,
					)
				}

				deleteResource = func(id string) (corerepository.Resource, error) {
					resource, err := repositories.TestSuite.Delete(
						t.Context(),
						id,
					)
					if err != nil {
						return corerepository.Resource{}, err
					}

					return corerepository.Resource(*resource), nil
				}
			}

			expected := resources[testCase.index]

			actual, err := get(expected.ID)
			if err != nil {
				t.Fatalf("Unexpected get error: %v", err)
			}

			assertEqualResources(t, expected, actual)

			// Other kinds must not be visible through Get.
			for i := range resources {
				if i == testCase.index {
					continue
				}

				_, err := get(resources[i].ID)
				if !errors.Is(err, corerepository.ErrNotFound) {
					t.Fatalf(
						"Expected resource %q to be isolated, got error: %v",
						resources[i].ID,
						err,
					)
				}
			}

			// List must contain only resources of this kind.
			listed, err := list()
			if err != nil {
				t.Fatalf("Unexpected list error: %v", err)
			}

			if len(listed) != 1 {
				t.Fatalf(
					"Expected 1 resource, got %d",
					len(listed),
				)
			}

			assertEqualResources(t, expected, listed[0])

			// Prefix matches resources from all three kinds,
			// but lookup must return only the current kind.
			ids, err := getIDsByID("shared-")
			if err != nil {
				t.Fatalf("Unexpected GetIDsByID error: %v", err)
			}

			if len(ids) != 1 || ids[0] != expected.ID {
				t.Fatalf(
					"Expected IDs [%s], got %v",
					expected.ID,
					ids,
				)
			}

			ids, err = getIDsByName("shared-")
			if err != nil {
				t.Fatalf("Unexpected GetIDsByName error: %v", err)
			}

			if len(ids) != 1 || ids[0] != expected.ID {
				t.Fatalf(
					"Expected IDs [%s], got %v",
					expected.ID,
					ids,
				)
			}

			deleted, err := deleteResource(expected.ID)
			if err != nil {
				t.Fatalf("Unexpected delete error: %v", err)
			}

			assertEqualResources(t, expected, deleted)

			// Deleting one kind must not touch the same table rows
			// belonging to the other kinds.
			for i := range resources {
				if i == testCase.index {
					continue
				}

				var err error

				switch i {
				case 0:
					_, err = repositories.Environment.Get(
						t.Context(),
						resources[i].ID,
					)

				case 1:
					_, err = repositories.Test.Get(
						t.Context(),
						resources[i].ID,
					)

				case 2:
					_, err = repositories.TestSuite.Get(
						t.Context(),
						resources[i].ID,
					)
				}

				if err != nil {
					t.Fatalf(
						"Deleting %s unexpectedly affected resource %q: %v",
						testCase.name,
						resources[i].ID,
						err,
					)
				}
			}
		})
	}
}
