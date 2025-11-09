



### Run Unit Tests Only

```bash
make test
```

This runs tests with the `-short` flag, skipping integration tests.

### Run Integration Tests Only

```bash
make test-integration
```

This automatically starts the test database and runs all integration tests.

### Run All Tests

```bash
make test-all
```

Runs both unit and integration tests.

### Run Tests with Coverage

```bash
make test-coverage
```

Generates a coverage report in `coverage.html`.

### Stop Test Database

```bash
make test-db-stop
```

## Test Database Configuration

Integration tests connect to the same PostgreSQL database configured in your Docker Compose setup:

- **Host**: localhost
- **Port**: 8010 (configured in `dev-env/local.env`)
- **Database**: postgres
- **Schema**: Auto-initialized from `dev-env/sql/create_tables.sql`

### Database Isolation

Each integration test:
1. Cleans the database before running (via `setupTest()`)
2. Truncates all tables in the correct order
3. Can optionally clean up after completion

This ensures tests are isolated and don't interfere with each other.

## Using Test Builders

Test builders provide a fluent interface for creating test data:

### Task Builder

```go
task := builders.NewTaskBuilder().
    WithName("My Task").
    WithDescription("Task description").
    HighPriority().
    Important().
    DueTomorrow().
    MustBuild()
```

### User Builder

```go
user := builders.NewUserBuilder().
    WithUserName("john_doe").
    WithEmail("john@example.com").
    WithFullName("John", "Doe").
    Build()
```

### UserGroup Builder

```go
group := builders.NewUserGroupBuilder().
    WithName("Development Team").
    WithDescription("Our dev team").
    TeamType().
    MustBuild()
```

### TaskContainer Builder

```go
container := builders.NewTaskContainerBuilder().
    WithName("TODO List").
    WithUsergroupId(groupId).
    WithActivityLevel(5).
    Build()
```

## Writing Integration Tests

### Basic Structure

```go
func TestMyFeature(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    cleanup := setupTest(t)
    defer cleanup()

    // Test code here
}
```

### Example: Testing Task Creation

```go
func TestCreateTask(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    cleanup := setupTest(t)
    defer cleanup()

    // Setup repositories
    taskRepo := repository.NewTaskRepository(testDB)

    // Create test data using builders
    task := builders.NewTaskBuilder().
        WithName("Test Task").
        MustBuild()

    // Execute test
    createdTask, err := taskRepo.CreateTask(containerId, *task)

    // Assertions
    require.NoError(t, err)
    assert.Equal(t, "Test Task", createdTask.TaskName)
}
```

## Test Categories

### User Group Flow Tests

**File**: `integration/usergroup_flow_test.go`

Tests user group management including:
- Creating groups with admin users
- Adding/removing members
- Changing user roles
- Deleting groups

### Task Management Flow Tests

**File**: `integration/task_management_flow_test.go`

Tests task operations including:
- Creating tasks in containers
- Updating task details
- Toggling completion and importance
- Deleting tasks
- Querying tasks by container/group

### Repository Complex Queries Tests

**File**: `integration/repository_complex_queries_test.go`

Tests complex database queries involving:
- Multi-table joins (users with roles)
- Group-level task queries
- Container-level filtering
- Important task filtering
- User lookup methods

## CI/CD Integration

The GitHub Actions workflow automatically runs tests. Add this to your workflow:

```yaml
- name: Run Integration Tests
  run: |
    make test-db-start
    make test-integration
```

## Best Practices

1. **Always use `setupTest()`** - Ensures clean database state
2. **Use test builders** - More readable and maintainable test data
3. **Test one thing per test** - Keep tests focused and clear
4. **Use descriptive test names** - Follow `TestFeature_Scenario_ExpectedResult` pattern
5. **Clean up resources** - Use `defer cleanup()` pattern
6. **Check for `-short` flag** - Skip integration tests during quick test runs

## Troubleshooting

### Tests fail with connection errors

Ensure Docker Compose database is running:
```bash
make start
# or
make test-db-start
```

### Database has stale data

Integration tests should clean the database automatically. If issues persist:
```bash
make down
make start
```

### Tests are slow

Integration tests are slower than unit tests. Use:
```bash
make test  # Unit tests only
```

For development, run specific test files:
```bash
go test -v ./tests/integration/usergroup_flow_test.go
```

## Future Improvements

Potential enhancements for the test suite:

1. **Transaction-based isolation** - Use database transactions instead of truncation
2. **Test fixtures** - Pre-seeded data for common scenarios
3. **API client helper** - End-to-end HTTP request testing
4. **Performance tests** - Load testing for critical paths
5. **Parallel test execution** - Run tests concurrently with isolated databases
