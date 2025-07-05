# Claude Code Rules for Clojure CLI Project

## Code Organization Rules

### Single File Architecture
- All Clojure code should be consolidated into a single file (`src/main.clj`)
- Avoid multiple namespace dependencies unless absolutely necessary
- Keep the main entry point (`-main`) in the same file as business logic

### Namespace Declaration
- Use minimal namespace declarations: `(ns main)`
- Only include `:gen-class` if creating JAR files for Java interop
- For direct execution with `clojure -M -m main`, `:gen-class` is optional
- Include test dependencies: `(:require [clojure.test :refer :all])`

## Testing Standards

### Test-Driven Development (TDD)
- Follow Kent Beck's TDD cycle: Red → Green → Refactor
- Always write failing tests first
- Implement minimal code to pass tests
- Refactor while keeping tests green

### Test Organization
- Include tests in the same file as implementation
- Use `deftest` with descriptive names ending in `-test`
- Group related assertions with `testing` blocks
- Test edge cases and typical use cases

### Test Execution
- Run tests with: `clojure -M -e "(require 'main) (clojure.test/run-tests 'main)"`
- Ensure all tests pass before considering implementation complete

## Code Quality Guidelines

### Function Implementation
- Write pure functions when possible
- Use descriptive function names
- Keep functions focused on single responsibility
- Use `cond` for multiple conditional branches

### Code Style
- Use idiomatic Clojure constructs (`zero?` instead of `(= 0 ...)`)
- Extract complex expressions to `let` bindings for clarity
- Use appropriate sequence functions (`filter`, `map`, `reduce`)
- Prefer built-in functions over manual implementations

## Development Workflow

### Implementation Process
1. Write failing test (Red phase)
2. Implement minimal code to pass (Green phase)
3. Refactor for clarity and performance
4. Verify all tests still pass

### Code Documentation
- Functions should be self-documenting through clear names
- Avoid unnecessary comments unless complex algorithms require explanation
- Focus on readable, idiomatic Clojure code

## Example Patterns

### Prime Number Function
```clojure
(defn prime? [n]
  (cond
    (< n 2) false
    (= n 2) true
    (even? n) false
    :else (let [sqrt-n (int (Math/sqrt n))]
            (not-any? #(zero? (mod n %)) (range 3 (inc sqrt-n) 2)))))
```

### FizzBuzz Pattern
```clojure
(defn fizzbuzz [n]
  (cond
    (= (mod n 15) 0) "FizzBuzz"
    (= (mod n 3) 0) "Fizz"
    (= (mod n 5) 0) "Buzz"
    :else (str n)))
```

### Test Structure
```clojure
(deftest function-name-test
  (testing "Description of what is being tested"
    (is (= expected-value (function-call args)))
    (is (= another-expected (function-call other-args)))))
```