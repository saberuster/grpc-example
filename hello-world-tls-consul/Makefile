

genproto:
	protoc -I pb/ pb/*.proto --go_out=plugins=grpc:pb

genssl:
	openssl req -x509 -nodes -newkey rsa:2048 -keyout private.key -out public.pem -days 3650 -subj "/CN=*"

runconsul:
	docker run --name=dev-consul1.4 --rm -e CONSUL_BIND_INTERFACE=eth0 -p8500:8500 consul:1.4.0-rc1