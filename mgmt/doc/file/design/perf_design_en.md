### Performance Testing

##### Objective
- Obtain the system's concurrent data support under specific environments to guide users in environment planning and anticipate system capabilities.

##### Implementation
- Acquisition of Environment Configuration Information:
    - Manual Input
    - Automatic Retrieval
    - High Concurrency Testing for Single Interface
    - Concurrent Testing for Multiple Interfaces Based on Weights or Ratios

##### Test Types
- Baseline Testing: Acquisition and Establishment of Baseline Data (Primary Focus Initially)
- Load Testing: Gradually increasing pressure to consume resources up to a certain percentage, observing system speed and stability.
- Stress Testing: Pushing the environment to its limits, consuming resources above 99%, and observing system speed and stability.
- Stability Testing: Determining the amount of data under which the system's speed and stability are guaranteed.
- Concurrency Testing: 1k / 10k / 100k / ...

#### Plan
##### Preliminary Data Preparation
- Information of the Tested Environment
    - Host: IP / HTTP / HTTPS
    - Authentication: Username / Password / Token
- Environment Information Collection
    - CPU
    - Memory (Mem)
    - Disk
    - Network
- Interface Concurrent Request Data
    - Response Time
    - Concurrency Number
    - Transactions Per Second (TPS)
    - Error Rate
- Automatic Collection of Performance Testing Data, with Automatic Backup of JMeter Test Results
- Automated Generation of Test Reports, Including the Necessary Data Collected Automatically