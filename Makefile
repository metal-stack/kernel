GO111MODULE := on

.PHONY: all
all:
	docker build \
		--build-arg KERNEL_MAJOR="v5.x" \
		--build-arg KERNEL_VERSION=${{ matrix.node }} \
    	--build-arg KERNEL_SERIES="mainline" \
		--tag metal-stack/kernel \
		.
	docker export $(docker create metal-stack/kernel /dev/null) > kernel.tar
	tar -xf kernel.tar metal-kernel
	mv metal-kernel metal-kernel-${KERNEL_VERSION}
	