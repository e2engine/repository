-- ============================================================================
-- CRUD
-- ============================================================================

-- name: CreateTestSuiteExecution :one
INSERT INTO test_suite_executions (
    id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_suite_id,
    test_suite_name,
    summary
)
VALUES (
           sqlc.arg(id),
           sqlc.arg(started_at),
           sqlc.arg(finished_at),
           sqlc.arg(status),
           sqlc.arg(environment_id),
           sqlc.arg(environment_name),
           sqlc.arg(test_suite_id),
           sqlc.arg(test_suite_name),
           sqlc.narg(summary)
       )
    RETURNING
    id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_suite_id,
    test_suite_name,
    summary;

-- name: GetTestSuiteExecution :one
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE tse.id = sqlc.arg(id);

-- name: DeleteTestSuiteExecution :one
DELETE FROM test_suite_executions
WHERE id = sqlc.arg(id)
    RETURNING
    id,
    started_at,
    finished_at,
    status,
    environment_id,
    environment_name,
    test_suite_id,
    test_suite_name,
    summary;


-- ============================================================================
-- Prefix lookup
-- ============================================================================

-- name: GetTestSuiteExecutionIDsByID :many
SELECT id
FROM test_suite_executions
WHERE id LIKE sqlc.arg(prefix) || '%'
ORDER BY id ASC;


-- name: GetTestSuiteExecutionIDsByTestSuiteID :many
SELECT id
FROM test_suite_executions
WHERE test_suite_id LIKE sqlc.arg(prefix) || '%'
ORDER BY test_suite_id ASC, id ASC;


-- ============================================================================
-- First page: ID
-- ============================================================================

-- name: ListTestSuiteExecutionsByIDFirstAsc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
ORDER BY tse.id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByIDFirstDesc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
ORDER BY tse.id DESC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Cursor page: ID
-- ============================================================================

-- name: ListTestSuiteExecutionsByIDAfterAsc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE tse.id > CAST(sqlc.arg(anchor_id) AS TEXT)
ORDER BY tse.id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByIDAfterDesc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE tse.id < CAST(sqlc.arg(anchor_id) AS TEXT)
ORDER BY tse.id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByIDBeforeAsc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE tse.id < CAST(sqlc.arg(anchor_id) AS TEXT)
ORDER BY tse.id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByIDBeforeDesc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE tse.id > CAST(sqlc.arg(anchor_id) AS TEXT)
ORDER BY tse.id ASC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- First page: StartedAt
-- ============================================================================

-- name: ListTestSuiteExecutionsByStartedAtFirstAsc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
ORDER BY tse.started_at ASC, tse.id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByStartedAtFirstDesc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
ORDER BY tse.started_at DESC, tse.id DESC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Cursor page: StartedAt
-- ============================================================================

-- name: ListTestSuiteExecutionsByStartedAtAfterAsc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE (tse.started_at, tse.id) >
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY tse.started_at ASC, tse.id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByStartedAtAfterDesc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE (tse.started_at, tse.id) <
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY tse.started_at DESC, tse.id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByStartedAtBeforeAsc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE (tse.started_at, tse.id) <
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY tse.started_at DESC, tse.id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByStartedAtBeforeDesc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE (tse.started_at, tse.id) >
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY tse.started_at ASC, tse.id ASC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- First page: FinishedAt
-- ============================================================================

-- name: ListTestSuiteExecutionsByFinishedAtFirstAsc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
ORDER BY tse.finished_at ASC, tse.id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByFinishedAtFirstDesc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
ORDER BY tse.finished_at DESC, tse.id DESC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Cursor page: FinishedAt
-- ============================================================================

-- name: ListTestSuiteExecutionsByFinishedAtAfterAsc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE (tse.finished_at, tse.id) >
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY tse.finished_at ASC, tse.id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByFinishedAtAfterDesc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE (tse.finished_at, tse.id) <
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY tse.finished_at DESC, tse.id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByFinishedAtBeforeAsc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE (tse.finished_at, tse.id) <
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY tse.finished_at DESC, tse.id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByFinishedAtBeforeDesc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE (tse.finished_at, tse.id) >
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY tse.finished_at ASC, tse.id ASC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- First page: Status
-- ============================================================================

-- name: ListTestSuiteExecutionsByStatusFirstAsc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
ORDER BY tse.status ASC, tse.id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByStatusFirstDesc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
ORDER BY tse.status DESC, tse.id DESC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Cursor page: Status
-- ============================================================================

-- name: ListTestSuiteExecutionsByStatusAfterAsc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE (tse.status, tse.id) >
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY tse.status ASC, tse.id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByStatusAfterDesc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE (tse.status, tse.id) <
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY tse.status DESC, tse.id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByStatusBeforeAsc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE (tse.status, tse.id) <
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY tse.status DESC, tse.id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListTestSuiteExecutionsByStatusBeforeDesc :many
SELECT
    tse.id,
    tse.started_at,
    tse.finished_at,
    tse.status,
    tse.environment_id,
    tse.environment_name,
    tse.test_suite_id,
    tse.test_suite_name,
    (
        SELECT COUNT(*)
        FROM test_executions te
        WHERE te.test_suite_execution_id = tse.id
    ) AS tests_count,
    tse.summary
FROM test_suite_executions tse
WHERE (tse.status, tse.id) >
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY tse.status ASC, tse.id ASC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Status
-- ============================================================================

-- name: GetTestSuiteExecutionStatus :one
SELECT status
FROM test_suite_executions
WHERE id = sqlc.arg(id);


-- ============================================================================
-- State transitions
-- ============================================================================

-- name: SetTestSuiteExecutionRunning :one
UPDATE test_suite_executions
SET
    started_at = sqlc.arg(started_at),
    status = 'running'
WHERE id = sqlc.arg(id)
RETURNING id;


-- name: SetTestSuiteExecutionCompleted :one
UPDATE test_suite_executions
SET
    finished_at = sqlc.arg(finished_at),
    status = sqlc.arg(status),
    summary = sqlc.narg(summary)
WHERE id = sqlc.arg(id)
RETURNING id;