PRAGMA foreign_keys = ON;

CREATE TABLE resources (
                           id          TEXT NOT NULL PRIMARY KEY,
                           kind        TEXT NOT NULL,
                           version     TEXT NOT NULL,
                           name        TEXT NOT NULL,
                           description TEXT NOT NULL DEFAULT '',
                           spec        BLOB NOT NULL,
                           created_at  DATETIME NOT NULL,
                           updated_at  DATETIME NOT NULL,

                           CONSTRAINT resources_kind_check
                               CHECK (
                                   kind IN (
                                            'environment',
                                            'test',
                                            'test_suite'
                                       )
                                   ),

                           CONSTRAINT resources_id_not_empty_check
                               CHECK (length(id) > 0),

                           CONSTRAINT resources_version_not_empty_check
                               CHECK (length(version) > 0),

                           CONSTRAINT resources_name_not_empty_check
                               CHECK (length(name) > 0),

                           CONSTRAINT resources_spec_not_empty_check
                               CHECK (length(spec) > 0),

                           CONSTRAINT resources_created_at_not_zero_check
                               CHECK (substr(CAST(created_at AS TEXT), 1, 10) <> '0001-01-01'),

                           CONSTRAINT resources_updated_at_not_zero_check
                               CHECK (substr(CAST(updated_at AS TEXT), 1, 10) <> '0001-01-01'),

                           CONSTRAINT resources_kind_name_unique
                               UNIQUE (kind, name)
);

-- Supports resource lookup and cursor ordering by ID within a kind.
CREATE INDEX resources_kind_id_idx
    ON resources (
                  kind,
                  id
        );

-- Supports prefix lookup and deterministic cursor ordering by name.
CREATE INDEX resources_kind_name_id_idx
    ON resources (
                  kind,
                  name,
                  id
        );

-- Supports cursor ordering by version.
CREATE INDEX resources_kind_version_id_idx
    ON resources (
                  kind,
                  version,
                  id
        );

-- Supports cursor ordering by creation time.
CREATE INDEX resources_kind_created_at_id_idx
    ON resources (
                  kind,
                  created_at,
                  id
        );

-- Supports cursor ordering by update time.
CREATE INDEX resources_kind_updated_at_id_idx
    ON resources (
                  kind,
                  updated_at,
                  id
        );

CREATE TABLE test_suite_executions (
                                       id          TEXT NOT NULL PRIMARY KEY,
                                       started_at  DATETIME NOT NULL,
                                       finished_at DATETIME NOT NULL,
                                       status      TEXT NOT NULL,

                                       environment_id   TEXT NOT NULL,
                                       environment_name TEXT NOT NULL,

                                       test_suite_id   TEXT NOT NULL,
                                       test_suite_name TEXT NOT NULL,

                                       summary BLOB,

                                       CONSTRAINT test_suite_executions_id_not_empty_check
                                           CHECK (length(id) > 0),

                                       CONSTRAINT test_suite_executions_status_check
                                           CHECK (
                                               status IN (
                                                          'scheduled',
                                                          'running',
                                                          'passed',
                                                          'failed',
                                                          'error'
                                                   )
                                               ),

                                       CONSTRAINT environment_id_not_empty_check
                                           CHECK (length(environment_id) > 0),

                                       CONSTRAINT environment_name_not_empty_check
                                           CHECK (length(environment_name) > 0),

                                       CONSTRAINT test_suite_id_not_empty_check
                                           CHECK (length(test_suite_id) > 0),

                                       CONSTRAINT test_suite_name_not_empty_check
                                           CHECK (length(test_suite_name) > 0)
);

-- Supports cursor ordering by start time.
CREATE INDEX test_suite_executions_started_at_id_idx
    ON test_suite_executions (
                              started_at,
                              id
        );

-- Supports cursor ordering by finish time.
CREATE INDEX test_suite_executions_finished_at_id_idx
    ON test_suite_executions (
                              finished_at,
                              id
        );

