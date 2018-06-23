# Rider

Simple Microservice backend application

Rider is composed for 3 microservices:

    - The Bike Service which handles a fleet of bikes,
    
    - The Trip Service abstracts the Trip Workflow,
    
    - The API Gateway (the only one that is been exposed)

## How to install

`git clone git@github.com/EarvinKayonga/rider rider` to checkout this baby

`cd rider`

`docker-compose up --build -d` to run the gateway on port 8080

## How to use


- GET `/health`                                 health check       
- GET `/bikes?cursor={cursor}&limit={limit}`    list bikes
- GET `/bike/{bikeID}`                          bike description

- POST `/trip/track`: add a point location to a trip

```
        { 
            "trip_id": string,
            "location": {
                "lng": int,
                "lat": int
            }
        }
```

- POST `/trip/start`: starts a trip

```
        { 
            "bike_id": string,
            "location": {
                "lng": int,
                "lat": int
            }
        }
```

- POST `/trip/end`: ends a trip

```
        { 
            "trip_id": string,
            "location": {
                "lng": int,
                "lat": int
            }
        }
```

## Observations

Only the happy path is implemented. There is no implementation of error handling 
across the microservices boundaries which needs its own layer.
The Makefile and Dockerfile are kept very simple. And the docker images are small (the binaries are stripped).
