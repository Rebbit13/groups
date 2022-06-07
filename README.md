# Groups REST API
A simple go rest api based on gin and gorm. 
For migrations uses rubenv/sql-migrate. 
For database uses PostgreSQL.
For API documentation swagger and swaggo/swag.
For logging uses zap logger.

### Group entity
Group has name and subgroups (other group).
Can have members (humans)

### Human entity
Human has name, surname, birthdate and can 
be a member of some groups

# Run
---
**NOTE**

You need docker and docker-compose for running project

---
Go to the app dir and run the command:

> cp ./example.env ./.env

### Test data
If you need test data you can set env var in .env
TEST_DATA to value true
> TEST_DATA=true

Change environment vars if you need and then:

> docker-compose up

# Author
Reshetnikov Evgeniy | Rebbit13
* Telegran @rebbit13
* email errb13@gmail.com

# Routes
You can see swagger doc on the path "/v1/swagger/index.html"

### Regenerate swagger doc
Use this command in the source project dir if you need regenerate swagger doc:
> swag init -g cmd/main.go 