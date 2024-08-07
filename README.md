# Dataflow API


## Architecture

In my project, I am adhering to clean architecture principles by structuring it into three distinct layers:
1. Handler layer: Handles HTTP-related concerns such as JSON decoding, request validation, and response formatting.
2. Service layer: Contains the business logic of the application
3. Repository level: Responsible for saving and retrieving data

Note: I keep the context (ctx) at the repository level because, in a real-world application that uses a database, it’s important to manage request context. This way, if we change our database implementation in the future, the interface will remain the same, making it easier to adapt without needing to change the rest of the code.


## How to run a server

Using docker:

`docker build -t dataflow-api . && docker run -p 8080:8080 dataflow-api`

Or using docker-compose:

`docker-compose up`


## How to run tests

`make test`

