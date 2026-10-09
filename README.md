# E2Engine Repository

Repository implementations for E2Engine.

Part of [E2Engine](https://e2engine.dev), an open-source platform for declarative end-to-end testing.

This module provides persistence for E2Engine resources and execution data.

## Installation

```bash
go get github.com/e2engine/repository
```

## Packages

**repository**

Repository interfaces and common types.

**sqlite**

SQLite implementation of the E2Engine repositories.

## Development

Run tests:

```bash
make test
```

Run the linter:

```bash
make lint
```

Run race detection:

```bash
make test-race
```

## E2Engine

This repository is part of E2Engine.

- [core](https://github.com/e2engine/core) — core domain model, execution logic, and public APIs
- [repository](https://github.com/e2engine/repository) — persistence implementations
- [runner-local](https://github.com/e2engine/runner-local) — local test execution
- [cli](https://github.com/e2engine/cli) — command-line interface
- [tests](https://github.com/e2engine/tests) — end-to-end tests for E2Engine
- [demo](https://github.com/e2engine/demo) — executable demonstration system and E2Engine usage examples
- [instrumentation-go](https://github.com/e2engine/instrumentation-go) — Go instrumentation library for E2Engine

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE).