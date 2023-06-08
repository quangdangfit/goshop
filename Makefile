doc:
	swag fmt && swag init

unittest:
	go test -timeout 9000s -a -v -coverpkg=./... ./test

mock:
	#mockgen -source=./internal/interfaces/IArenaHeroRepository.go -destination=./mocks/IArenaHeroRepository.go  --build_flags=--mod=vendor -package=mocks
