# simple-kvdb
Very basic Key-Value store with REST API and UI based on BoldDB

## Build

A simple build script which pulls the dependency and builds the binary
```bash
go run build.go
```

## REST API Usage

### List buckets

```bash
curl -X GET http://localhost:8081/api/

[{"bucket":"bucket1"}, {"bucket":"bucket2"}]
```

### Create bucket

```bash
curl -X POST http://localhost:8081/api/bucket3
```

### Delete bucket

```bash
curl -X DELETE http://localhost:8081/api/bucket3
```

### List Key-Value pairs

```bash
curl -X GET http://localhost:8081/api/bucket1

[{"key":"key1","value":"value1"},{"key":"key2","value":"value2"}]
```

### Create Key-Value pair

```bash
curl -X POST http://localhost:8081/api/bucket1/key1/value1
```

### Delete Key-Value pair

```bash
curl -X DELETE http://localhost:8081/api/bucket1/key1
```

### Get value

```bash
curl -X GET http://localhost:8081/api/bucket1/key1

{"value":"value1"}
```

