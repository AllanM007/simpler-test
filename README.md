### Backend Engineer Take Home Test @Simpler

- This is a simple CRUD REST API written in Go for a product resource microservice with unit tests and full test coverage.

- To run the project you have the option of running the docker compose file to set up both the api server and it's database with sample data and run the project locally.
  ```
  docker compose up -d
  ```

- This will spin up two docker containers one of which is the golang api and the other is the postgres database with some sample data for testing

<h3>Documentation</h3>
- The project has swagger docs generated from the rest endpoints that provide better context on the implementation of each api which can be accessed from:
  
  ```
  localhost:8080/api/swagger/index.html
  ```