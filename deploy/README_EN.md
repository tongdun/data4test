##### System Requirements
- Linux/macOS/Windows system
- MySQL database

##### Production Environment
- Download the release package
- Unpack the tar file: tar -xvf data4test_20XX0X0X.tgz
- Navigate to the file directory: cd deploy
- The default package provides the Linux x86 version. If other versions are needed, download the respective package and replace it.
- Download the required environment execution files to the ./deploy directory
- Rename the package: mv data4test_xxx data4test

##### Deployment in a New Environment
- Create a database: create database data4test;
- Import the initialization database file: mysql -h x.x.x.x -u user -p data4test < ./mgmt/sql/init.sql
- Modify the configuration file: config.json, fill in all placeholders according to the actual situation
- Edit config.json: Update database and path information as needed
- The default listening port is 9088. If there is a conflict, it can be changed according to the actual situation.

##### Update in an Existing Environment
- Change log: ./mgmt/doc/file/update/change_log.md
- If there are scheduled tasks, they do not need to be executed in the new environment. You can pause running tasks. Once the tasks are redeployed, they will automatically restart.
- Import the latest updated SQL: mysql -h x.x.x.x -u user -p data4test < ./mgmt/sql/update.sql (Pull the corresponding SQL from the log for updating)
- Terminate existing processes: ps aux | grep -i data4test | head -n 1 | awk '{print $2}' | xargs kill -9

##### Deployment System
- nohup ./data4test >/dev/null 2>&1 & # No logging to nohup.out
- nohup ./data4test & # Automatically logs to nohup.out

##### Access and Login
- Access: http://10.0.X.X:9088
- Default user: admin/ admin