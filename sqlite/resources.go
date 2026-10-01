package sqlite

import (
	"context"
	"database/sql"

	"github.com/e2engine/core/repository"

	db "github.com/e2engine/repository/sqlite/internal/sqlc"
)

const (
	resourceKindEnvironment = "environment"
	resourceKindTest        = "test"
	resourceKindTestSuite   = "test_suite"
)

type resource interface {
	repository.Environment |
		repository.Test |
		repository.TestSuite
}

type resourceRepository[E resource] struct {
	queries db.Querier
	kind    string
}

var (
	_ repository.EnvironmentRepository = (*resourceRepository[repository.Environment])(nil)
	_ repository.TestSuiteRepository   = (*resourceRepository[repository.TestSuite])(nil)
)

func newEnvironmentRepository(
	queries db.Querier,
) *resourceRepository[repository.Environment] {
	return newResourceRepository[repository.Environment](
		queries,
		resourceKindEnvironment,
	)
}

func newTestSuiteRepository(
	queries db.Querier,
) *resourceRepository[repository.TestSuite] {
	return newResourceRepository[repository.TestSuite](
		queries,
		resourceKindTestSuite,
	)
}

func newResourceRepository[E resource](
	queries db.Querier,
	kind string,
) *resourceRepository[E] {
	return &resourceRepository[E]{
		queries: queries,
		kind:    kind,
	}
}

func (r *resourceRepository[E]) Create(
	ctx context.Context,
	params *E,
) (*E, error) {
	res := repository.Resource(*params)

	created, err := r.queries.CreateResource(
		ctx,
		db.CreateResourceParams{
			ID:          res.ID,
			Kind:        r.kind,
			Version:     res.Version,
			Name:        res.Name,
			Description: res.Description,
			Spec:        res.Spec,
			CreatedAt:   res.CreatedAt,
			UpdatedAt:   res.UpdatedAt,
		},
	)
	if err != nil {
		return nil, mapError(err)
	}

	result := resourceFromDB[E](&created)

	return &result, nil
}

func (r *resourceRepository[E]) Get(
	ctx context.Context,
	id string,
) (*E, error) {
	res, err := r.queries.GetResource(
		ctx,
		db.GetResourceParams{
			Kind: r.kind,
			ID:   id,
		},
	)
	if err != nil {
		return nil, mapError(err)
	}

	result := resourceFromDB[E](&res)

	return &result, nil
}

func (r *resourceRepository[E]) Delete(
	ctx context.Context,
	id string,
) (*E, error) {
	res, err := r.queries.DeleteResource(
		ctx,
		db.DeleteResourceParams{
			Kind: r.kind,
			ID:   id,
		},
	)
	if err != nil {
		return nil, mapError(err)
	}

	result := resourceFromDB[E](&res)

	return &result, nil
}

