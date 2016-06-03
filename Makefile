Dockerfile.stamp: Dockerfile
	docker build -t courier-packager .
	@touch Dockerfile.stamp

bin/courier.rpm: Dockerfile.stamp
	mkdir -p dist/usr/local/bin && cp bin/courier dist/usr/local/bin/courier
	chmod 0755 dist/usr/local/bin/courier
	docker run --rm -u `whoami` -w `pwd` -v /etc/passwd:/etc/passwd -v /etc/group:/etc/group -v `pwd`:`pwd` courier-packager "fpm -s dir -t rpm -v 0.0.1 -n courier -a x86_64 -C dist"


