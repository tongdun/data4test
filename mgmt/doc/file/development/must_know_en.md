##### Code Integration Notice
```
1. The changes made must be recorded in the change_log.md file, including what problems they solve. The format for recording is as follows: Change Log
2. If there are SQL changes, they must be appended to the update.sql file. The format for recording is as follows: SQL Changes
3. If new features are added, instructions on how to use them must be provided in the corresponding documentation, along with example explanations.
   - (1) For changes to actions, record them in the action_design.md file. If there are existing instructions, you can use them as a reference. For templates, add instructions to the template_instructions.yml file as needed.
   - (2) For changes to assertions, record them in the assert_design.md file. If there are existing instructions, you can use them as a reference. For templates, add instructions to the template_instructions.yml file as needed.
   - (3) For changes to parameters, record them in the parameter_design.md file. If there are existing instructions, you can use them as a reference. For templates, add instructions to the template_instructions.yml file as needed.
4. Please conduct self-verification before integrating new code. If self-verification is not conducted, please explain in the change log.
5. If the new code is a temporary solution, please explain in the change log.
6. When adding new functions, please name them according to their actual functions and follow industry naming conventions.
```

##### Change Log
```
Date: 20XX-XX-XX
Number: [Type] #Issue or Requirement or Optimization Number XXXX

Type:
1. Bug: Represents a feature that has been integrated with code-related issues
2. Optimize: Represents a feature that has been integrated with improvements in usability or logic
3. Feature: Represents a new feature that has been integrated

Description: #Issue or Requirement or Optimization Number Fill in the actual number or omit if there is no relevant number
```

Example:
```
##### Date: 20XX-XX-XX
1. [Bug] #5459679 The scene testing history environment type in the console-domain running scene is incorrect due to an assignment error
2. [Optimize] #6423134 When editing associated data files in the scene list, single-line editing is now supported, and data files can be moved up and down
3. [Feature] #4654561 The task list now supports关联数据编辑 with single-line editing and up-down movement for associated scenes
```


##### SQL Changes
```
-- Date: 20XX-XX-XX
SQL Statement 1;
SQL Statement 2;
```

Example:
```
-- Date: 2021-09-01
alter table scene_data change yaml file_name longtext null comment '文件名称';
alter table scene_data change json content longtext null comment '文件内容';
```