doc:
	swag fmt && swag init

unittest:
	go test -timeout 9000s -a -v -coverprofile=coverage.out -coverpkg=./... ./... 2>&1 | tee report.out

mock:
	mockgen -source=./app/repositories/user.go -destination=./mocks/IUserRepository.go  --build_flags=--mod=vendor -package=mocks
	mockgen -source=./app/repositories/order.go -destination=./mocks/IOrderRepository.go  --build_flags=--mod=vendor -package=mocks
	mockgen -source=./app/repositories/product.go -destination=./mocks/IProductRepository.go  --build_flags=--mod=vendor -package=mocks
