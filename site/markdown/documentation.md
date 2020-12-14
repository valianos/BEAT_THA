# BEAT Take Home Assignment

This document serves as a brief description on the implemented microservice.

## High Level Introduction

The implemented microservice exposes a single API path. When server runs, 
`POST` requests made to `{server}:8080/calculate` will trigger a process to use and expose
the already existing and encapsulated by this service, ETA microservices.

The protocol dictates that requests must bare the following `Calculate` entities to get validated.

### Point entity

| Property     | Type                 | Description |
| ----------   |:-------------------  | :-------------------------------   |
| lat          | number(float)        | The latitude of the point.         |
| lng          | number(float)        | The longitude point of the route.  |

### Calculate entity

| Property          | Type                 | Description |
| -------------     |:-------------------  | :-------------------------------   |
| origin            | Point                | The starting point of the route.   |
| destination       | Point                | The end point of the route.        |
| provider          | string               | The preferred microservice to use. |

_example_

    {
        "origin":{
            "lat":37.0816818,
            "lng":23.5035676
        },
        "destination":{
            "lat":37.5142881,
            "lng":22.5093049
        },
        "provider": "ETAServiceA"
    }

If everything is proper, the service will respond with the following entity.

### Service Response Entity

| Property          | Type                 | Description |
| -------------     |:-------------------  | :-------------------------------                |
| eta               | Number               | The ETA returned from the chosen microservice.  |
| provider          | string               | The chosen microservice.                        |

### Possible failures

 * A status `405` will get responded if other method than `POST` is used.
 * A status `400` will get responded if request has invalid or missing body.
 * A status `500` will get responded if encapsulated microservices fail to properly respond.

## Key components of the design

The project has been broken-down into two modules, namely the `protocol` and the `etaMicroservice`

The `protocol` module contains the structs and functions that describe the protocols to and from 
the user as well as to and from the microservices. It also exposes functions for deserializing
protocol objects. 

The `etaMicroservice` is the module containing the whole implemented functionality. This consists 
of the three following packages, each of which serves a certain purpose.

### Logger

Every decent application needs a proper logger. A semantically useful format is always helpful
when hunting down issues. Here, a basic logger has been implemented just to illustrate this need.

### HttpUtil

This package contains the basic functions for performing HTTP requests.

### Main

The main package of the application consists of the main file that has the main go file, 
which starts the aforementioned endpoint. It starts the server and defines the handler 
method that will handle and serve all incoming requests. In this scope, the external service 
are used as an abstraction.

This is achieved with the use of the `externalMicroservice` go file, which encapsulates the 
whole logic of the external microservices. There a factory method is exposed for abstractly 
creating the desired microservice(s). Then a method for properly using each of them is defined.
Finally, it completely hides the external protocols from the main function, it only exposes 
the inner protocol. 

This way, the main functionality is completely agnostic of the underlying microservices. This 
enables the application to be more flexible in implementing new endpoints or even deprecating 
old ones etc.

In addition, a constants file is defined, where useful constants, globally required, are 
gathered and defined. 

Finally, there lies a basic functionality unit test that assures that the endpoint runs according 
to given specifications.

## Patterns used

### Factory method

As stated already, for creating the desired external microservice, a factory method was exposed. 
It makes possible to create everything needed for using the appropriate endpoint just be deserializing 
the received `Calculate` entity. This hides the complexity and the implementation specifics from the main
function.

### Singletons

The singleton pattern was chosen for implementing the external microservices as this is a design pattern 
that restricts the implementation specifics of them to only one file, while maintaining a global point 
of access to them. In addition, it only gets constructed once. For this purpose the function init was 
chosen.

### Interfaces

By defining the general interface of an external microservice, it was possible to hide implementation 
specifics of each of them. By defining a function that will use the implementations of this interface,
the abstraction for the main function was complete.

### Channels

Channels were first imported in the project as a means to sync data that were passed to several go routines,
when implementation went to use such logic.
The most essential purpose though was that it made it fairly easy to select the quickest response between 
the two external services, when this is required.

## Future work

The test needs more work to be done in order to be a full test.

 * The external microservices should be mocked, it is not right to depend on any external service, and it 
   will give full control to implement more test cases by giving access to the responses returned.
 * Right now the only functionality that is being tested is that if the services are up-and-running, the
   microservice will operate as expected. We should test what will happen if the services are own or take 
   too long to respond.
   
In the main file, a more clear pipeline would be a great next improvement. This would be beneficiary for a
real-life application as it would make efficient use of I/O and multiple CPUs.