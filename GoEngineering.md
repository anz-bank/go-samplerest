# Engineering practices

This project illustrates how we incorporate good practices when working in go.
 * Testing
   * Testing practices
   * Code coverage
 * Documentation
   * README.md (build, run, test)
   * API specs
 * Project structure as in the [golang-standards project layout](https://github.com/golang-standards/project-layout)
 * Linting
 * CI

### Tools
 * [golangci-lint](https://github.com/golangci/golangci-lint). Meta-linter that can be integrated into ci tools. Highly configurable via command line args or in json/yaml/toml.

 * go modules. Package dependency management introduced in go 1.11

### Libraries

It is **Highly** recommended to use the following libraries across all applicable projects.

 * [testify](https://github.com/stretchr/testify). Provides better test organization, test suites, test setup/teardown.

 * [kingpin](https://github.com/alecthomas/kingpin). Better command line argument parsing.

 * [go-chi/chi](https://github.com/go-chi/chi). For defining application and api endpoints. This library is specific to http server based applications.
   * The [go-chi/render](https://github.com/go-chi/render) library also provides nice http request/response managing.

 * [logrus](https://github.com/sirupsen/logrus). Great logging solution
 `go get github.com/sirupsen/logrus`
