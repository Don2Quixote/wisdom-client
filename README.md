# wisdom-client

[Server's code](https://github.com/don2quixote/wisdom-server)

## ENV Configuration:
| name           | type   | description                            |
| ---------------| ------ | -------------------------------------- |
| WOW_HOST       | string | Port to laynch TCP server              |
| MAX_COMPLEXITY | int    | Bytes to read for the challenge number |

## How to launch
Using Makefile:
```
make run
```

Or using docker manually:
```
docker build -f build/Dockerfile -t wisdom-client .
docker run --name wisdom-client --rm --network=host \
    -e WOW_HOST=localhost:4444 \
    -e MAX_COMPLEXITY=16 \
    wisdom-client
```

`MAX_COMPLEXITY` value should match server's one.