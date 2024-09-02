### API Management

#### API Changes and Tracking

1. **Specification Check for New API Definitions and Change Detection**
    - **Specification Check (with continually added checkpoints based on API specifications)**
        - Verify if newly added APIs return unique identifiers
        - Ensure that query data operations return `data` content
        - Check for incorrect parameter placement, e.g., including Body parameters in GET requests or Query parameters in POST requests will fail
    - **Change Detection**
        - Application Level: Count of newly added APIs, deleted APIs, modified APIs, and unchanged APIs
        - API Level: Addition of new fields/parameters, deletion of fields/parameters, changes in field/parameter types

2. **Specification and Change Checks During API Changes via API Console or Dedicated Management Page**

3. **History Query Support for Specification Check Results**

4. **History Query Support for Change Detection Results**

#### API Data Reuse

- **Select an API to Automatically Generate Test Data**
    - Standard data files in YAML or JSON format for testing purposes.