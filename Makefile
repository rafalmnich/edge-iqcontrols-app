version="0.0.2"
version_file=package/debian/opt/thingsplex/$(app_name)/VERSION
working_dir=$(shell pwd)
app_name=iqcontrols
remote_host = "fhtunnel@52.58.200.103"
remote_port = "8000"

clean:
	-rm ./src/$(app_name)
	-rm ./package/debian/opt/thingsplex/$(app_name)/$(app_name)
	find package/debian -name ".DS_Store" -delete

build-go:
	cd ./src; go build -o $(app_name) main.go; cd ../

build-go-arm:
	cd ./src; GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o $(app_name) main.go; cd ../

build-go-amd:
	cd ./src; GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(app_name) main.go; cd ../

configure-arm:
	sed -i.bak "1s/.*/$(version)/" $(version_file)
	rm $(version_file).bak
	sed -i.bak "s/Version: .*/Version: $(version)/" package/debian/DEBIAN/control
	sed -i.bak "s/Architecture: .*/Architecture: armhf/" package/debian/DEBIAN/control
	rm package/debian/DEBIAN/control.bak

configure-amd64:
	sed -i.bak "1s/.*/$(version)/" $(version_file)
	rm $(version_file).bak
	sed -i.bak "s/Version: .*/Version: $(version)/" package/debian/DEBIAN/control
	sed -i.bak "s/Architecture: .*/Architecture: amd64/" package/debian/DEBIAN/control
	rm package/debian/DEBIAN/control.bak

package-deb-lint:
	docker run -w /root -v $(working_dir)/package/build/:/root/ -it eddelbuettel/lintian lintian $(app_name)_$(version)_armhf.deb --no-tag-display-limit

package-deb:
	@echo "Packaging application using Thingsplex debian package layout..."
	mkdir -p package/debian/var/log/thingsplex/$(app_name)
	mkdir -p package/debian/opt/thingsplex/$(app_name)/data
	chmod -R 755 package/debian
	chmod 644 package/debian/opt/thingsplex/$(app_name)/defaults/*
	chmod 644 package/debian/opt/thingsplex/$(app_name)/VERSION
	chmod 644 package/debian/usr/lib/systemd/system/$(app_name).service
	cp ./src/$(app_name) package/debian/opt/thingsplex/$(app_name)
	docker run --rm -v ${working_dir}:/build -w /build --name debuild debian dpkg-deb --build package/debian
	@echo "Done"

deb-arm: clean configure-arm build-go-arm package-deb
	mv package/debian.deb package/build/$(app_name)_$(version)_armhf.deb

deb-amd: clean configure-amd64 build-go-amd package-deb
	mv package/debian.deb package/build/$(app_name)_$(version)_amd64.deb

upload:
	scp -P ${remote_port} package/build/$(app_name)_$(version)_armhf.deb $(remote_host):~/

run: build-go
	cd ./testdata; ../src/$(app_name) -c ./; cd ../

lint:
	cd src; golangci-lint run; cd ..

test:
	cd src; go test; cd ..

.phony : clean
