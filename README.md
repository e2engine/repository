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

## Links

- [E2Engine website](https://e2engine.dev)
- [E2Engine documentation](https://e2engine.dev/docs)
- [E2Engine on GitHub](https://github.com/e2engine)

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE).