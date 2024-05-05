# Accel-One Contacts API (interview exercise)

## General

This repository aims to provide a scaffold and simple implementation of a Contacts API.

Build & Run like this (port is optional):
```shell
go build
./accelone-contacts --port 8000
```
or
```shell
go run main.go
```

Generate random contacts to test the API with:
```shell
./generate-contacts.sh
```

## Features

- Logging
- JSON error responses
- Pagination for Get all
- Concurrency handling

## Design

- Modular approach, easy to integrate new APIs and reuse config such as Json Content-type
- Easy to extend interfaces such as if a database is used for persistence as opposed to in-memory.
- Used generic package for scalability to encompass different services,models and apis to achieve separation by tiers (business layer, data layer, model layer).

## Considerations

- Does not implement uniqueness in contact attributes besides `Id`.
- Does not index by any other attribute besides `Id`.
- Pagination is suboptimal due to map needed to be traversed before slicing the results.
  > - This limitation is due to the in-memory approach with would require manual bookkeeping of items as list in parallel to the map.
  > - In databases sequential indexes can be used to skip records as needed