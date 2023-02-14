NAME=mify_tmp
TARGET_PATH=$HOME/.cache/"$NAME"
rm -rf "$TARGET_PATH"

go run ./cmd/mify/ init "$NAME" -p $HOME/.cache || exit 2
go run ./cmd/mify/ add service service1 service2 -p "$TARGET_PATH" || exit 2
go run ./cmd/mify/ add client service1 --to service2 -p "$TARGET_PATH" || exit 2
go run ./cmd/mify/ remove client service1 --to service2 -p "$TARGET_PATH" || exit 2
go run ./cmd/mify/ add client service1 --to service2 -p "$TARGET_PATH" || exit 2
go run ./cmd/mify/ add frontend --template nuxtjs front -p "$TARGET_PATH" || exit 2
go run ./cmd/mify/ add client front --to service1 -p "$TARGET_PATH" || exit 2
go run ./cmd/mify/ add api-gateway -p "$TARGET_PATH" || exit 2

cd "$TARGET_PATH/go-services" || exit 2
go mod tidy || exit 2
go build ./cmd/service1 ./cmd/service2 || exit 2

cd "$TARGET_PATH/js-services/front" || exit 2
yarn install && yarn build || exit 2

git add * && git commit -m "init" || exit 2
