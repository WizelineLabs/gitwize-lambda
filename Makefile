.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/update_one_repo functions/update_one_repo/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/load_full_one_repo functions/load_full_one_repo/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/update_all_repos functions/update_all_repos/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/load_metrics functions/load_metrics/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/git_native_example functions/git_native_example/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
