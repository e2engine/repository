package tests

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
	"time"

	coreapi "github.com/e2engine/core/api"
	corerepository "github.com/e2engine/core/repository"
	_ "modernc.org/sqlite"

	"github.com/e2engine/repository/sqlite"
	"github.com/e2engine/repository/sqlite/config"
)

var defaultConfig = &config.Config{
	Name:         "e2engine.sqlite",
	Driver:       "sqlite",
	MaxOpenConns: 1,
	MaxIdleConns: 1,
}

func getRepositories(t *testing.T) coreapi.Repositories {
	t.Helper()

	dbPath := sqlite.GetPath(t.TempDir(), defaultConfig)
	dbConn, err := sqlite.Open(t.Context(), dbPath, defaultConfig)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	return sqlite.NewRepositories(dbConn)
}

const (
	idPrefixEnvironment        = "111"
	idPrefixTest               = "222"
	idPrefixTestSuite          = "333"
	idPrefixTestExecution      = "444"
	idPrefixTestSuiteExecution = "555"
	idBody                     = "8ab887149ec5aa583c564a659a0b0a933ef27bacb675b081a2a752d8a389"
)

func getResource(t *testing.T, kind string, i int) corerepository.Resource {
	t.Helper()

	return corerepository.Resource{
		ID:          idPrefixEnvironment + idBody + strconv.Itoa(i),
		Kind:        kind,
		Version:     strconv.Itoa(i) + ".0.0",
		Name:        "Resource " + kind + " " + strconv.Itoa(i) + " Name",
		Description: "This is the resource" + strconv.Itoa(i) + " for testing.",
		Spec:        json.RawMessage(`{"key": "Resource` + strconv.Itoa(i) + `"}`),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func getTestExecution(
	t *testing.T,
	i int,
	testSuiteExecutionID string,
	startedAt, finishedAt time.Time,
	status string,
) corerepository.TestExecution {
	t.Helper()

	return corerepository.TestExecution{
		ID:                   idPrefixTestExecution + idBody + strconv.Itoa(i),
		TestSuiteExecutionID: testSuiteExecutionID,
		StartedAt:            startedAt,
		FinishedAt:           finishedAt,
		Status:               status,
		EnvironmentID:        idPrefixEnvironment + idBody + strconv.Itoa(i),
		EnvironmentName:      "Environment " + strconv.Itoa(i) + " Name",
		TestID:               idPrefixTest + idBody + strconv.Itoa(i),
		TestName:             "Test " + strconv.Itoa(i) + " Name",
		Summary:              json.RawMessage(`{"key": "TestExecution` + strconv.Itoa(i) + `"}`),
	}
}

func getDefaultTestExecution(t *testing.T, i int) corerepository.TestExecution {
	t.Helper()

	testSuiteExecutionID := idPrefixTestSuiteExecution + idBody + "1"
	startedAt := time.Now().Add(-time.Duration(i) * time.Minute)
	finishedAt := startedAt.Add(time.Duration(i+1) * time.Minute)
	status := "passed"

	return getTestExecution(t, i, testSuiteExecutionID, startedAt, finishedAt, status)
}

func getTestSuiteExecution(
	t *testing.T,
	i int,
	startedAt, finishedAt time.Time,
	status string,
	summary json.RawMessage,
) corerepository.TestSuiteExecution {
	t.Helper()

	return corerepository.TestSuiteExecution{
		ID:              idPrefixTestSuiteExecution + idBody + strconv.Itoa(i),
		StartedAt:       startedAt,
		FinishedAt:      finishedAt,
		Status:          status,
		EnvironmentID:   idPrefixEnvironment + idBody + strconv.Itoa(i),
		EnvironmentName: "Environment " + strconv.Itoa(i) + " Name",
		TestSuiteID:     idPrefixTestSuite + idBody + strconv.Itoa(i),
		TestSuiteName:   "TestSuite " + strconv.Itoa(i) + " Name",
		Summary:         summary,
	}
}

func getDefaultTestSuiteExecution(t *testing.T, i int) corerepository.TestSuiteExecution {
	t.Helper()

	startedAt := time.Now().Add(-time.Duration(i) * time.Minute)
	finishedAt := startedAt.Add(time.Duration(i+1) * time.Minute)
	status := "passed"
	// summary := json.RawMessage(`{"key": "TestSuiteExecution` + strconv.Itoa(i) + `"}`)

	return getTestSuiteExecution(t, i, startedAt, finishedAt, status, nil)
}

func getResourcesForList(t *testing.T, kind string) []corerepository.Resource {
	t.Helper()

	names := []string{
		"Name 03",
		"Name 01",
		"Name 05",
		"Name 02",
		"Name 04",
	}
	versions := []string{
		"Version 02",
		"Version 01",
		"Version 03",
		"Version 01",
		"Version 02",
	}
	createdAtOffsets := []time.Duration{
		2 * time.Second,
		3 * time.Second,
		time.Second,
		2 * time.Second,
		3 * time.Second,
	}
	updatedAtOffsets := []time.Duration{
		4 * time.Second,
		time.Second,
		3 * time.Second,
		time.Second,
		2 * time.Second,
	}
	timestamp := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	resources := make([]corerepository.Resource, len(names))

	for i := range resources {
		resource := getResource(t, kind, i+1)
		resource.Name = names[i]
		resource.Version = versions[i]
		resource.CreatedAt = timestamp.Add(createdAtOffsets[i])
		resource.UpdatedAt = timestamp.Add(updatedAtOffsets[i])
		resources[i] = resource
	}

	return resources
}

func getTestExecutionsForList(t *testing.T, tseID string) []corerepository.TestExecution {
	t.Helper()

	startedAtOffsets := []time.Duration{
		2 * time.Second,
		3 * time.Second,
		time.Second,
		2 * time.Second,
		3 * time.Second,
	}
	finishedAtOffsets := []time.Duration{
		4 * time.Second,
		time.Second,
		3 * time.Second,
		time.Second,
		2 * time.Second,
	}
	timestamp := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	tes := make([]corerepository.TestExecution, len(startedAtOffsets))

	for i := range tes {
		te := getDefaultTestExecution(t, i+1)
		te.TestSuiteExecutionID = tseID
		te.StartedAt = timestamp.Add(startedAtOffsets[i])
		te.FinishedAt = timestamp.Add(finishedAtOffsets[i])
		tes[i] = te
	}

	return tes
}

func createEnvironmentResources(
	t *testing.T,
	repository corerepository.EnvironmentRepository,
	resources []corerepository.Resource,
) {
	t.Helper()

	for i := range resources {
		environment := corerepository.Environment(resources[i])
		created, err := repository.Create(t.Context(), &environment)
		if err != nil {
			t.Fatalf("Unexpected environment%d create error: %v", i+1, err)
		}
		if created == nil {
			t.Fatalf("Expected created%d environment to be non-nil", i+1)
		}
		assertEqualResources(t, resources[i], corerepository.Resource(*created))
	}
}

func createTestResources(
	t *testing.T,
	repository corerepository.TestRepository,
	resources []corerepository.Resource,
) {
	t.Helper()

	for i := range resources {
		test := corerepository.Test(resources[i])
		created, err := repository.Create(t.Context(), &test)
		if err != nil {
			t.Fatalf("Unexpected test%d create error: %v", i+1, err)
		}
		if created == nil {
			t.Fatalf("Expected created%d test to be non-nil", i+1)
		}
		assertEqualResources(t, resources[i], corerepository.Resource(*created))
	}
}

func createTestSuiteResources(
	t *testing.T,
	repository corerepository.TestSuiteRepository,
	resources []corerepository.Resource,
) {
	t.Helper()

	for i := range resources {
		testSuite := corerepository.TestSuite(resources[i])
		created, err := repository.Create(t.Context(), &testSuite)
		if err != nil {
			t.Fatalf("Unexpected test suite%d create error: %v", i+1, err)
		}
		if created == nil {
			t.Fatalf("Expected created%d test suite to be non-nil", i+1)
		}
		assertEqualResources(t, resources[i], corerepository.Resource(*created))
	}
}

func createTestExecutions(
	t *testing.T,
	repository corerepository.TestExecutionRepository,
	tes []corerepository.TestExecution,
) {
	t.Helper()

	for i := range tes {
		te := tes[i]
		created, err := repository.Create(t.Context(), &te)
		if err != nil {
			t.Fatalf("Unexpected test execution%d create error: %v", i+1, err)
		}
		if created == nil {
			t.Fatalf("Expected created%d test execution to be non-nil", i+1)
		}
		assertEqualTestExecutions(t, te, *created)
	}
}

func environmentsAt(
	resources []corerepository.Resource,
	indexes ...int,
) []corerepository.Environment {
	environments := make([]corerepository.Environment, len(indexes))

	for i, index := range indexes {
		environments[i] = corerepository.Environment(resources[index])
	}

	return environments
}

func makeSequentialIndexes(count int) []int {
	indexes := make([]int, count)
	for i := range indexes {
		indexes[i] = i
	}
	return indexes
}

func assertEqualEnvironments(
	t *testing.T,
	expected, actual []corerepository.Environment,
) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Fatalf("Expected %d environments, got %d", len(expected), len(actual))
	}

	for i := range expected {
		assertEqualResources(
			t,
			corerepository.Resource(expected[i]),
			corerepository.Resource(actual[i]),
		)
	}
}

