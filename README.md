# Sample Rest

This project is intended to illustrate good engineering practices for golang applications. See [GoEngineering](docs/GoEngineering.md) for more information.

### Build, Run, Test
 #### Simple Run
The service can be built and run in one command from the project directory

`> go run cmd/petserver`

This starts a server listening on port 4852. The port can be configured using the `-p <port>` flag.

`> go run cmd/petserver -p 9412`

### Installation

Install petserver on your local machine using

`> go install cmd/petserver`

Run with

`> petserver <args>`
