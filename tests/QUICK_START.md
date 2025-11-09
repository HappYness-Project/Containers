
### 2. Run Integration Tests

```bash
make test-integration
```

This will:
- Start the test database if not running
- Run all integration tests
- Show detailed test results

### 3. View Results

You should see output like:

```
Running integration tests...
=== RUN   TestUserGroupFlow_CreateGroupWithAdmin
--- PASS: TestUserGroupFlow_CreateGroupWithAdmin (0.05s)
=== RUN   TestTaskFlow_CreateTaskInContainer
--- PASS: TestTaskFlow_CreateTaskInContainer (0.04s)
...
PASS
ok      github.com/happYness-Project/taskManagementGolang/tests/integration    2.431s
```

## Quick Commands

| Command | Description |
|---------|-------------|
| `make test` | Run unit tests only (fast) |
| `make test-integration` | Run integration tests only |
| `make test-all` | Run all tests |
| `make test-coverage` | Generate coverage report |
| `make test-db-start` | Start database only |
| `make test-db-stop` | Stop database only |

## Example: Running a Single Test

```bash
go test -v ./tests/integration/ -run TestUserGroupFlow_CreateGroupWithAdmin
```

## Example: Using Test Builders in Your Code

```go
// Create a user
user := builders.NewUserBuilder().
    WithUserName("alice").
    WithEmail("alice@example.com").
    Build()

// Create a task
task := builders.NewTaskBuilder().
    WithName("Write documentation").
    HighPriority().
    Important().
    DueTomorrow().
    MustBuild()

// Create a group
group := builders.NewUserGroupBuilder().
    WithName("Engineering").
    TeamType().
    MustBuild()
```

## What Was Created

### Test Infrastructure
- ✅ Test setup with automatic database cleanup
- ✅ Test data builders for all domain entities
- ✅ Database helper functions

### Integration Tests (30+ test cases)
- ✅ User group management flow (7 tests)
- ✅ Task management flow (8 tests)
- ✅ Complex repository queries (8+ tests)

### Makefile Targets
- ✅ `test-integration` - Run integration tests
- ✅ `test-all` - Run all tests
- ✅ `test-coverage` - Generate coverage reports
- ✅ `test-db-start/stop` - Control test database

### Documentation
- ✅ Comprehensive README
- ✅ This quick start guide

## Test Coverage

The integration tests cover:

**User & Group Management:**
- Creating groups with admin users
- Adding/removing members
- Changing user roles
- User lookup by various fields
- Complex join queries with roles

**Task Management:**
- Creating tasks in containers
- Updating task details
- Toggling completion status
- Toggling importance
- Deleting tasks
- Querying tasks by container/group

**Repository Layer:**
- Multi-table joins
- Foreign key relationships
- Complex filtering (important tasks only)
- Group-level aggregations

## Next Steps

1. **Run the tests** to ensure everything works
2. **Add more tests** as you develop new features
3. **Check coverage** with `make test-coverage`
4. **Integrate with CI/CD** in your GitHub Actions workflow

## Need Help?

- See [tests/README.md](./README.md) for detailed documentation
- Check existing tests for examples
- Use test builders to create clean test data
