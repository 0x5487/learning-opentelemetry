# learning jaegertacing

1.  setup jaegertacing
```
docker run --rm -p 6831:6831/udp -p 6832:6832/udp -p 16686:16686 jaegertracing/all-in-one --log-level=debug
```
1. go to dashboard `http://localhost:16686`