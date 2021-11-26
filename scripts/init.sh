rm -rf /tmp/mify_tmp
go run ./cmd/mify/ init mify_tmp -p /tmp/
go run ./cmd/mify/ add service service1 service2 -p /tmp/mify_tmp
go run ./cmd/mify/ add client service1 --to service2 -p /tmp/mify_tmp
go run ./cmd/mify/ remove client service1 --to service2 -p /tmp/mify_tmp
go run ./cmd/mify/ add client service1 --to service2 -p /tmp/mify_tmp
go run ./cmd/mify/ add service --lang js front -p /tmp/mify_tmp
go run ./cmd/mify/ add client front --to service1 -p /tmp/mify_tmp
(cd /tmp/mify_tmp/go_services && go mod tidy)
