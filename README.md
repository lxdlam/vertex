# Vertex

Vertex is a pure go memory storage engine which is compatible with RESP. **Now working in progress**.

## Usage

No binary or package available by now.

Clone the repo then start `cmd/server/main.go` and `cmd/client/main.go` manually.

## Status

It's at a really early stage of development. Currently supported feature:

- TCP connection and full RESP support.
- String, List, Hash and Set and a subset of the core commands are supported.

## Limitations

- The whole system is built above the GC of go.
- Performance may poor now since no benchmark has been performed.
- The hash and set are just simple proxies to go's map.
- Only one database is supported.
- RESP is supported, but do not fit the real communication environment what means you may not use any redis client to communicate by now.
- And more...

## Road map

Coming soon...
