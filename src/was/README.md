# Gate

Gate is the combination server of various protocol, which supports HTTP, WebSocket, gRPC, TCP and udp. 
Probably will add more from time to time. As universe interface all access will be through Gate service.

## HTTP
- /ping : Health checking with 'pong'
- /air/{city} : Retrieve air quality for external site.
- /metrics : Instrumenting for Prometheus scrap

## gRPC

## graphQL
- / : Playground for graphQL, run query, mutation, etc.

## WebSocket
- / : Embedding page for WebSocket demo.
- ws://${host}/echo : WebSocket based echo service.
- ws://${host}/spy : WebSocket for realtime Kubernetes information.
  
## TCP
- : Telnet to listen port and experience TCP based echo service (delimit '\n').

## UDP 