package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/e2engine/core/repository"

	db "github.com/e2engine/repository/sqlite/internal/sqlc"
)

//go:generate go run ./internal/generate/testsuiteexecution

type testSuiteExecutionRepository struct {
	store db.Store
}

var _ repository.TestSuiteExecutionRepository = (*testSuiteExecutionRepository)(nil)

func newTestSuiteExecutionRepository(
	store db.Store,
) *testSuiteExecutionRepository {
	return &testSuiteExecutionRepository{
		store: store,
	}
}

func (r *testSuiteExecutionRepository) Create(
	ctx context.Context,
	params *repository.TestSuiteExecution,
) (*repository.TestSuiteExecution, error) {
	raw, err := r.store.CreateTestSuiteExecution(
		ctx,
		db.CreateTestSuiteExecutionParams{
			ID:              params.ID,
			StartedAt:       params.StartedAt,
			FinishedAt:      params.FinishedAt,
			Status:          string(params.Status),
			EnvironmentID:   params.EnvironmentID,
			EnvironmentName: params.EnvironmentName,
			TestSuiteID:     params.TestSuiteID,
			TestSuiteName:   params.TestSuiteName,
		},
	)
	if err != nil {
		return nil, mapError(err)
	}

	return &repository.TestSuiteExecution{
		ID:              raw.ID,
		StartedAt:       raw.StartedAt,
		FinishedAt:      raw.FinishedAt,
		Status:          repository.ExecutionStatus(raw.Status),
		EnvironmentID:   raw.EnvironmentID,
		EnvironmentName: raw.EnvironmentName,
		TestSuiteID:     raw.TestSuiteID,
		TestSuiteName:   raw.TestSuiteName,
		Summary:         raw.Summary,
	}, nil
}

