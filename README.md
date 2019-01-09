# Sample Rest

This project is intended to illustrate good engineering practices for golang applications.

### Third party library

It is **Highly** recommended to use the following libraries across all applications.

 * **testify**. Provides better test organization, test suites, test setup/teardown.
 github.com/stretchr/testify
 `go get github.com/stretchr/testify`

 * **kingpin**. Better command line argument parsing.
 github.com/alecthomas/kingpin
 `go get gopkg.in/alecthomas/kingpin.v2`

 * **chi**. For defining application and api endpoints. The `render` library also provides nice http request/response managing.
 github.com/go-chi/chi
 `go get github.com/go-chi/chi`
 `go get github.com/go-chi/render`

 * **logrus**. Great logging solution
 `go get github.com/sirupsen/logrus`

 *