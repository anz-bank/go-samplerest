# Sample Rest

This project is intended to illustrate good engineering practices for golang applications. See [GoEngineering](docs/GoEngineering.md) for more information.

### Build, Run, Test
#### Simple Run
The service can be built and run in one command from the project directory

`> go run cmd/petserver`

This starts a server listening on port 4852. The port can be configured using the `-p <port>` flag.

`> go run cmd/petserver -p 9412`

#### Instal

Install petserver on your local machine using

`> go install cmd/petserver`

Run with

`> petserver <args>`

#### Test

To run unit tests, use the command
`> go test ./...`

##### Coverage
The same command can be used to generate a code coverage report using some options
`> go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out`

This will create a file `coverage.out` that contains information about code coverage (in unreadable form)
The second command opens a browser and displays coverage information.
