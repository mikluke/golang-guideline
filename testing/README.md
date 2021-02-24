# Packaging

## Index

- [Why do we use tests](#why-do-we-use-tests?)
- [Packaging](#packaging)
- [Mocking](#mocking)

## Topics

### Why do we use tests?

- To avoid manual testing which is very slow and inefficient
- To avoid regression bugs
- To prove interface' design usability
- To speed up development process

#### Which process helps to have code covered?

We suggest following [test-driven development process](https://en.wikipedia.org/wiki/Test-driven_development)  
That process obliges you to code test first and after that code implementation. This way helps you to tune your interface to be more usable.

#### How much code should be covered?

##### Maximal coverage
- Critical code
- Business logic
- Validation logic

##### Medium coverage
- Object allocations
- Small utility functions
- Logging
 
##### Low coverage
- "One-line" functions

##### Don't cover
- `main` function


#### What should I do if my code is already done?

Definitely, you should cover your code with tests. If this looks like impossible, revise your code' design. **Good designed code is always testable.** 

#### What can I do if my manager says that writing tests is a very complicated business?

You have to agree. Code covering by tests is the complicated business, but the manual testing is **more** complicated. After that you have to show reason why we should cover our code by tests. In the end, if it won't persuade your manager, just leave this place and get another job with smart management personnel. 

### Packaging

#### _test.go

Put your tests near to code files
```
serivice
| - service.go
| - service_test.go
```

If you need to use assets within tests put it into a `test` subdirectory
```
image
| - test
    | - image.png
| - reader.go
| - reader_test.go
```

#### Mocked implementation

There are two acceptable ways to store mocks:

##### If you get interfaces close to using place?

Store mocked implementation in {interface_name}_mock_test.go
```
service
| - service.go # contains the Storage interface
| - service_test.go # your tests
| - storage_mock_test.go # mock implementation of the Storage interface
```

Benefits:
- You will be able only to use it in `_test.go` go files 
- A mock implementation won't be imported into the `vendor` directory by `go mod`

##### If you get interfaces close to implementation?

Store mocked implementation in a detached directory with `mock` name
```
storage
| - storage.go # interface
| - mock
    | - storage.go
| - mongo
    | - mongo.go
    | - mongo_test.go
```

Benefits:
- Your code's clients can use mocked implementation for tests directly
- Will be imported into the `vendor` directory

### Mocking

If your code works with external objects(hid by interfaces), it's a very good idea to generate [mocks](https://en.wikipedia.org/wiki/Mock_object).    
We suggest using [gomock](https://github.com/golang/mock) as tool for generating mocks.

We strongly advise to follow next rules while using mocks:
- Use one shared controller for all mock objects - it's possible with `gomock` and gives more opportunities)
- Define order of calls - changing the order of calls usually produces bugs
- Use `Do` and `DoAndReturn` functions to do extra assertions or to get generated values (e.g generated ID)
- Use gomock within different goroutines with caution