func (r *resourceRepository[E]) GetIDsByID(
	ctx context.Context,
	id string,
) ([]string, error) {
	ids, err := r.queries.GetResourceIDsByID(
		ctx,
		db.GetResourceIDsByIDParams{
			Kind: r.kind,
			Prefix: sql.NullString{
				String: id,
				Valid:  true,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *resourceRepository[E]) GetIDsByName(
	ctx context.Context,
	name string,
) ([]string, error) {
	ids, err := r.queries.GetResourceIDsByName(
		ctx,
		db.GetResourceIDsByNameParams{
			Kind: r.kind,
			Prefix: sql.NullString{
				String: name,
				Valid:  true,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *resourceRepository[E]) List(
	ctx context.Context,
	params *repository.ListParams,
) ([]E, error) {
	if params.Limit <= 0 {
		return nil, repository.ErrInvalidPageLimit
	}

	pageLimit := int64(params.Limit)

	var (
		resources []db.Resource
		err       error
	)

	if params.ID == "" {
		resources, err = r.listFirstPage(
			ctx,
			params.OrderBy,
			params.OrderDirection,
			pageLimit,
		)
	} else {
		resources, err = r.listCursorPage(
			ctx,
			params,
			pageLimit,
		)
	}

	if err != nil {
		return nil, err
	}

	return resourcesFromDB[E](resources), nil
}

func (r *resourceRepository[E]) listFirstPage(
	ctx context.Context,
	orderBy repository.OrderBy,
	direction repository.OrderDirection,
	pageLimit int64,
) ([]db.Resource, error) {
	switch {
	case orderBy == repository.OrderByID &&
		direction == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByIDFirstAsc(
			ctx,
			db.ListResourcesByIDFirstAscParams{
				Kind:      r.kind,
				PageLimit: pageLimit,
			},
		)

	case orderBy == repository.OrderByID &&
		direction == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByIDFirstDesc(
			ctx,
			db.ListResourcesByIDFirstDescParams{
				Kind:      r.kind,
				PageLimit: pageLimit,
			},
		)

	case orderBy == repository.OrderByName &&
		direction == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByNameFirstAsc(
			ctx,
			db.ListResourcesByNameFirstAscParams{
				Kind:      r.kind,
				PageLimit: pageLimit,
			},
		)

	case orderBy == repository.OrderByName &&
		direction == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByNameFirstDesc(
			ctx,
			db.ListResourcesByNameFirstDescParams{
				Kind:      r.kind,
				PageLimit: pageLimit,
			},
		)

	case orderBy == repository.OrderByVersion &&
		direction == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByVersionFirstAsc(
			ctx,
			db.ListResourcesByVersionFirstAscParams{
				Kind:      r.kind,
				PageLimit: pageLimit,
			},
		)

	case orderBy == repository.OrderByVersion &&
		direction == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByVersionFirstDesc(
			ctx,
			db.ListResourcesByVersionFirstDescParams{
				Kind:      r.kind,
				PageLimit: pageLimit,
			},
		)

	case orderBy == repository.OrderByCreatedAt &&
		direction == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByCreatedAtFirstAsc(
			ctx,
			db.ListResourcesByCreatedAtFirstAscParams{
				Kind:      r.kind,
				PageLimit: pageLimit,
			},
		)

	case orderBy == repository.OrderByCreatedAt &&
		direction == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByCreatedAtFirstDesc(
			ctx,
			db.ListResourcesByCreatedAtFirstDescParams{
				Kind:      r.kind,
				PageLimit: pageLimit,
			},
		)

	case orderBy == repository.OrderByUpdatedAt &&
		direction == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByUpdatedAtFirstAsc(
			ctx,
			db.ListResourcesByUpdatedAtFirstAscParams{
				Kind:      r.kind,
				PageLimit: pageLimit,
			},
		)

	case orderBy == repository.OrderByUpdatedAt &&
		direction == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByUpdatedAtFirstDesc(
			ctx,
			db.ListResourcesByUpdatedAtFirstDescParams{
				Kind:      r.kind,
				PageLimit: pageLimit,
			},
		)

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

func (r *resourceRepository[E]) listCursorPage(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.Resource, error) {
	switch params.OrderBy {
	case repository.OrderByID:
		return r.listByID(ctx, params, pageLimit)

	case repository.OrderByName:
		return r.listByName(ctx, params, pageLimit)

	case repository.OrderByVersion:
		return r.listByVersion(ctx, params, pageLimit)

	case repository.OrderByCreatedAt:
		return r.listByCreatedAt(ctx, params, pageLimit)

	case repository.OrderByUpdatedAt:
		return r.listByUpdatedAt(ctx, params, pageLimit)

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

//nolint:dupl // Explicit dispatch to distinct sqlc-generated query methods is clearer than abstracting it.
func (r *resourceRepository[E]) listByID(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.Resource, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByIDAfterAsc(
			ctx,
			db.ListResourcesByIDAfterAscParams{
				Kind:      r.kind,
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByIDAfterDesc(
			ctx,
			db.ListResourcesByIDAfterDescParams{
				Kind:      r.kind,
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByIDBeforeAsc(
			ctx,
			db.ListResourcesByIDBeforeAscParams{
				Kind:      r.kind,
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByIDBeforeDesc(
			ctx,
			db.ListResourcesByIDBeforeDescParams{
				Kind:      r.kind,
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

//nolint:dupl // Explicit dispatch to distinct sqlc-generated query methods is clearer than abstracting it.
func (r *resourceRepository[E]) listByName(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.Resource, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByNameAfterAsc(
			ctx,
			db.ListResourcesByNameAfterAscParams{
				Kind:        r.kind,
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByNameAfterDesc(
			ctx,
			db.ListResourcesByNameAfterDescParams{
				Kind:        r.kind,
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByNameBeforeAsc(
			ctx,
			db.ListResourcesByNameBeforeAscParams{
				Kind:        r.kind,
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByNameBeforeDesc(
			ctx,
			db.ListResourcesByNameBeforeDescParams{
				Kind:        r.kind,
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

//nolint:dupl // Explicit dispatch to distinct sqlc-generated query methods is clearer than abstracting it.
func (r *resourceRepository[E]) listByVersion(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.Resource, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByVersionAfterAsc(
			ctx,
			db.ListResourcesByVersionAfterAscParams{
				Kind:        r.kind,
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByVersionAfterDesc(
			ctx,
			db.ListResourcesByVersionAfterDescParams{
				Kind:        r.kind,
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByVersionBeforeAsc(
			ctx,
			db.ListResourcesByVersionBeforeAscParams{
				Kind:        r.kind,
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByVersionBeforeDesc(
			ctx,
			db.ListResourcesByVersionBeforeDescParams{
				Kind:        r.kind,
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

//nolint:dupl // Explicit dispatch to distinct sqlc-generated query methods is clearer than abstracting it.
func (r *resourceRepository[E]) listByCreatedAt(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.Resource, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByCreatedAtAfterAsc(
			ctx,
			db.ListResourcesByCreatedAtAfterAscParams{
				Kind:            r.kind,
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByCreatedAtAfterDesc(
			ctx,
			db.ListResourcesByCreatedAtAfterDescParams{
				Kind:            r.kind,
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByCreatedAtBeforeAsc(
			ctx,
			db.ListResourcesByCreatedAtBeforeAscParams{
				Kind:            r.kind,
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByCreatedAtBeforeDesc(
			ctx,
			db.ListResourcesByCreatedAtBeforeDescParams{
				Kind:            r.kind,
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

//nolint:dupl // Explicit dispatch to distinct sqlc-generated query methods is clearer than abstracting it.
func (r *resourceRepository[E]) listByUpdatedAt(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.Resource, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByUpdatedAtAfterAsc(
			ctx,
			db.ListResourcesByUpdatedAtAfterAscParams{
				Kind:            r.kind,
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByUpdatedAtAfterDesc(
			ctx,
			db.ListResourcesByUpdatedAtAfterDescParams{
				Kind:            r.kind,
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListResourcesByUpdatedAtBeforeAsc(
			ctx,
			db.ListResourcesByUpdatedAtBeforeAscParams{
				Kind:            r.kind,
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListResourcesByUpdatedAtBeforeDesc(
			ctx,
			db.ListResourcesByUpdatedAtBeforeDescParams{
				Kind:            r.kind,
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

func resourceFromDB[E resource](value *db.Resource) E {
	return E(repository.Resource{
		ID:          value.ID,
		Kind:        value.Kind,
		Version:     value.Version,
		Name:        value.Name,
		Description: value.Description,
		Spec:        value.Spec,
		CreatedAt:   value.CreatedAt,
		UpdatedAt:   value.UpdatedAt,
	})
}

func resourcesFromDB[E resource](values []db.Resource) []E {
	result := make([]E, len(values))

	for i := range values {
		result[i] = resourceFromDB[E](&values[i])
	}

	return result
}
