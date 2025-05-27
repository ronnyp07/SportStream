nats str add SPORTSTREAM --server=nats://nats:4222 --config ./streams/sportstream.json

nats con add SPORTSTREAM sportstream_docker_updated --server=nats://nats:4222 --config ./consumers/sportstream_docker_updated.json