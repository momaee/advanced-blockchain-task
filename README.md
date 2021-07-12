## Running the app

```bash
# install and start DDB
$ docker-compose up -d

# install go and its modules and gRPC-gateway

# build protos
$ make

# run project
$ go run . 
or 
$ ./micro

#help
$ go run . -h
or 
$ ./micro -h

# rest api port
8080

# gRPC api port
50051

# api url
POST http://localhost:8080/backend_task/supplies
Body example:
{
    "contractAddress":"0x70e36f6BF80a52b3B46b3aF8e106CC0ed743e8e4"
}
