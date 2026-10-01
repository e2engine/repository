-- name: CreateTestTag :exec
INSERT INTO test_tags (
    test_id,
    tag
)
VALUES (
           sqlc.arg(test_id),
           sqlc.arg(tag)
       );

-- name: DeleteTestTagsByTestID :exec
DELETE FROM test_tags
WHERE test_id = sqlc.arg(test_id);

-- name: GetTestsByTag :many
SELECT r.*
FROM resources r
         JOIN test_tags tt
              ON tt.test_id = r.id
WHERE tt.tag = sqlc.arg(tag)
  AND r.kind = 'test'
ORDER BY r.id ASC;