func assertEqualResources(t *testing.T, expected, actual corerepository.Resource) {
	t.Helper()
	if expected.ID != actual.ID {
		t.Errorf("Expected resource ID %v, got %v", expected.ID, actual.ID)
	}
	if expected.Kind != actual.Kind {
		t.Errorf("Expected resource Kind %v, got %v", expected.Kind, actual.Kind)
	}
	if expected.Version != actual.Version {
		t.Errorf("Expected resource Version %v, got %v", expected.Version, actual.Version)
	}
	if expected.Name != actual.Name {
		t.Errorf("Expected resource Name %v, got %v", expected.Name, actual.Name)
	}
	if expected.Description != actual.Description {
		t.Errorf("Expected resource Description %v, got %v", expected.Description, actual.Description)
	}
	if !reflect.DeepEqual(expected.Spec, actual.Spec) {
		t.Errorf("Expected resource Spec %v, got %v", expected.Spec, actual.Spec)
	}
	if !expected.CreatedAt.Equal(actual.CreatedAt) {
		t.Errorf("Expected resource CreatedAt %v, got %v", expected.CreatedAt, actual.CreatedAt)
	}
	if !expected.UpdatedAt.Equal(actual.UpdatedAt) {
		t.Errorf("Expected resource UpdatedAt %v, got %v", expected.UpdatedAt, actual.UpdatedAt)
	}
}

