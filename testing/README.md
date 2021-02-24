# Testing

## Index

- [Why do we use tests](#why-do-we-use-tests?)
- [Packaging](#packaging)
- [Formatting](#formatting)
- [Mocking](#mocking)
- [Table driven tests](#table-driven-tests)
- [Integration tests](#integration-tests)

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
- Transport code
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

You have to agree. Code covering by tests is the complicated business, but the manual testing is **more** complicated. After that you have to show reasons why we should cover our code by tests. In the end, if it won't persuade your manager, just leave this place and get another job with smart management personnel. 

### Packaging

#### _test.go

Put your tests near to code files
```
serivice
| - service.go
| - service_test.go
```

#### Test data
If you need to use assets within tests put it into a `testdata` subdirectory
```
image
| - testdata
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

### Formatting

#### Avoid spaces in test names
Output of `go test` replaces spaces with underscores and this way it will be more complicated to find failed test with `find in directory` IDE function.

#### Do not number tests, use clarification postfix
Use next template for test functions
```
Test{StructName}_{MethodName}_{Clarification}
```
  
!Never name like this
```
TestService_GetUser2
```

#### Define constants
If you use some specific data in several tests, define it before tests.

Example
```
const specificAddr = "0x12345678910"
var requestMetaData = MetaData {
    IP: "1.2.3.4",
    UserAgent: "hello! It's me!",
    Cookie: []Cookie{...}
    Headers: []Header{...},
} 
```

#### Think twice if you need suites
TestSuites can be useful but using this can make your code more complicated. Using [table driven tests](#table-driven-tests) should be enough at almost all cases.

#### Use testify package
Use testify [assert](https://github.com/stretchr/testify/tree/master/assert) and [require](https://github.com/stretchr/testify/tree/master/require) packages. It will make your tests more readable.

### Mocking

If your code works with external objects(hid by interfaces), it's a very good idea to generate [mocks](https://en.wikipedia.org/wiki/Mock_object).    
We suggest using [gomock](https://github.com/golang/mock) as tool for generating mocks.

We strongly advise to follow next rules while using mocks:
- Use one shared controller for all mock objects - it's possible with `gomock` and gives more opportunities
- Define order of calls - changing the order of calls usually produces bugs
- Use `Do` and `DoAndReturn` functions to do extra assertions or to get generated values (e.g. generated ID)
- Use gomock within different goroutines with caution

### Table driven tests

At almost all cases [table driven tests](https://github.com/golang/go/wiki/TableDrivenTests) are very useful.

#### Conditions for using table driven tests

- Test structure(input/output params and execution flow) is identical for every testcase
- There are more than 2-3 testcases

Examples:
- validation tests
- transport tests
- utils tests

#### Parameters

There are 2 good ways to define input parameters: all-params-definition and option-applying.  
Use most applicable to your function way, but you should always define output params for both ways.

##### All parameters definition

Define all parameters in test struct. Use this way when your function has a simple interface.

```
// function under test
func deletePost(ctx context.Context, user uint64, id uint64) error
```

```
// struct for table driven test

tt := []struct{
    // test name
    name string
    
    // input params
    ctx context.Context
    user uint64
    id uint64
    
    // will be returned by mock objects
    storageErr error
    producerErr error
    
    // output values
    err error
} {
    {
        name: "success",
        ctx: context.Background(),
        user: 1,
        id: 2,
        // you can omit next lines
        storageErr: nil,
        producerErr: nil,
        err: nil
    },
}
```

##### Option applying

Define default params object, copy and patch that to use in test (like with option pattern)  
Very useful for functions with huge input structs.

```
// function under test
func createPost(ctx context.Context, post service.Post) error
```

```
// struct for table driven test

type params struct {
    // input params
    ctx context.Context
    user uint64
    id uint64
    title string
}

defaultParams := Params {
    ctx: context.Background()
    user: 1,
    id: 2,
    title: "title",
}

tt := []struct{
    // test name
    name string

    // this func will return parameters for test  
    getParams func() params
        
    // output values
    err error
} {
    {
        name: "success",
        getParams: func() params {
            p := defaultParams
            p.title = ""
            return p
        },
        err: nil
    },
}

for i := range tt {
    tc := tt[i]
    t.Run(tc.name, func(t *testing.T) {
        t.Parallel()
        
        params := tc.getParams()
        
        ...
    })
}
```

### Integration tests

With [docker](https://www.docker.com), it is very easy to set up any service or database and run integration tests as unit tests.

To speed up tests execution you can hide integration tests by build tag.  

Example
```
//+build integration
```

#### Useful packages

Use can use one of these libraries to start containers:

- [dockertest](https://github.com/ory/dockertest)
- [testcontainers](https://github.com/testcontainers/testcontainers-go)

These packages are pretty same, so you can use any of it

#### Wait while service is setting up

Always define a [wait strategy](https://github.com/testcontainers/testcontainers-go/blob/master/wait/wait.go) to make your test determined.

Example
```
req := testcontainers.ContainerRequest{
		Image:        "postgres:12",
		Env:          map[string]string{"POSTGRES_PASSWORD": "root"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
}
```

#### Clean environment

We define two ways of keeping clean environment for integration tests:

- cleaning up after every test (easy to implement, impossible to run parallel)
- run tests in isolated environment (more complicated implementation) 

##### Examples

1) cleanUp func
```
func cleanUp(t *testing.T) {
    _, err := db.ExecContext(ctx, `DELETE FROM post`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `DELETE FROM user`)
	require.NoError(t, err)
}

func Test_GetUser(t *testing.T) {
    defer cleanUp(t)
    
    // test implementation
}
```

2) getDB func
```
func getDB(t *testing.T) {
    return client.Database(fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
}

func Test_SetPost(t *testing.T) {
    s := NewStorage(getDB(t))        
        
    // test implementation

}
```

3) isolated func
```
func isolated(t *testing.T, f func (db *mongo.Database)) {
    db := client.Database(fmt.Sprintf("%s-%d", t.Name(), rand.Int31())
   
    f(db)
   
    if err := db.Drop(context.Background()); err != nil {
        logrus.WithError(err).Fatal("failed to drop database")
    }
}

func Test_SetPost(t *testing.T) {
    isolated(func (db *mongo.Database) {
        s := NewStorage(db)
        
        // test implementation
    })
}
```

#### Example

You can find example [here](example/integration-test/storage_test.go)
