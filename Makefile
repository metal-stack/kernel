KERNEL_SERIES := $(or $(KERNEL_SERIES),mainline)
KERNEL_VERSION := $(or $(KERNEL_VERSION),5.4.60)
KERNEL_MAJOR := $(or $(KERNEL_MAJOR),v5.x)

.PHONY: build
build:
	rm -f kernel.tar metal-kernel
	docker build --build-arg KERNEL_MAJOR=$(KERNEL_MAJOR) \
							 --build-arg KERNEL_VERSION=$(KERNEL_VERSION) \
							 --build-arg KERNEL_SERIES=$(KERNEL_SERIES) \
							 --tag metal-stack/kernel .
	docker export $(shell docker create metal-stack/kernel /dev/null) > kernel.tar
