# 🧪 Yinsen Testing Strategy & Harness
*Comprehensive Quality Assurance for Satori Platform*

## 🎯 Testing Philosophy

**Core Principle**: *"Move fast without breaking things"*
- Test the critical path, not every edge case
- Regression protection for core features
- Performance benchmarks as first-class citizens
- Simple, maintainable test suite

## 🛡️ Test Architecture

### Test Pyramid Structure
```
    🔺 E2E Tests (Few, High-Value)
   🔺🔺 Integration Tests (API + Components)
  🔺🔺🔺 Unit Tests (Go Standard Library)
 🔺🔺🔺🔺 Static Analysis (go vet, golint)
```

### Test Categories

#### 1. Unit Tests (`tests/unit/`)
**Framework**: Go's built-in `testing` package + `testify/assert`
**Coverage Target**: Core business logic functions
**Execution**: `go test ./...`

```go
// Example test structure
func TestSenseiAgent_ProcessQuery(t *testing.T) {
    agent := NewSenseiAgent(mockConfig)
    
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"simple query", "hello", "greeting response", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := agent.ProcessQuery(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Contains(t, result, tt.expected)
        })
    }
}
```

#### 2. API Tests (`tests/api/`)
**Framework**: Custom HTTP test harness
**Coverage**: Sensei server endpoints
**Execution**: Dedicated test runner

```go
// HTTP endpoint validation
func TestSenseiAPI_HealthCheck(t *testing.T) {
    server := setupTestServer()
    defer server.Close()
    
    resp, err := http.Get(server.URL + "/health")
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    body, _ := io.ReadAll(resp.Body)
    assert.Contains(t, string(body), "healthy")
}
```

#### 3. Integration Tests (`tests/integration/`)
**Framework**: Docker Compose for dependencies (test only!)
**Coverage**: Zilliz Milvus + n8n interactions
**Execution**: CI/CD pipeline

```bash
# Integration test setup
docker-compose -f docker-compose.test.yml up -d milvus n8n
go test -tags=integration ./tests/integration/...
docker-compose -f docker-compose.test.yml down
```

#### 4. End-to-End Tests (`tests/e2e/`)
**Framework**: Playwright + Go HTTP clients
**Coverage**: Full user workflows
**Execution**: Pre-deployment validation

## 🚀 Performance Testing

### Benchmark Suite
**Location**: `tests/performance/`
**Framework**: Go's `testing.B` benchmarks

```go
func BenchmarkSenseiAgent_ProcessQuery(b *testing.B) {
    agent := NewSenseiAgent(prodConfig)
    query := "complex semantic analysis request"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := agent.ProcessQuery(query)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Performance Targets
- **Query Response**: < 100ms (95th percentile)
- **Memory Usage**: < 512MB steady state
- **Milvus Operations**: < 50ms vector search
- **n8n Workflow Triggers**: < 200ms end-to-end

## 🔄 Test Data Management

### Fixtures & Mocks
**Location**: `tests/fixtures/`
**Strategy**: 
- Minimal, realistic test data
- No sensitive information
- Version controlled
- Easy to regenerate

```go
// Test data factory
func CreateTestPersona() *Persona {
    return &Persona{
        ID:          "test-persona-001",
        Name:        "Test Assistant",
        Personality: "helpful and concise",
        Model:       "claude-3-7-sonnet",
    }
}
```

### Database State
- **Unit Tests**: In-memory mocks
- **Integration Tests**: Isolated test database
- **E2E Tests**: Dedicated test environment

## 🛠️ Test Automation

### CI/CD Pipeline Integration
```yaml
# .github/workflows/test.yml (example)
test:
  runs-on: ubuntu-latest
  steps:
    - name: Run Unit Tests
      run: go test -v ./...
    
    - name: Run Integration Tests
      run: |
        docker-compose -f docker-compose.test.yml up -d
        go test -tags=integration ./tests/integration/...
        docker-compose -f docker-compose.test.yml down
    
    - name: Performance Benchmarks
      run: go test -bench=. ./tests/performance/
```

### Test Runners
- **Local Development**: `make test`
- **CI/CD**: GitHub Actions / GitLab CI
- **Manual**: Individual test commands

## 📊 Test Reporting

### Coverage Metrics
- **Target**: 80% line coverage for core modules
- **Tool**: `go test -coverprofile=coverage.out`
- **Visualization**: HTML reports via `go tool cover`

### Performance Monitoring
- Benchmark trends over time
- Memory leak detection
- Response time degradation alerts

## 🚨 Regression Testing

### Critical Path Protection
1. **User Authentication**: Persona loading + validation
2. **Query Processing**: Input → Sensei → Response
3. **Vector Search**: Milvus semantic operations
4. **Workflow Triggers**: n8n integration points
5. **Error Handling**: Graceful degradation

### Test Execution Schedule
- **Pre-commit**: Unit tests (< 30 seconds)
- **PR/MR**: Unit + Integration tests
- **Pre-deployment**: Full test suite + performance
- **Post-deployment**: Smoke tests + health checks

## 🔧 Test Utilities

### Custom Assertions
```go
// Domain-specific test helpers
func AssertValidPersona(t *testing.T, persona *Persona) {
    assert.NotEmpty(t, persona.ID)
    assert.NotEmpty(t, persona.Name)
    assert.Contains(t, []string{"claude-3-7-sonnet", "claude-3-opus"}, persona.Model)
}

func AssertResponseTime(t *testing.T, duration time.Duration, threshold time.Duration) {
    assert.Less(t, duration, threshold, "Response time exceeded threshold")
}
```

### Test Environment Management
```bash
#!/bin/bash
# scripts/test-env.sh
set -e

echo "Setting up test environment..."
export MILVUS_HOST=localhost:19530
export N8N_HOST=localhost:5678
export SENSEI_ENV=test

# Start dependencies
docker-compose -f docker-compose.test.yml up -d

# Wait for services
./scripts/wait-for-it.sh localhost:19530 -- echo "Milvus ready"
./scripts/wait-for-it.sh localhost:5678 -- echo "n8n ready"

echo "Test environment ready!"
```

## 📝 Test Documentation

### Writing Good Tests
1. **Clear naming**: TestFeature_Scenario_Expectation
2. **AAA Pattern**: Arrange, Act, Assert
3. **Single responsibility**: One concept per test
4. **Readable failures**: Descriptive error messages
5. **Maintainable**: Easy to update with code changes

### Test Case Documentation
- Link to feature specifications
- Document test rationale
- Explain complex setup/teardown
- Note external dependencies

---

## 🎖️ Testing Checklist

### Before Feature Complete
- [ ] Unit tests for new functions
- [ ] Integration tests for external APIs
- [ ] Performance benchmarks (if applicable)
- [ ] Error scenario validation
- [ ] Documentation updates

### Before Release
- [ ] Full test suite passes
- [ ] Performance benchmarks stable
- [ ] No memory leaks detected
- [ ] Critical path regression tests green
- [ ] Test coverage targets met

---

*"Quality is not an act, but a habit. Testing is the foundation of confidence."*

**Author**: Yinsen
**Version**: 1.0
**Last Updated**: ${new Date().toISOString()}