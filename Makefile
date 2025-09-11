.ONESHELL:
KERNEL_VERSION := $(or ${KERNEL_VERSION},6.12.46)

.PHONY: all
all: build
	
.PHONY: build
build:
	docker build --build-arg KERNEL_MAJOR="v6.x" \
                 --build-arg KERNEL_VERSION=${KERNEL_VERSION} \
                 --build-arg KERNEL_SERIES="mainline" \
                 --tag metal-stack/kernel .

.PHONY: save
save:
	docker export $(shell docker create metal-stack/kernel /dev/null) > kernel.tar
	tar -xf kernel.tar metal-kernel
	mv metal-kernel metal-kernel-${KERNEL_VERSION}
