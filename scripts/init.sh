rm -rf /tmp/mify_tmp
go run ./cmd/mify/ init mify_tmp -p /tmp/
go run ./cmd/mify/ add service service1 -p /tmp/mify_tmp
tree -a /tmp/mify_tmp/