func assertEqualTestExecutions(t *testing.T, expected, actual corerepository.TestExecution) {
	t.Helper()
	if expected.ID != actual.ID {
		t.Errorf("Expected testExecution ID %v, got %v", expected.ID, actual.ID)
	}
	if expected.TestSuiteExecutionID != actual.TestSuiteExecutionID {
		t.Errorf("Expected testExecution TestSuiteExecutionID %v, got %v", expected.TestSuiteExecutionID, actual.TestSuiteExecutionID)
	}
	if !expected.StartedAt.Equal(actual.StartedAt) {
		t.Errorf("Expected testExecution StartedAt %v, got %v", expected.StartedAt, actual.StartedAt)
	}
	if !expected.FinishedAt.Equal(actual.FinishedAt) {
		t.Errorf("Expected testExecution FinishedAt %v, got %v", expected.FinishedAt, actual.FinishedAt)
	}
	if expected.Status != actual.Status {
		t.Errorf("Expected testExecution Status %v, got %v", expected.Status, actual.Status)
	}
	if expected.EnvironmentID != actual.EnvironmentID {
		t.Errorf("Expected testExecution EnvironmentID %v, got %v", expected.EnvironmentID, actual.EnvironmentID)
	}
	if expected.EnvironmentName != actual.EnvironmentName {
		t.Errorf("Expected testExecution EnvironmentName %v, got %v", expected.EnvironmentName, actual.EnvironmentName)
	}
	if expected.TestID != actual.TestID {
		t.Errorf("Expected testExecution TestID %v, got %v", expected.TestID, actual.TestID)
	}
	if expected.TestName != actual.TestName {
		t.Errorf("Expected testExecution TestName %v, got %v", expected.TestName, actual.TestName)
	}
	if !reflect.DeepEqual(expected.Summary, actual.Summary) {
		t.Errorf("Expected testExecution Summary %v, got %v", expected.Summary, actual.Summary)
	}
}

func assertEqualTestSuiteExecutions(t *testing.T, expected, actual corerepository.TestSuiteExecution) {
	t.Helper()
	if expected.ID != actual.ID {
		t.Errorf("Expected testSuiteExecution ID %v, got %v", expected.ID, actual.ID)
	}
	if !expected.StartedAt.Equal(actual.StartedAt) {
		t.Errorf("Expected testSuiteExecution StartedAt %v, got %v", expected.StartedAt, actual.StartedAt)
	}
	if !expected.FinishedAt.Equal(actual.FinishedAt) {
		t.Errorf("Expected testSuiteExecution FinishedAt %v, got %v", expected.FinishedAt, actual.FinishedAt)
	}
	if expected.Status != actual.Status {
		t.Errorf("Expected testSuiteExecution Status %v, got %v", expected.Status, actual.Status)
	}
	if expected.EnvironmentID != actual.EnvironmentID {
		t.Errorf("Expected testSuiteExecution EnvironmentID %v, got %v", expected.EnvironmentID, actual.EnvironmentID)
	}
	if expected.EnvironmentName != actual.EnvironmentName {
		t.Errorf("Expected testSuiteExecution EnvironmentName %v, got %v", expected.EnvironmentName, actual.EnvironmentName)
	}
	if expected.TestSuiteID != actual.TestSuiteID {
		t.Errorf("Expected testSuiteExecution TestSuiteID %v, got %v", expected.TestSuiteID, actual.TestSuiteID)
	}
	if expected.TestSuiteName != actual.TestSuiteName {
		t.Errorf("Expected testSuiteExecution TestSuiteName %v, got %v", expected.TestSuiteName, actual.TestSuiteName)
	}
	if expected.TestsCount != actual.TestsCount {
		t.Errorf("Expected testSuiteExecution TestsCount %v, got %v", expected.TestsCount, actual.TestsCount)
	}
	if !reflect.DeepEqual(expected.Summary, actual.Summary) {
		t.Errorf("Expected testSuiteExecution Summary %v, got %v", expected.Summary, actual.Summary)
	}
}
