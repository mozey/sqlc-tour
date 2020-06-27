[Introducing sqlc](https://conroy.org/introducing-sqlc)

The repos mostly follows the tutorial above, 
but with a bit more structure 

Create and start `postgres` docker container

    docker pull "postgres:12.3"

    docker create --name postgres \
    -p 5432:5432 \
    -e POSTGRES_HOST_AUTH_METHOD="trust" \
    -e DB_HOST=docker.for.mac.host.internal \
    postgres

    docker start postgres
    

Confirm it's working with `pgcli`
    
    pip3 install pgcli
    
    pgcli postgresql://postgres@localhost:5432


Generate db code and run it

    brew install kyleconroy/sqlc/sqlc
    
    sqlc generate
    
    go run *.go
    
    
**TODO** Bulk inserts
[How does this work for bulk operations?](https://github.com/kyleconroy/sqlc/issues/216)
[Efficient bulk imports via `pq.CopyIn`](https://github.com/kyleconroy/sqlc/issues/218)

