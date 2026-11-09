#### Development Environment

#### Runtime Environment Requirements
- Go version: the go.mod `go` directive requires ≥ 1.23; using the Go 1.27.3 toolchain is recommended for building.
- The new Go toolchain is installed at ~/sdk/go1.27.3; run `unset GOROOT` before use (to avoid the old GOROOT in the shell profile pointing to 1.17, which causes compile errors).
- On China networks, set `GOSUMDB=off` when running `go mod tidy` (GOPROXY already points to the Alibaba Cloud mirror).
- linux/arm64 cross-packaging requires Go ≥ 1.25 (below this version it errors with "ARM external linker must be gold").

#### Console Frontend
- Command1: yarn install
- Command2: yarn run - select dev mode

#### Backend
##### Import MySQL
- Full mode
  - [InitAllSQL](../../../sql/init_all_XXX.sql)
- Incremental mode
  - [InitSQL](../../../sql/init.sql)
  - [UpdateSQL](../../sql/update.sql)

##### Update Configuration File
- [Configuration File](../../../../config.json) Fill in the information based on actual situation.

##### Code to Start the Service
- Command: go run main.go / sudo go run main.go

##### Login
- Default Access: http://127.0.0.1:9088
- Default User: admin/ admin

#### Parameter Changes
- To use intelligent testing, go to [Environment - System Parameters] and configure the "aiRunEngine" parameter according to the actual LLM.
- To execute non-standard data, go to [Environment - System Parameters] and configure the "scriptRunEngine" parameter according to the actual execution engine.

##### Others
- If using Xmind to import use cases for product lists, the environment needs to install xmind2case first.
- Install xmind2case: https://pypi.org/project/xmind2case/
- Use a relatively old version of Xmind software to save xmind files, as the latest xmind2case cannot parse them.
