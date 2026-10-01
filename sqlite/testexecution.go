package sqlite

import (
	"context"
	"database/sql"

	"github.com/e2engine/core/repository"

	db "github.com/e2engine/repository/sqlite/internal/sqlc"
)

type testExecutionRepository struct {
	queries db.Querier
}

var _ repository.TestExecutionRepository = (*testExecutionRepository)(nil)

func newTestExecutionRepository(queries db.Querier) *testExecutionRepository {
	return &testExecutionRepository{
		queries: queries,
	}
}

func (r *testExecutionRepository) Create(
	ctx context.Context,
	params *repository.TestExecution,
) (*repository.TestExecution, error) {
	created, err := r.queries.CreateTestExecution(
		ctx,
		db.CreateTestExecutionParams{
			ID: params.ID,
			TestSuiteExecutionID: sql.NullString{
				String: params.TestSuiteExecutionID,
				Valid:  params.TestSuiteExecutionID != "",
			},
			StartedAt:       params.StartedAt,
			FinishedAt:      params.FinishedAt,
			Status:          params.Status,
			EnvironmentID:   params.EnvironmentID,
			EnvironmentName: params.EnvironmentName,
			TestID:          params.TestID,
			TestName:        params.TestName,
			Summary:         params.Summary,
		},
	)
	if err != nil {
		return nil, mapError(err)
	}

	result := testExecutionFromDB(&created)

	return &result, nil
}

