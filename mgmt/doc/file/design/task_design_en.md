#### Task Types
- Custom Task: Uses a Crontab expression to set the execution frequency as needed.
- One-Time Task: Executes once and stops, can be re-executed repeatedly.
- Daily Task: Sets tasks to execute on a daily basis.
- Weekly Task: Sets tasks to execute based on weekly maintenance schedules.

#### Task Execution Scope
- Data-based Tasks
- Scenario-based Tasks

#### Task Execution Environment
- Can be associated with a single environment: for functional testing, performance testing, etc.
- Can be associated with multiple environments: for multi-environment testing in the information technology and innovation sector, allowing for execution across multiple environments simultaneously.
- Can be executed without associating with an environment: directly using the environment associated with the scenario.

#### One-Click Task Export
- All information used upstream and downstream of the task is exported as SQL.
- Includes information from forms such as Task/Product/Application/Scenario/Data/Assertion Templates/System Parameters, etc.
- SQL Generation Rules: Skip product and application configurations if they exist; insert if not present, update if present.
- The creator of the data is uniformly set to the information of the person performing the export operation, and the product or associated product belongs to the information of the first product in the exported task.

#### Crontab Expression
```
* * * * *
- - - - -
| | | | |
| | | | +-------------- Day of the week (0 - 6) (Sunday is 0)
| | | +---------------- Month (1 - 12)
| | +------------------ Day of the month (1 - 31)
| +---------------------- Hour (0 - 23)
+------------------------- Minute (0 - 59)
```
This Crontab expression allows you to specify when a task should be executed based on minutes, hours, days of the month, months, and days of the week.