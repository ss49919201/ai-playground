# DI Container Design

## Overview

This DI (Dependency Injection) container provides a lightweight solution for managing dependencies and avoiding tight coupling between components.

## Architecture

### Core Components

- **DIContainer Interface**: Defines the contract for dependency registration and resolution
- **Container Implementation**: Uses a Map-based registry to store service factories
- **Service Registration**: Dependencies are registered with string keys and factory functions
- **Service Resolution**: Services are resolved by key and instantiated on-demand

### Key Features

- **Type Safety**: Generic methods ensure type safety at compile time
- **Factory Pattern**: Services are created using factory functions for flexibility
- **Error Handling**: Throws descriptive errors for unregistered services
- **Lazy Instantiation**: Services are created only when requested

## Usage

### Basic Registration and Resolution

```typescript
const container = createContainer();

// Register services
container.register('db', () => newDB());
container.register('searchUserUsecase', () => 
  newSearchUserUsecase(container.get<DB>('db'))
);

// Resolve services
const usecase = container.get<ReturnType<typeof newSearchUserUsecase>>('searchUserUsecase');
```

### Handler Implementation

Before DI Container (tight coupling):
```typescript
export const SearchUserHandler = (input: { ids: number[]; limit: number }) => {
  const db = newDB();                    // Direct instantiation
  const usecase = newSearchUserUsecase(db); // Manual dependency injection
  return usecase.exec(input);
};
```

After DI Container (loose coupling):
```typescript
export const SearchUserHandler = (input: { ids: number[]; limit: number }) => {
  const usecase = container.get<ReturnType<typeof newSearchUserUsecase>>('searchUserUsecase');
  return usecase.exec(input);
};
```

## Benefits

1. **Decoupling**: Handlers no longer need to know how to construct their dependencies
2. **Testability**: Easy to mock dependencies by registering test implementations
3. **Flexibility**: Can swap implementations without changing consumer code
4. **Centralized Configuration**: All dependency wiring is defined in one place

## Design Decisions

- **String-based Keys**: Simple and flexible, though could be replaced with symbols for better type safety
- **Factory Functions**: Allow for complex initialization logic and ensure fresh instances
- **Singleton Pattern**: Container instance is shared globally for consistency
- **No Lifecycle Management**: Services are created on-demand without automatic disposal

## Future Enhancements

- Add lifecycle management (singleton, transient scopes)
- Implement symbol-based keys for better type safety
- Add circular dependency detection
- Support for dependency decoration/interception