func (r *testExecutionRepository) Get(
	ctx context.Context,
	id string,
) (*repository.TestExecution, error) {
	res, err := r.queries.GetTestExecution(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	result := testExecutionFromDB(&res)

	return &result, nil
}

func (r *testExecutionRepository) GetStatus(
	ctx context.Context,
	id string,
) (string, error) {
	res, err := r.queries.GetTestExecutionStatus(ctx, id)
	if err != nil {
		return "", mapError(err)
	}

	return res, nil
}

func (r *testExecutionRepository) Delete(
	ctx context.Context,
	id string,
) (*repository.TestExecution, error) {
	res, err := r.queries.DeleteTestExecution(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	result := testExecutionFromDB(&res)

	return &result, nil
}

func (r *testExecutionRepository) GetIDsByID(
	ctx context.Context,
	id string,
) ([]string, error) {
	ids, err := r.queries.GetTestExecutionIDsByID(
		ctx,
		sql.NullString{
			String: id,
			Valid:  true,
		},
	)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *testExecutionRepository) GetByTestSuiteExecutionID(
	ctx context.Context,
	testSuiteExecutionID string,
) ([]repository.TestExecution, error) {
	values, err := r.queries.GetTestExecutionsByTestSuiteExecutionID(
		ctx,
		testSuiteExecutionID,
	)
	if err != nil {
		return nil, err
	}

	return testExecutionsFromDB(values), nil
}

func (r *testExecutionRepository) List(
	ctx context.Context,
	params *repository.ListParams,
) ([]repository.TestExecution, error) {
	if params.Limit <= 0 {
		return nil, repository.ErrInvalidPageLimit
	}

	pageLimit := int64(params.Limit)

	var (
		executions []db.TestExecution
		err        error
	)

	if params.ID == "" {
		executions, err = r.listTestExecutionsFirstPage(
			ctx,
			params.OrderBy,
			params.OrderDirection,
			pageLimit,
		)
	} else {
		executions, err = r.listTestExecutionsCursorPage(
			ctx,
			params,
			pageLimit,
		)
	}

	if err != nil {
		return nil, err
	}

	return testExecutionsFromDB(executions), nil
}

func (r *testExecutionRepository) listTestExecutionsFirstPage(
	ctx context.Context,
	orderBy repository.OrderBy,
	direction repository.OrderDirection,
	pageLimit int64,
) ([]db.TestExecution, error) {
	switch {
	case orderBy == repository.OrderByID &&
		direction == repository.OrderDirectionAsc:
		return r.queries.ListTestExecutionsByIDFirstAsc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByID &&
		direction == repository.OrderDirectionDesc:
		return r.queries.ListTestExecutionsByIDFirstDesc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByStartedAt &&
		direction == repository.OrderDirectionAsc:
		return r.queries.ListTestExecutionsByStartedAtFirstAsc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByStartedAt &&
		direction == repository.OrderDirectionDesc:
		return r.queries.ListTestExecutionsByStartedAtFirstDesc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByFinishedAt &&
		direction == repository.OrderDirectionAsc:
		return r.queries.ListTestExecutionsByFinishedAtFirstAsc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByFinishedAt &&
		direction == repository.OrderDirectionDesc:
		return r.queries.ListTestExecutionsByFinishedAtFirstDesc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByStatus &&
		direction == repository.OrderDirectionAsc:
		return r.queries.ListTestExecutionsByStatusFirstAsc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByStatus &&
		direction == repository.OrderDirectionDesc:
		return r.queries.ListTestExecutionsByStatusFirstDesc(
			ctx,
			pageLimit,
		)

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

func (r *testExecutionRepository) listTestExecutionsCursorPage(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.TestExecution, error) {
	switch params.OrderBy {
	case repository.OrderByID:
		return r.listTestExecutionsByID(
			ctx,
			params,
			pageLimit,
		)

	case repository.OrderByStartedAt:
		return r.listTestExecutionsByStartedAt(
			ctx,
			params,
			pageLimit,
		)

	case repository.OrderByFinishedAt:
		return r.listTestExecutionsByFinishedAt(
			ctx,
			params,
			pageLimit,
		)

	case repository.OrderByStatus:
		return r.listTestExecutionsByStatus(
			ctx,
			params,
			pageLimit,
		)

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

func (r *testExecutionRepository) listTestExecutionsByID(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.TestExecution, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListTestExecutionsByIDAfterAsc(
			ctx,
			db.ListTestExecutionsByIDAfterAscParams{
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListTestExecutionsByIDAfterDesc(
			ctx,
			db.ListTestExecutionsByIDAfterDescParams{
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListTestExecutionsByIDBeforeAsc(
			ctx,
			db.ListTestExecutionsByIDBeforeAscParams{
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListTestExecutionsByIDBeforeDesc(
			ctx,
			db.ListTestExecutionsByIDBeforeDescParams{
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

//nolint:dupl // Explicit dispatch to distinct sqlc-generated query methods is clearer than abstracting it.
func (r *testExecutionRepository) listTestExecutionsByStartedAt(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.TestExecution, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListTestExecutionsByStartedAtAfterAsc(
			ctx,
			db.ListTestExecutionsByStartedAtAfterAscParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListTestExecutionsByStartedAtAfterDesc(
			ctx,
			db.ListTestExecutionsByStartedAtAfterDescParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListTestExecutionsByStartedAtBeforeAsc(
			ctx,
			db.ListTestExecutionsByStartedAtBeforeAscParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListTestExecutionsByStartedAtBeforeDesc(
			ctx,
			db.ListTestExecutionsByStartedAtBeforeDescParams{
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
func (r *testExecutionRepository) listTestExecutionsByFinishedAt(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.TestExecution, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListTestExecutionsByFinishedAtAfterAsc(
			ctx,
			db.ListTestExecutionsByFinishedAtAfterAscParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListTestExecutionsByFinishedAtAfterDesc(
			ctx,
			db.ListTestExecutionsByFinishedAtAfterDescParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListTestExecutionsByFinishedAtBeforeAsc(
			ctx,
			db.ListTestExecutionsByFinishedAtBeforeAscParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListTestExecutionsByFinishedAtBeforeDesc(
			ctx,
			db.ListTestExecutionsByFinishedAtBeforeDescParams{
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
func (r *testExecutionRepository) listTestExecutionsByStatus(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.TestExecution, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListTestExecutionsByStatusAfterAsc(
			ctx,
			db.ListTestExecutionsByStatusAfterAscParams{
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListTestExecutionsByStatusAfterDesc(
			ctx,
			db.ListTestExecutionsByStatusAfterDescParams{
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.queries.ListTestExecutionsByStatusBeforeAsc(
			ctx,
			db.ListTestExecutionsByStatusBeforeAscParams{
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.queries.ListTestExecutionsByStatusBeforeDesc(
			ctx,
			db.ListTestExecutionsByStatusBeforeDescParams{
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

func (r *testExecutionRepository) SetRunning(
	ctx context.Context,
	params *repository.SetRunningParams,
) error {
	if _, err := r.queries.SetTestExecutionRunning(
		ctx,
		db.SetTestExecutionRunningParams{
			ID:        params.ID,
			StartedAt: params.StartedAt,
		},
	); err != nil {
		return mapError(err)
	}

	return nil
}

func (r *testExecutionRepository) SetCompleted(
	ctx context.Context,
	params *repository.SetCompletedParams,
) error {
	if _, err := r.queries.SetTestExecutionCompleted(
		ctx,
		db.SetTestExecutionCompletedParams{
			ID:         params.ID,
			FinishedAt: params.FinishedAt,
			Status:     params.Status,
			Summary:    params.Summary,
		},
	); err != nil {
		return mapError(err)
	}

	return nil
}

func testExecutionsFromDB(
	values []db.TestExecution,
) []repository.TestExecution {
	result := make([]repository.TestExecution, len(values))

	for i := range values {
		result[i] = testExecutionFromDB(&values[i])
	}

	return result
}

func testExecutionFromDB(
	value *db.TestExecution,
) repository.TestExecution {
	return repository.TestExecution{
		ID:                   value.ID,
		TestSuiteExecutionID: fromNullString(value.TestSuiteExecutionID),
		StartedAt:            value.StartedAt,
		FinishedAt:           value.FinishedAt,
		Status:               value.Status,
		EnvironmentID:        value.EnvironmentID,
		EnvironmentName:      value.EnvironmentName,
		TestID:               value.TestID,
		TestName:             value.TestName,
		Summary:              value.Summary,
	}
}

func fromNullString(value sql.NullString) string {
	if value.Valid {
		return value.String
	}

	return ""
}
