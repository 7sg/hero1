#What would i test and how?

I will test all the function which are part of business logic.
For example, i will test service layer, repository layer and domain layer. Because these are the layers which will change frequently in future.
Each `.go` file in these layer will have file with suffix `_test.go`.
Every test file is intended to test the functions of source code file.
I have created `_test.go` files with empty test functions in these layers.
I will use table driven test, where every test function will have list of test cases, and that can be run in a for loop.


#Repository layer test
Since our test required interaction with db, we can start a mongo db instance in docker-compose, then run the tests on top of it.
all the steps can be automated in a script file or make file.
I do not prefer mocking of db calls here, because that will not give me confidence whether the library/integration i have done with db is correct or not.

#Service layer test
Same goes for service layer test as well, i will first start db, then run all test cases of service layer.
I will test all the positive and negative cases. In idle case every line of code should be covered.



#Additional test(API tests)
Go provides framework to test each rest endpoints of the services, i will prefer to write these as well after above tests.
We can start all dependencies first, which is only db in our case, in a container.


#Thoughts on future extension of service/logic
In case business logic become complex
   - If it is related to event(single responsibility), then i will add functions in event package,
   - otherwise i will create new services and repository, which will have its own packages.

Future business logic addition will grow domain, service and repository layers.
   - I will identify new storage level entities, and create respective domain , repositories and services in respective packages.


#Thoughts on managing huge number of activities

Currently we are fetching the events from db, with out any caching.
###Scaling the read traffic
1. Currenlty i am using pattern match to do text search. Instead of this we can create text index on `message` field, which can support text search.
   and it is more efficient then pattern match.
   
2. we can cached some data with respective filers and ttl, with the assumption that client might get stale data for sometime.
   once the ttl has been passed, on new request call be made to db and new set of events will be fetched and updated in cache.
   something like redis will work here.
   
3. Read from the slaves of mongodb, but there might be some lag with master.


###Scaling the write traffic
1. Using the multi-masters approach, we can scale the write traffic. Like we can have many masters which is responsilbe 
   for some part of the data(Sharding). and each master node will have its own set of replicas.
  
2. Async write, with the assumption that we can show the updated data after sometime, but not immediately.
   In this case, we can use a message bus(something like kafka). we first write the event to bus and then asyncronsly write the events to db from bus.
   
  
#How to run this service
This service supports both `GRPC` and `HTTP` servers.
GRPC server runs on 9999 and HTTP server runs on 8080.
This service use mongodb as storage to store the event/activity.
I have provided Makefile in root directory, which can be used for generating stubs from proto message files.
and this make file can be used for quickly starting and stopping the service with its dependency(mongodb).
###To start service
`make start`

###To stop service
`make stop`

#Swagger UI
Swagger ui is automatically updated when we run `make generate`, it generates all stubs and swagger files.
you can access swagger at `http://localhost:8080/swagger-ui/`
###Sample Request
Search events with out any filter,

`curl -X GET "http://localhost:8080/v1/event" -H "accept: application/json"`

Search events with filters,

`curl -X GET "http://localhost:8080/v1/event?email=diego%40hero1.com&environment=production&component=orders&text=success&date=2018-05-12T00%3A00%3A00Z" -H "accept: application/json"`

Save event

`curl -X POST \
   http://localhost:8080/v1/event \
   -H 'Content-Type: application/json' \
   -d '{
   "created_at":1526123095,
   "email": "diego@hero1.com",
   "environment": "production",
   "component": "orders",
   "message": "the buyer #123456 has placed an order successfully",
   "data": "{ \"order_id\": 123, \"amount\":300, \"created_at\":1526123095}"
 }'`


#Structure of the service
`hero1/.githooks` will contain commit message format, pre commit checks etc.

`hero1/.github` will contain github workflow, for continuous integration and deployments.

`hero1/docs` will contain architecture, code, and db schema related documents etc.

`hero1/go-services` will contain go services. Currently it has one `historical-events`

`go-services/historical-events/api/proto` will contain proto message files, generated packages and files.

`go-services/historical-events/api/swagger-ui` will contain swagger ui.

`go-services/historical-events/internal` has packages for configuration, domain, infra(grpc/http server, db), repository and service.

`go-services/historical-events/Dockerfile` docker file is for `historical-events` service.

`hero1/infra` has service environments variables and docker compose for running service locally.
