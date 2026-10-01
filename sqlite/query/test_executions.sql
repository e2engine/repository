-- ============================================================================
-- CRUD
-- ============================================================================

-- name: CreateTestExecution :one
INSERT INTO test_executions (
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
)
VALUES (
           sqlc.arg(id),
           sqlc.narg(test_suite_execution_id),
           sqlc.arg(started_at),
           sqlc.arg(finished_at),
           sqlc.arg(status),
           sqlc.arg(environment_id),
           sqlc.arg(environment_name),
           sqlc.arg(test_id),
           sqlc.arg(test_name),
           sqlc.arg(summary)
       )
    RETURNING
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary;


-- name: GetTestExecution :one
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE id = sqlc.arg(id);


-- name: DeleteTestExecution :one
DELETE FROM test_executions
WHERE id = sqlc.arg(id)
    RETURNING
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary;


-- ============================================================================
-- Prefix lookup
-- ============================================================================

-- name: GetTestExecutionIDsByID :many
SELECT id
FROM test_executions
WHERE id LIKE sqlc.arg(prefix) || '%'
ORDER BY id ASC;


-- ============================================================================
-- Parent lookup
-- ============================================================================

-- name: GetTestExecutionIDsByTestSuiteExecutionID :many
SELECT id
FROM test_executions
WHERE test_suite_execution_id = sqlc.arg(test_suite_execution_id)
ORDER BY id ASC;


-- name: GetTestExecutionsByTestSuiteExecutionID :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE test_suite_execution_id =
      CAST(sqlc.arg(test_suite_execution_id) AS TEXT)
ORDER BY id ASC;


-- ============================================================================
-- First page: ID
-- ============================================================================

-- name: ListTestExecutionsByIDFirstAsc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
ORDER BY id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestExecutionsByIDFirstDesc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
ORDER BY id DESC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Cursor page: ID
-- ============================================================================

-- PositionAfter + ASC
-- name: ListTestExecutionsByIDAfterAsc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE id > CAST(sqlc.arg(anchor_id) AS TEXT)
ORDER BY id ASC
    LIMIT sqlc.arg(page_limit);


-- PositionAfter + DESC
-- name: ListTestExecutionsByIDAfterDesc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE id < CAST(sqlc.arg(anchor_id) AS TEXT)
ORDER BY id DESC
    LIMIT sqlc.arg(page_limit);


-- PositionBefore + ASC
-- Results are queried in reverse order.
-- The repository must restore ASC order.
-- name: ListTestExecutionsByIDBeforeAsc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE id < CAST(sqlc.arg(anchor_id) AS TEXT)
ORDER BY id DESC
    LIMIT sqlc.arg(page_limit);


-- PositionBefore + DESC
-- Results are queried in reverse order.
-- The repository must restore DESC order.
-- name: ListTestExecutionsByIDBeforeDesc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE id > CAST(sqlc.arg(anchor_id) AS TEXT)
ORDER BY id ASC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- First page: StartedAt
-- ============================================================================

-- name: ListTestExecutionsByStartedAtFirstAsc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
ORDER BY started_at ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestExecutionsByStartedAtFirstDesc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
ORDER BY started_at DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Cursor page: StartedAt
-- ============================================================================

-- PositionAfter + ASC
-- name: ListTestExecutionsByStartedAtAfterAsc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE (started_at, id) >
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY started_at ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- PositionAfter + DESC
-- name: ListTestExecutionsByStartedAtAfterDesc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE (started_at, id) <
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY started_at DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- PositionBefore + ASC
-- Results are queried in reverse order.
-- The repository must restore ASC order.
-- name: ListTestExecutionsByStartedAtBeforeAsc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE (started_at, id) <
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY started_at DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- PositionBefore + DESC
-- Results are queried in reverse order.
-- The repository must restore DESC order.
-- name: ListTestExecutionsByStartedAtBeforeDesc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE (started_at, id) >
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY started_at ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- First page: FinishedAt
-- ============================================================================

-- name: ListTestExecutionsByFinishedAtFirstAsc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
ORDER BY finished_at ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestExecutionsByFinishedAtFirstDesc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
ORDER BY finished_at DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Cursor page: FinishedAt
-- ============================================================================

-- PositionAfter + ASC
-- name: ListTestExecutionsByFinishedAtAfterAsc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE (finished_at, id) >
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY finished_at ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- PositionAfter + DESC
-- name: ListTestExecutionsByFinishedAtAfterDesc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE (finished_at, id) <
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY finished_at DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- PositionBefore + ASC
-- Results are queried in reverse order.
-- The repository must restore ASC order.
-- name: ListTestExecutionsByFinishedAtBeforeAsc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE (finished_at, id) <
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY finished_at DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- PositionBefore + DESC
-- Results are queried in reverse order.
-- The repository must restore DESC order.
-- name: ListTestExecutionsByFinishedAtBeforeDesc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE (finished_at, id) >
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY finished_at ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- First page: Status
-- ============================================================================

-- name: ListTestExecutionsByStatusFirstAsc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
ORDER BY status ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestExecutionsByStatusFirstDesc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
ORDER BY status DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Cursor page: Status
-- ============================================================================

-- PositionAfter + ASC
-- name: ListTestExecutionsByStatusAfterAsc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE (status, id) >
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY status ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- PositionAfter + DESC
-- name: ListTestExecutionsByStatusAfterDesc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE (status, id) <
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY status DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- PositionBefore + ASC
-- Results are queried in reverse order.
-- The repository must restore ASC order.
-- name: ListTestExecutionsByStatusBeforeAsc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE (status, id) <
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY status DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- PositionBefore + DESC
-- Results are queried in reverse order.
-- The repository must restore DESC order.
-- name: ListTestExecutionsByStatusBeforeDesc :many
SELECT
    id,
    test_suite_execution_id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_id,
    test_name,
    summary
FROM test_executions
WHERE (status, id) >
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY status ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Status
-- ============================================================================

-- name: GetTestExecutionStatus :one
SELECT status
FROM test_executions
WHERE id = sqlc.arg(id);


-- ============================================================================
-- State transitions
-- ============================================================================

-- name: SetTestExecutionRunning :one
UPDATE test_executions
SET
    started_at = sqlc.arg(started_at),
    status = 'running'
WHERE id = sqlc.arg(id)
RETURNING id;


-- name: SetTestExecutionCompleted :one
UPDATE test_executions
SET
    finished_at = sqlc.arg(finished_at),
    status = sqlc.arg(status),
    summary = sqlc.arg(summary)
WHERE id = sqlc.arg(id)
RETURNING id;