-- Supports cursor ordering by status.
CREATE INDEX test_suite_executions_status_id_idx
    ON test_suite_executions (
                              status,
                              id
        );

-- Supports filtering executions by environment.
CREATE INDEX test_suite_executions_environment_id_idx
    ON test_suite_executions (
                              environment_id
        );

-- Supports prefix lookup by test-suite ID with deterministic ID ordering.
CREATE INDEX test_suite_executions_test_suite_id_id_idx
    ON test_suite_executions (
                              test_suite_id,
                              id
        );

CREATE TABLE test_tags (
                           test_id TEXT NOT NULL,
                           tag     TEXT NOT NULL,

                           CONSTRAINT test_tags_pk
                               PRIMARY KEY (test_id, tag),

                           CONSTRAINT test_tags_test_id_not_empty_check
                               CHECK (length(test_id) > 0),

                           CONSTRAINT test_tags_tag_not_empty_check
                               CHECK (length(tag) > 0),

                           CONSTRAINT test_tags_test_fk
                               FOREIGN KEY (test_id)
                                   REFERENCES resources (id)
                                   ON UPDATE CASCADE
                                   ON DELETE CASCADE
);

-- Supports resolving tests by tag with deterministic test ID ordering.
CREATE INDEX test_tags_tag_test_id_idx
    ON test_tags (
                  tag,
                  test_id
        );

CREATE TABLE test_executions (
                                 id                      TEXT NOT NULL PRIMARY KEY,
                                 test_suite_execution_id TEXT,

                                 started_at  DATETIME NOT NULL,
                                 finished_at DATETIME NOT NULL,
                                 status      TEXT NOT NULL,

                                 environment_id   TEXT NOT NULL,
                                 environment_name TEXT NOT NULL,

                                 test_id   TEXT NOT NULL,
                                 test_name TEXT NOT NULL,

                                 summary BLOB,

                                 CONSTRAINT test_executions_id_not_empty_check
                                     CHECK (length(id) > 0),

                                 CONSTRAINT test_executions_status_check
                                     CHECK (
                                         status IN (
                                                  'scheduled',
                                                  'running',
                                                  'passed',
                                                  'failed',
                                                  'error'
                                             )
                                         ),

                                 CONSTRAINT test_executions_environment_id_not_empty_check
                                     CHECK (length(environment_id) > 0),

                                 CONSTRAINT test_executions_environment_name_not_empty_check
                                     CHECK (length(environment_name) > 0),

                                 CONSTRAINT test_executions_test_id_not_empty_check
                                     CHECK (length(test_id) > 0),

                                 CONSTRAINT test_executions_test_name_not_empty_check
                                     CHECK (length(test_name) > 0),

                                 CONSTRAINT test_executions_test_suite_execution_fk
                                     FOREIGN KEY (test_suite_execution_id)
                                         REFERENCES test_suite_executions (id)
                                         ON UPDATE CASCADE
                                         ON DELETE CASCADE
);

-- Supports prefix lookup by parent test-suite execution ID with deterministic
-- child execution ID ordering.
CREATE INDEX test_executions_test_suite_execution_id_id_idx
    ON test_executions (
                        test_suite_execution_id,
                        id
        );

-- Supports cursor ordering by start time.
CREATE INDEX test_executions_started_at_id_idx
    ON test_executions (
                        started_at,
                        id
        );

-- Supports cursor ordering by finish time.
CREATE INDEX test_executions_finished_at_id_idx
    ON test_executions (
                        finished_at,
                        id
        );

-- Supports cursor ordering by status.
CREATE INDEX test_executions_status_id_idx
    ON test_executions (
                        status,
                        id
        );

-- Supports filtering executions by environment.
CREATE INDEX test_executions_environment_id_idx
    ON test_executions (
                        environment_id
        );

-- Supports filtering executions by test.
CREATE INDEX test_executions_test_id_idx
    ON test_executions (
                        test_id
        );