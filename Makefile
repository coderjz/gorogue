build:
	go build -o ./gorogue ./main/main.go

clean: 
	rm -f gorogue
	rm -rf dist

release:
	for GOOS in darwin linux windows; do \
		for GOARCH in 386 amd64; do \
			mkdir -p dist/$$GOOS/$$GOARCH; \
			if [[ $$GOOS == "windows" ]] ; then \
				GOARCH=$$GOARCH GOOS=$$GOOS go build -o ./dist/$$GOOS/$$GOARCH/gorogue.exe ./main/main.go; \
			else \
				GOARCH=$$GOARCH GOOS=$$GOOS go build -o ./dist/$$GOOS/$$GOARCH/gorogue ./main/main.go; \
			fi \
		done \
	done
