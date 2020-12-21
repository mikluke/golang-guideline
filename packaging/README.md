# Packaging

## Index

- [Repository naming](#repos-naming)
- [Package naming](#package-naming)
- [Project structure](#project-structure)

## Topics

### Repo's naming

We suggest three ways for repository naming:

- by domain
- by purpose
- by branded name (e.g ryuk, caesar, mars)

#### By domain

Preferable services naming with using domain driven design

##### Examples

- users
- tokens
- profiles
- notifications
- lobbies

#### By purpose

Use a main purpose of application as a name

##### Examples

- users-provider
- authenticator
- profiles-manager
- notificator
- matchmaker

#### By branded name

If you have project with branded name this way is only possible to name your repo. A main area to use this way is opensource

##### Examples

- cerberus
- рeimdallr
- alfred
- bing
- cross

#### NB

We recommend not to mesh `by-domain` and `by-purpose` ways for a one project.

---

### Package naming

Package name should present certain layer or domain which is present in package.

Use singular nouns as packages' names. Using plural names is acceptable if a package presents plural entity as singular (e.g `bytes`).

##### Examples

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```
- users (package presents `user` model)
- handler (package provides all server layer's functionality)
- header (package presents headers struct)
```

</td><td>

```
- user
- server
- headers
```

</td></tr>
</tbody></table>

---

### Project structure

This project structure was inspired by and based on [this project](https://github.com/golang-standards/project-layout), but we added some changes.

You may adapt this layout for your purposes, but if you want follow this guideline try to ask us before deviation.

*ToDo: add boilerplate which will be an example*

#### High-level description

There are directories which we can have on top level of our hierarchy.

##### assets

Assets to go along with your repository (images, logos, etc)

##### build

Although this directory is not commitable, we define it as directory which contains built files.

##### cmd

Main applications for this project.

The directory name for each application should match the name of the executable you want to have (e.g., /cmd/myapp).

##### configs

Put configuration files, its templates, or default configs here. For example: confd, consul templates, prototool or linter config.

##### docs

Design and user documents (in addition to your godoc generated documentation).

##### internal

Private application and library code. This is the code you don't want others importing in their applications or libraries. 

##### pkg

Code that may be used by external applications.

Put here shared clients, models and other stuff.

##### scripts

Put any(ci, migrations, analysis) scripts here.

##### static

Put there all static files which should be copied to container or build bundle.

##### tools

If you develop additional tools for your project put your executable packages here.

##### vendor

Result of running `go mod vendor`. You may commit this directory if you want(for example to reduce CI stage execution time).

#### Internal directory structure

We also describe how the internal directory should be structured. You may follow this rules if you want your code to be standardized and cozy.

##### server

Put your server's interface implementation packages here.

Server package's code uses own models or imported from pkg package(if you defined client there)

NB: Code placed in this directory shouldn't have any business logic, just (de)serialization, validation, and running service methods code.

Example:

- /internal/server/grpc
- /internal/server/nsq (msg brocker consumber)
- /internal/server/http
- /internal/server/udp

###### entity

A package to define service(repo)-wide models. You may use subpackages, if you have lots of models and you want to divide it logically, but remember, if you use subpackeges don't put files into the root directory.

Examples:

- /internal/entity/user/
- /internal/entity/token/

###### service

This package should contain all business logic. Your service ought be already devided by [DDD](https://en.wikipedia.org/wiki/Domain-driven_design) principles, but if you need to serve two logical domains you may use subpackages.

Service package's code uses `/internal/entities` models.

###### storage

This package contains common storage interfaces and its implementations.

If you have to split this package to subpackages it's okay. Just remember one rule: only implementation packages can be placed near an interface file.

NB: You're able to generate interfaces' mocks in `mock` directory(as implementation).

NB: Implementation package should be named as technology called.

Example:

```
internal
|- storage
    |- storage.go // interface
    |- mock
        |- mock.go
    |- postgres
         |- postgres.go // implementation
         |- postgres_test.go // integration tests
```

Example:
```
internal
|- storage
    |- user
        |- user.go // interface
        |- mock
            |- mock.go // mock implementation
        |- postgres
            |- postgres.go // implementation
            |- postgres_test.go // integration tests
    |- token
        |- token.go // interface
        |- mock
            |- mock.go // mock implementation
        |- mongo
            |- mongo.go // implementation
            |- mongo_test.go // integration tests
```

###### producer

This package should look like storage package: interface and its implementations.

Example:
```
internal
|- producer
    |- producer.go // interface
    |- mock
        |- mock.go // mock implementation
    |- nsq
        |- nsq.go // nsq implementation
        |- nsq_test.go // integration tests
    |- kafka
        |- kafka.go // kafka implementation
        |- kafka_test.go // integration tests
    |- fan_out // for example we push event to all msg brokers
        |- fan_out.go
        |- fan_out_test.go // we should use there mock implementation
```

###### x
Use this package as parent of all your extenstions.

NB: Don't forget to define interface and cover your code with tests.

Example:
```
internal
|- x
    |- health
    |- email
    |- encrypter
```

#### Pkg directory structure

Put packages and files however as you want to see them used. Think about your users.

Example:
```
pkg
|- client
    |- request.go
    |- response.go
    |- client.go
    |- client_test.go
```