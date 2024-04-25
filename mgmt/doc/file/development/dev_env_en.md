#### Development Environment
#### Console Frontend
- Command1: yarn install
- Command2: yarn run - select dev mode

#### Backend
##### Import MySQL
- [InitSQL](../../../sql/init.sql)
- [UpdateSQL](../../sql/update.sql)

##### Update Configuration File
- [Configuration File](../../../../config.json) Fill in the information based on actual situation.

##### Code to Start the Service
- Command: go run main.go / sudo go run main.go

##### Login
- Default Access: http://127.0.0.1:9088
- Default User: admin/ admin

##### Others
- If using Xmind to import use cases for product lists, the environment needs to install xmind2case first.
- Install xmind2case: https://pypi.org/project/xmind2case/
- Use a relatively old version of Xmind software to save xmind files, as the latest xmind2case cannot parse them.