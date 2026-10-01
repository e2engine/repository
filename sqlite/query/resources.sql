-- resources.sql

-- ============================================================================
-- CRUD
-- ============================================================================

-- name: CreateResource :one
INSERT INTO resources (
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
)
VALUES (
           sqlc.arg(id),
           sqlc.arg(kind),
           sqlc.arg(version),
           sqlc.arg(name),
           sqlc.arg(description),
           sqlc.arg(spec),
           sqlc.arg(created_at),
           sqlc.arg(updated_at)
       )
    RETURNING
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at;


-- name: GetResource :one
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND id = sqlc.arg(id);


-- name: DeleteResource :one
DELETE FROM resources
WHERE kind = sqlc.arg(kind)
  AND id = sqlc.arg(id)
    RETURNING
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at;


-- ============================================================================
-- Prefix lookup
-- ============================================================================

-- name: GetResourceIDsByID :many
SELECT id
FROM resources
WHERE kind = sqlc.arg(kind)
  AND id LIKE sqlc.arg(prefix) || '%'
ORDER BY id ASC;


-- name: GetResourceIDsByName :many
SELECT id
FROM resources
WHERE kind = sqlc.arg(kind)
  AND name LIKE sqlc.arg(prefix) || '%'
ORDER BY name ASC, id ASC;


-- ============================================================================
-- First page: ID
-- ============================================================================

-- name: ListResourcesByIDFirstAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
ORDER BY id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByIDFirstDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
ORDER BY id DESC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Cursor page: ID
-- ============================================================================

-- name: ListResourcesByIDAfterAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND id > CAST(sqlc.arg(anchor_id) AS TEXT)
ORDER BY id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByIDAfterDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND id < CAST(sqlc.arg(anchor_id) AS TEXT)
ORDER BY id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByIDBeforeAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND id < CAST(sqlc.arg(anchor_id) AS TEXT)
ORDER BY id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByIDBeforeDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND id > CAST(sqlc.arg(anchor_id) AS TEXT)
ORDER BY id ASC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- First page: Name
-- ============================================================================

-- name: ListResourcesByNameFirstAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
ORDER BY name ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByNameFirstDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
ORDER BY name DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Cursor page: Name
-- ============================================================================

-- name: ListResourcesByNameAfterAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (name, id) >
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY name ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByNameAfterDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (name, id) <
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY name DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByNameBeforeAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (name, id) <
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY name DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByNameBeforeDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (name, id) >
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY name ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- First page: Version
-- ============================================================================

-- name: ListResourcesByVersionFirstAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
ORDER BY version ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByVersionFirstDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
ORDER BY version DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Cursor page: Version
-- ============================================================================

-- name: ListResourcesByVersionAfterAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (version, id) >
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY version ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByVersionAfterDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (version, id) <
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY version DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByVersionBeforeAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (version, id) <
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY version DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByVersionBeforeDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (version, id) >
      (sqlc.arg(anchor_value), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY version ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- First page: CreatedAt
-- ============================================================================

-- name: ListResourcesByCreatedAtFirstAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
ORDER BY created_at ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByCreatedAtFirstDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
ORDER BY created_at DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Cursor page: CreatedAt
-- ============================================================================

-- name: ListResourcesByCreatedAtAfterAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (created_at, id) >
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY created_at ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByCreatedAtAfterDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (created_at, id) <
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY created_at DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByCreatedAtBeforeAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (created_at, id) <
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY created_at DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByCreatedAtBeforeDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (created_at, id) >
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY created_at ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- First page: UpdatedAt
-- ============================================================================

-- name: ListResourcesByUpdatedAtFirstAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
ORDER BY updated_at ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByUpdatedAtFirstDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
ORDER BY updated_at DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- ============================================================================
-- Cursor page: UpdatedAt
-- ============================================================================

-- name: ListResourcesByUpdatedAtAfterAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (updated_at, id) >
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY updated_at ASC, id ASC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByUpdatedAtAfterDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (updated_at, id) <
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY updated_at DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByUpdatedAtBeforeAsc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (updated_at, id) <
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY updated_at DESC, id DESC
    LIMIT sqlc.arg(page_limit);


-- name: ListResourcesByUpdatedAtBeforeDesc :many
SELECT
    id,
    kind,
    version,
    name,
    description,
    spec,
    created_at,
    updated_at
FROM resources
WHERE kind = sqlc.arg(kind)
  AND (updated_at, id) >
      (sqlc.arg(anchor_timestamp), CAST(sqlc.arg(anchor_id) AS TEXT))
ORDER BY updated_at ASC, id ASC
    LIMIT sqlc.arg(page_limit);