func (r *testSuiteExecutionRepository) Get(
	ctx context.Context,
	id string,
) (*repository.TestSuiteExecution, error) {
	raw, err := r.store.GetTestSuiteExecution(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	return &repository.TestSuiteExecution{
		ID:              raw.ID,
		StartedAt:       raw.StartedAt,
		FinishedAt:      raw.FinishedAt,
		Status:          repository.ExecutionStatus(raw.Status),
		EnvironmentID:   raw.EnvironmentID,
		EnvironmentName: raw.EnvironmentName,
		TestSuiteID:     raw.TestSuiteID,
		TestSuiteName:   raw.TestSuiteName,
		TestsCount:      int(raw.TestsCount),
		Summary:         raw.Summary,
	}, nil
}

func (r *testSuiteExecutionRepository) GetStatus(
	ctx context.Context,
	id string,
) (repository.ExecutionStatus, error) {
	res, err := r.store.GetTestSuiteExecutionStatus(ctx, id)
	if err != nil {
		return "", mapError(err)
	}

	return repository.ExecutionStatus(res), nil
}

func (r *testSuiteExecutionRepository) Delete(
	ctx context.Context,
	id string,
) (*repository.TestSuiteExecution, error) {
	raw, err := r.store.DeleteTestSuiteExecution(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	return &repository.TestSuiteExecution{
		ID:              raw.ID,
		StartedAt:       raw.StartedAt,
		FinishedAt:      raw.FinishedAt,
		Status:          repository.ExecutionStatus(raw.Status),
		EnvironmentID:   raw.EnvironmentID,
		EnvironmentName: raw.EnvironmentName,
		TestSuiteID:     raw.TestSuiteID,
		TestSuiteName:   raw.TestSuiteName,
		Summary:         raw.Summary,
	}, nil
}

func (r *testSuiteExecutionRepository) GetIDsByID(
	ctx context.Context,
	id string,
) ([]string, error) {
	ids, err := r.store.GetTestSuiteExecutionIDsByID(
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

func (r *testSuiteExecutionRepository) List(
	ctx context.Context,
	params *repository.ListParams,
) ([]repository.TestSuiteExecution, error) {
	if params.Limit <= 0 {
		return nil, repository.ErrInvalidPageLimit
	}

	pageLimit := int64(params.Limit)

	var (
		executions []testSuiteExecutionRow
		err        error
	)

	if params.ID == "" {
		executions, err = r.listFirstPage(
			ctx,
			params.OrderBy,
			params.OrderDirection,
			pageLimit,
		)
	} else {
		executions, err = r.listCursorPage(
			ctx,
			params,
			pageLimit,
		)
	}

	if err != nil {
		return nil, err
	}

	return testSuiteExecutionsFromDB(executions), nil
}

func (r *testSuiteExecutionRepository) listFirstPage(
	ctx context.Context,
	orderBy repository.OrderBy,
	direction repository.OrderDirection,
	pageLimit int64,
) ([]testSuiteExecutionRow, error) {
	switch {
	case orderBy == repository.OrderByID &&
		direction == repository.OrderDirectionAsc:
		list, err := r.store.ListTestSuiteExecutionsByIDFirstAsc(
			ctx,
			pageLimit,
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByIDFirstAsc(list), nil

	case orderBy == repository.OrderByID &&
		direction == repository.OrderDirectionDesc:
		list, err := r.store.ListTestSuiteExecutionsByIDFirstDesc(
			ctx,
			pageLimit,
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByIDFirstDesc(list), nil

	case orderBy == repository.OrderByStartedAt &&
		direction == repository.OrderDirectionAsc:
		list, err := r.store.ListTestSuiteExecutionsByStartedAtFirstAsc(
			ctx,
			pageLimit,
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByStartedAtFirstAsc(list), nil

	case orderBy == repository.OrderByStartedAt &&
		direction == repository.OrderDirectionDesc:
		list, err := r.store.ListTestSuiteExecutionsByStartedAtFirstDesc(
			ctx,
			pageLimit,
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByStartedAtFirstDesc(list), nil

	case orderBy == repository.OrderByFinishedAt &&
		direction == repository.OrderDirectionAsc:
		list, err := r.store.ListTestSuiteExecutionsByFinishedAtFirstAsc(
			ctx,
			pageLimit,
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByFinishedAtFirstAsc(list), nil

	case orderBy == repository.OrderByFinishedAt &&
		direction == repository.OrderDirectionDesc:
		list, err := r.store.ListTestSuiteExecutionsByFinishedAtFirstDesc(
			ctx,
			pageLimit,
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByFinishedAtFirstDesc(list), nil

	case orderBy == repository.OrderByStatus &&
		direction == repository.OrderDirectionAsc:
		list, err := r.store.ListTestSuiteExecutionsByStatusFirstAsc(
			ctx,
			pageLimit,
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByStatusFirstAsc(list), nil

	case orderBy == repository.OrderByStatus &&
		direction == repository.OrderDirectionDesc:
		list, err := r.store.ListTestSuiteExecutionsByStatusFirstDesc(
			ctx,
			pageLimit,
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByStatusFirstDesc(list), nil

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

func (r *testSuiteExecutionRepository) listCursorPage(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]testSuiteExecutionRow, error) {
	switch params.OrderBy {
	case repository.OrderByID:
		return r.listByID(ctx, params, pageLimit)

	case repository.OrderByStartedAt:
		return r.listByStartedAt(ctx, params, pageLimit)

	case repository.OrderByFinishedAt:
		return r.listByFinishedAt(ctx, params, pageLimit)

	case repository.OrderByStatus:
		return r.listByStatus(ctx, params, pageLimit)

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

func (r *testSuiteExecutionRepository) listByID(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]testSuiteExecutionRow, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		list, err := r.store.ListTestSuiteExecutionsByIDAfterAsc(
			ctx,
			db.ListTestSuiteExecutionsByIDAfterAscParams{
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByIDAfterAsc(list), nil

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		list, err := r.store.ListTestSuiteExecutionsByIDAfterDesc(
			ctx,
			db.ListTestSuiteExecutionsByIDAfterDescParams{
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByIDAfterDesc(list), nil

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		list, err := r.store.ListTestSuiteExecutionsByIDBeforeAsc(
			ctx,
			db.ListTestSuiteExecutionsByIDBeforeAscParams{
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByIDBeforeAsc(list), nil

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		list, err := r.store.ListTestSuiteExecutionsByIDBeforeDesc(
			ctx,
			db.ListTestSuiteExecutionsByIDBeforeDescParams{
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByIDBeforeDesc(list), nil

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

func (r *testSuiteExecutionRepository) listByStartedAt(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]testSuiteExecutionRow, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		list, err := r.store.ListTestSuiteExecutionsByStartedAtAfterAsc(
			ctx,
			db.ListTestSuiteExecutionsByStartedAtAfterAscParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByStartedAtAfterAsc(list), nil

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		list, err := r.store.ListTestSuiteExecutionsByStartedAtAfterDesc(
			ctx,
			db.ListTestSuiteExecutionsByStartedAtAfterDescParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByStartedAtAfterDesc(list), nil

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		list, err := r.store.ListTestSuiteExecutionsByStartedAtBeforeAsc(
			ctx,
			db.ListTestSuiteExecutionsByStartedAtBeforeAscParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByStartedAtBeforeAsc(list), nil

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		list, err := r.store.ListTestSuiteExecutionsByStartedAtBeforeDesc(
			ctx,
			db.ListTestSuiteExecutionsByStartedAtBeforeDescParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByStartedAtBeforeDesc(list), nil

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

func (r *testSuiteExecutionRepository) listByFinishedAt(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]testSuiteExecutionRow, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		list, err := r.store.ListTestSuiteExecutionsByFinishedAtAfterAsc(
			ctx,
			db.ListTestSuiteExecutionsByFinishedAtAfterAscParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByFinishedAtAfterAsc(list), nil

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		list, err := r.store.ListTestSuiteExecutionsByFinishedAtAfterDesc(
			ctx,
			db.ListTestSuiteExecutionsByFinishedAtAfterDescParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByFinishedAtAfterDesc(list), nil

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		list, err := r.store.ListTestSuiteExecutionsByFinishedAtBeforeAsc(
			ctx,
			db.ListTestSuiteExecutionsByFinishedAtBeforeAscParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByFinishedAtBeforeAsc(list), nil

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		list, err := r.store.ListTestSuiteExecutionsByFinishedAtBeforeDesc(
			ctx,
			db.ListTestSuiteExecutionsByFinishedAtBeforeDescParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByFinishedAtBeforeDesc(list), nil

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

func (r *testSuiteExecutionRepository) listByStatus(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]testSuiteExecutionRow, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		list, err := r.store.ListTestSuiteExecutionsByStatusAfterAsc(
			ctx,
			db.ListTestSuiteExecutionsByStatusAfterAscParams{
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByStatusAfterAsc(list), nil

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		list, err := r.store.ListTestSuiteExecutionsByStatusAfterDesc(
			ctx,
			db.ListTestSuiteExecutionsByStatusAfterDescParams{
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByStatusAfterDesc(list), nil

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		list, err := r.store.ListTestSuiteExecutionsByStatusBeforeAsc(
			ctx,
			db.ListTestSuiteExecutionsByStatusBeforeAscParams{
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByStatusBeforeAsc(list), nil

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		list, err := r.store.ListTestSuiteExecutionsByStatusBeforeDesc(
			ctx,
			db.ListTestSuiteExecutionsByStatusBeforeDescParams{
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)
		if err != nil {
			return nil, err
		}
		return normalizeListTestSuiteExecutionsByStatusBeforeDesc(list), nil

	default:
		return nil, repository.ErrInvalidPageTraversal
	}
}

func (r *testSuiteExecutionRepository) SetRunning(
	ctx context.Context,
	params *repository.SetRunningParams,
) error {
	return r.store.ExecTx(ctx, func(q *db.Queries) error {
		_, err := q.SetTestSuiteExecutionRunning(
			ctx,
			db.SetTestSuiteExecutionRunningParams{
				ID:        params.ID,
				StartedAt: params.StartedAt,
			},
		)
		if err == nil {
			return nil
		}

		if !errors.Is(err, sql.ErrNoRows) {
			return mapError(err)
		}

		_, err = q.GetTestSuiteExecutionStatus(ctx, params.ID)
		if err != nil {
			return mapError(err)
		}

		return repository.ErrInvalidExecutionStateTransition
	})
}

func (r *testSuiteExecutionRepository) SetCompleted(
	ctx context.Context,
	params *repository.SetCompletedParams,
) error {
	switch params.Status {
	case repository.ExecutionStatusPassed,
		repository.ExecutionStatusFailed,
		repository.ExecutionStatusError:
		// Valid terminal status.
	default:
		return repository.ErrInvalidExecutionStatus
	}

	return r.store.ExecTx(ctx, func(q *db.Queries) error {
		_, err := q.SetTestSuiteExecutionCompleted(
			ctx,
			db.SetTestSuiteExecutionCompletedParams{
				ID:         params.ID,
				FinishedAt: params.FinishedAt,
				Status:     string(params.Status),
				Summary:    params.Summary,
			},
		)
		if err == nil {
			return nil
		}

		if !errors.Is(err, sql.ErrNoRows) {
			return mapError(err)
		}

		_, err = q.GetTestSuiteExecutionStatus(ctx, params.ID)
		if err != nil {
			return mapError(err)
		}

		return repository.ErrInvalidExecutionStateTransition
	})
}

type testSuiteExecutionRow struct {
	ID              string    `json:"id"`
	StartedAt       time.Time `json:"started_at"`
	FinishedAt      time.Time `json:"finished_at"`
	Status          string    `json:"status"`
	EnvironmentID   string    `json:"environment_id"`
	EnvironmentName string    `json:"environment_name"`
	TestSuiteID     string    `json:"test_suite_id"`
	TestSuiteName   string    `json:"test_suite_name"`
	TestsCount      int64     `json:"tests_count"`
	Summary         []byte    `json:"summary"`
}

func testSuiteExecutionsFromDB(
	values []testSuiteExecutionRow,
) []repository.TestSuiteExecution {
	result := make([]repository.TestSuiteExecution, len(values))

	for i := range values {
		result[i] = testSuiteExecutionFromDB(&values[i])
	}

	return result
}

func testSuiteExecutionFromDB(
	value *testSuiteExecutionRow,
) repository.TestSuiteExecution {
	return repository.TestSuiteExecution{
		ID:              value.ID,
		StartedAt:       value.StartedAt,
		FinishedAt:      value.FinishedAt,
		Status:          repository.ExecutionStatus(value.Status),
		EnvironmentID:   value.EnvironmentID,
		EnvironmentName: value.EnvironmentName,
		TestSuiteID:     value.TestSuiteID,
		TestSuiteName:   value.TestSuiteName,
		TestsCount:      int(value.TestsCount),
		Summary:         value.Summary,
	}
}
