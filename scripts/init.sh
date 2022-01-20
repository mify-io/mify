TARGET_PATH=$HOME/.cache/mify_tmp
rm -rf $TARGET_PATH
go run ./cmd/mify/ init mify_tmp -p $HOME/.cache
go run ./cmd/mify/ add service service1 service2 -p $TARGET_PATH
# go run ./cmd/mify/ add client service1 --to service2 -p $TARGET_PATH
# go run ./cmd/mify/ remove client service1 --to service2 -p $TARGET_PATH
# go run ./cmd/mify/ add client service1 --to service2 -p $TARGET_PATH
# go run ./cmd/mify/ add frontend --template vue front -p $TARGET_PATH
# go run ./cmd/mify/ add client front --to service1 -p $TARGET_PATH
# (cd $TARGET_PATH/go_services && go mod tidy)
