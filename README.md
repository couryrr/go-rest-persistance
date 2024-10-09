# go-rest-persistance
This is a template to spin up a micro service REST api with persistence in go-lang

## Notes

For the use of micro-service there should not be a need for a lot directories. All files are stored in the root except for internal types which are stored in the `internal` package.

## Setup

### HTTP

The goal is to use the standard library `http/net`.

### Persistence

For simplicity sake the template is using SQLite

