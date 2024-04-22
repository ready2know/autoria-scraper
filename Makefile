DOWNLOADER_NAME=downloader
SCRAPPER_NAME=scrapper
PGIMPORTER_NAME=pgImporter

BUILD_FOLDER=build

build_downloader:
	GOARCH=amd64 GOOS=linux go build -o ${BUILD_FOLDER}/${DOWNLOADER_NAME}-linux cmd/${DOWNLOADER_NAME}/main.go
	GOARCH=arm64 GOOS=darwin go build -o ${BUILD_FOLDER}/${DOWNLOADER_NAME}-darwin cmd/${DOWNLOADER_NAME}/main.go

build_scrapper:
	GOARCH=amd64 GOOS=linux go build -o ${BUILD_FOLDER}/${SCRAPPER_NAME}-linux cmd/${SCRAPPER_NAME}/main.go
	GOARCH=arm64 GOOS=darwin go build -o ${BUILD_FOLDER}/${SCRAPPER_NAME}-darwin cmd/${SCRAPPER_NAME}/main.go

build_pgImporter:
	GOARCH=amd64 GOOS=linux go build -o ${BUILD_FOLDER}/${PGIMPORTER_NAME}-linux cmd/${PGIMPORTER_NAME}/main.go
	GOARCH=arm64 GOOS=darwin go build -o ${BUILD_FOLDER}/${PGIMPORTER_NAME}-darwin cmd/${PGIMPORTER_NAME}/main.go

run_downloader:
	go run cmd/${DOWNLOADER_NAME}/main.go

run_scrapper:
	go run cmd/${SCRAPPER_NAME}/main.go

run_pgImporter:
	go run cmd/${PGIMPORTER_NAME}/main.go

clean:
	 go clean
	 rm ${BUILD_FOLDER}/*