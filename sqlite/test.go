package sqlite

import (
	"context"
	"encoding/json"

	"github.com/e2engine/core/repository"

	db "github.com/e2engine/repository/sqlite/internal/sqlc"
)

var _ repository.TestRepository = (*testRepository)(nil)

type testRepository struct {
	*resourceRepository[repository.Test]
	store db.Store
}

func newTestRepository(
	store db.Store,
) *testRepository {
	return &testRepository{
		resourceRepository: newResourceRepository[repository.Test](
			store,
			resourceKindTest,
		),
		store: store,
	}
}

func (r *testRepository) Create(
	ctx context.Context,
	params *repository.Test,
) (*repository.Test, error) {
	var created db.Resource

	err := r.store.ExecTx(
		ctx,
		func(q *db.Queries) error {
			var err error

			created, err = q.CreateResource(
				ctx,
				db.CreateResourceParams{
					ID:          params.ID,
					Kind:        r.kind,
					Version:     params.Version,
					Name:        params.Name,
					Description: params.Description,
					Spec:        params.Spec,
					CreatedAt:   params.CreatedAt,
					UpdatedAt:   params.UpdatedAt,
				},
			)
			if err != nil {
				return err
			}

			tags, err := getTags(params.Spec)
			if err != nil {
				return err
			}

			for _, tag := range tags {
				if err := q.CreateTestTag(
					ctx,
					db.CreateTestTagParams{
						TestID: params.ID,
						Tag:    tag,
					},
				); err != nil {
					return err
				}
			}

			return nil
		},
	)
	if err != nil {
		return nil, mapError(err)
	}

	result := resourceFromDB[repository.Test](&created)

	return &result, nil
}

func (r *testRepository) GetByTag(
	ctx context.Context,
	tag string,
) ([]repository.Test, error) {
	tests, err := r.store.GetTestsByTag(ctx, tag)
	if err != nil {
		return nil, err
	}

	return resourcesFromDB[repository.Test](tests), nil
}

func getTags(spec json.RawMessage) ([]string, error) {
	type TestSpec struct {
		Tags []string `json:"tags,omitempty"`
	}

	var testSpec TestSpec
	if err := json.Unmarshal(spec, &testSpec); err != nil {
		return nil, err
	}

	return testSpec.Tags, nil
}
