# TestContainers.

[TestContainers](https://testcontainers.com/) is an opensource library that provided containerized environments to test your code.
The advantage of using containerized environments is that they can be set up and deleted without affecting your 
system and also can be run anywhere. As demonstrated in this code, they can be run in GitHub actions. 

This application is a CRUD Go application that deals with a single table in a MySQL database.
When tests are run, a Mysql Docker container is set up and the tests run inside it.

> rapando
