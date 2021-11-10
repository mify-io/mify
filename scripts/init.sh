rm -rf /tmp/mify_tmp
go run ./cmd/mify/ init mify_tmp -p /tmp/
go run ./cmd/mify/ add service service1 -p /tmp/mify_tmp
go run ./cmd/mify/ generate service1 -p /tmp/mify_tmp
(cd /tmp/mify_tmp/go_services && go get ./...)
tree -a /tmp/mify_tmp/