# Engineering practices

This project illustrates how we incorporate good practices when working in go.
 * Testing (testing practices, code coverage)
 * Linting
 * Project structure as in [standards-project-layout](https://github.com/golang-standards/project-layout)
 * CI
 * Documentation
   * README.md (build, run, test)
   * API specs

### Tools
 * **golangci-lint**. Meta-linter that can be integrated into ci tools. Highly configurable via command line args or in json/yaml/toml.

 * **go modules**. Package dependency management introduced in go 1.11

### Libraries

It is **Highly** recommended to use the following libraries across all applicable projects.

 * **testify**. Provides better test organization, test suites, test setup/teardown.
 github.com/stretchr/testify
 `go get github.com/stretchr/testify`

 * **kingpin**. Better command line argument parsing.
 github.com/alecthomas/kingpin
 `go get gopkg.in/alecthomas/kingpin.v2`

 * **chi**. For defining application and api endpoints. This library is specific to http server based applications. The `render` library also provides nice http request/response managing.
 github.com/go-chi/chi
 `go get github.com/go-chi/chi`
 `go get github.com/go-chi/render`

 * **logrus**. Great logging solution
 `go get github.com/sirupsen/logrus`

