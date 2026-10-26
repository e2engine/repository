package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"

	"github.com/e2engine/core/pkg/keys"
	"github.com/e2engine/core/repository"
	"github.com/ygrebnov/errorc"

	db "github.com/e2engine/repository/sqlite/internal/sqlc"
)

type testExecutionRepository struct {
	store db.Store
}

var _ repository.TestExecutionRepository = (*testExecutionRepository)(nil)

func newTestExecutionRepository(store db.Store) *testExecutionRepository {
	return &testExecutionRepository{
		store: store,
	}
}

func (r *testExecutionRepository) Create(
	ctx context.Context,
	params *repository.TestExecution,
) (*repository.TestExecution, error) {
	created, err := r.store.CreateTestExecution(
		ctx,
		db.CreateTestExecutionParams{
			ID: params.ID,
			TestSuiteExecutionID: sql.NullString{
				String: params.TestSuiteExecutionID,
				Valid:  params.TestSuiteExecutionID != "",
			},
			StartedAt:       params.StartedAt,
			FinishedAt:      params.FinishedAt,
			Status:          string(params.Status),
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
	res, err := r.store.GetTestExecution(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	result := testExecutionFromDB(&res)

	return &result, nil
}

func (r *testExecutionRepository) GetStatus(
	ctx context.Context,
	id string,
) (repository.ExecutionStatus, error) {
	res, err := r.store.GetTestExecutionStatus(ctx, id)
	if err != nil {
		return "", mapError(err)
	}

	return repository.ExecutionStatus(res), nil
}

func (r *testExecutionRepository) Delete(
	ctx context.Context,
	id string,
) (*repository.TestExecution, error) {
	res, err := r.store.DeleteTestExecution(ctx, id)
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
	ids, err := r.store.GetTestExecutionIDsByID(
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
	values, err := r.store.GetTestExecutionsByTestSuiteExecutionID(
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
		return r.store.ListTestExecutionsByIDFirstAsc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByID &&
		direction == repository.OrderDirectionDesc:
		return r.store.ListTestExecutionsByIDFirstDesc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByStartedAt &&
		direction == repository.OrderDirectionAsc:
		return r.store.ListTestExecutionsByStartedAtFirstAsc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByStartedAt &&
		direction == repository.OrderDirectionDesc:
		return r.store.ListTestExecutionsByStartedAtFirstDesc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByFinishedAt &&
		direction == repository.OrderDirectionAsc:
		return r.store.ListTestExecutionsByFinishedAtFirstAsc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByFinishedAt &&
		direction == repository.OrderDirectionDesc:
		return r.store.ListTestExecutionsByFinishedAtFirstDesc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByStatus &&
		direction == repository.OrderDirectionAsc:
		return r.store.ListTestExecutionsByStatusFirstAsc(
			ctx,
			pageLimit,
		)

	case orderBy == repository.OrderByStatus &&
		direction == repository.OrderDirectionDesc:
		return r.store.ListTestExecutionsByStatusFirstDesc(
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
		return r.store.ListTestExecutionsByIDAfterAsc(
			ctx,
			db.ListTestExecutionsByIDAfterAscParams{
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.store.ListTestExecutionsByIDAfterDesc(
			ctx,
			db.ListTestExecutionsByIDAfterDescParams{
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.store.ListTestExecutionsByIDBeforeAsc(
			ctx,
			db.ListTestExecutionsByIDBeforeAscParams{
				AnchorID:  params.ID,
				PageLimit: pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.store.ListTestExecutionsByIDBeforeDesc(
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

func (r *testExecutionRepository) listTestExecutionsByStartedAt(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.TestExecution, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.store.ListTestExecutionsByStartedAtAfterAsc(
			ctx,
			db.ListTestExecutionsByStartedAtAfterAscParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.store.ListTestExecutionsByStartedAtAfterDesc(
			ctx,
			db.ListTestExecutionsByStartedAtAfterDescParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.store.ListTestExecutionsByStartedAtBeforeAsc(
			ctx,
			db.ListTestExecutionsByStartedAtBeforeAscParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.store.ListTestExecutionsByStartedAtBeforeDesc(
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

func (r *testExecutionRepository) listTestExecutionsByFinishedAt(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.TestExecution, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.store.ListTestExecutionsByFinishedAtAfterAsc(
			ctx,
			db.ListTestExecutionsByFinishedAtAfterAscParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.store.ListTestExecutionsByFinishedAtAfterDesc(
			ctx,
			db.ListTestExecutionsByFinishedAtAfterDescParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.store.ListTestExecutionsByFinishedAtBeforeAsc(
			ctx,
			db.ListTestExecutionsByFinishedAtBeforeAscParams{
				AnchorTimestamp: params.ValueTimestamp,
				AnchorID:        params.ID,
				PageLimit:       pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.store.ListTestExecutionsByFinishedAtBeforeDesc(
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

func (r *testExecutionRepository) listTestExecutionsByStatus(
	ctx context.Context,
	params *repository.ListParams,
	pageLimit int64,
) ([]db.TestExecution, error) {
	switch {
	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.store.ListTestExecutionsByStatusAfterAsc(
			ctx,
			db.ListTestExecutionsByStatusAfterAscParams{
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	case params.Position == repository.PositionAfter &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.store.ListTestExecutionsByStatusAfterDesc(
			ctx,
			db.ListTestExecutionsByStatusAfterDescParams{
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionAsc:
		return r.store.ListTestExecutionsByStatusBeforeAsc(
			ctx,
			db.ListTestExecutionsByStatusBeforeAscParams{
				AnchorValue: params.ValueString,
				AnchorID:    params.ID,
				PageLimit:   pageLimit,
			},
		)

	case params.Position == repository.PositionBefore &&
		params.OrderDirection == repository.OrderDirectionDesc:
		return r.store.ListTestExecutionsByStatusBeforeDesc(
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

func (r *testExecutionRepository) GetTestSuiteExecutionStates(
	ctx context.Context,
	testSuiteExecutionID string,
) ([]repository.TestSuiteExecutionState, bool, error) {
	var (
		states           []repository.TestSuiteExecutionState
		alreadyCompleted bool
	)

	err := r.store.ExecTx(ctx, func(q *db.Queries) error {
		status, err := q.GetTestSuiteExecutionStatus(
			ctx,
			testSuiteExecutionID,
		)
		if err != nil {
			return mapError(err)
		}

		switch repository.ExecutionStatus(status) {
		case repository.ExecutionStatusScheduled,
			repository.ExecutionStatusRunning:
			// Continue with child state aggregation.

		case repository.ExecutionStatusPassed,
			repository.ExecutionStatusFailed,
			repository.ExecutionStatusError:
			alreadyCompleted = true
			return nil

		default:
			return repository.ErrInvalidExecutionStatus
		}

		rows, err := q.GetTestSuiteExecutionStats(
			ctx,
			testSuiteExecutionID,
		)
		if err != nil {
			return mapError(err)
		}

		states = make([]repository.TestSuiteExecutionState, len(rows))

		for i, row := range rows {
			if row.Count <= 0 || row.Count > int64(math.MaxInt) {
				return repository.ErrTooManyResults
			}

			finishedAt, err := timeFromInterface(row.LastTestFinishedAt)
			if err != nil {
				return err
			}

			states[i] = repository.TestSuiteExecutionState{
				Status:             repository.ExecutionStatus(row.Status),
				Count:              int(row.Count),
				LastTestFinishedAt: finishedAt,
			}
		}

		return nil
	})

	if err != nil {
		return nil, false, err
	}

	return states, alreadyCompleted, nil
}

func timeFromInterface(value any) (time.Time, error) {
	switch v := value.(type) {
	case nil:
		return time.Time{}, nil

	case time.Time:
		return v, nil

	case *time.Time:
		if v == nil {
			return time.Time{}, nil
		}

		return *v, nil

	case string:
		return parseSQLiteTime(v)

	default:
		return time.Time{}, errorc.With(
			repository.ErrInvalidType,
			errorc.String(keys.Validation, "expected timestamp"),
		)
	}
}

func parseSQLiteTime(value string) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,
		"2006-01-02 15:04:05.999999999 -0700 MST",
	}

	for _, format := range formats {
		if timestamp, err := time.Parse(format, value); err == nil {
			return timestamp, nil
		}
	}

	return time.Time{}, errorc.With(
		repository.ErrInvalidType,
		errorc.String(keys.Validation, "invalid timestamp"),
	)
}

func (r *testExecutionRepository) SetRunning(
	ctx context.Context,
	params *repository.SetRunningParams,
) error {
	return r.store.ExecTx(ctx, func(q *db.Queries) error {
		_, err := q.SetTestExecutionRunning(
			ctx,
			db.SetTestExecutionRunningParams{
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

		_, err = q.GetTestExecutionStatus(ctx, params.ID)
		if err != nil {
			return mapError(err)
		}

		return repository.ErrInvalidExecutionStateTransition
	})
}

func (r *testExecutionRepository) SetCompleted(
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
		_, err := q.SetTestExecutionCompleted(
			ctx,
			db.SetTestExecutionCompletedParams{
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

		_, err = q.GetTestExecutionStatus(ctx, params.ID)
		if err != nil {
			return mapError(err)
		}

		return repository.ErrInvalidExecutionStateTransition
	})
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
		Status:               repository.ExecutionStatus(value.Status),
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
