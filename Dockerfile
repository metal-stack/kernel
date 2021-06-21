FROM alpine:3.14

ARG KERNEL_MAJOR
ARG KERNEL_VERSION
ARG KERNEL_SERIES

RUN set -ex \
 && apk add \
    argp-standalone \
    automake \
    bash \
    bc \
    binutils-dev \
    bison \
    build-base \
    curl \
    diffutils \
    findutils \
    flex \
    git \
    gmp-dev \
    gnupg \
    installkernel \
    kmod \
    elfutils-dev \
    linux-headers \
    libunwind-dev \
    lz4 \
    mpc1-dev \
    mpfr-dev \
    ncurses-dev \
    patch \
    sed \
    squashfs-tools \
    tar \
    xz \
    xz-dev \
    zlib-dev \
    zstd

ENV KERNEL_SOURCE=https://www.kernel.org/pub/linux/kernel/${KERNEL_MAJOR}/linux-${KERNEL_VERSION}.tar.xz
ENV KERNEL_SHA256_SUMS=https://www.kernel.org/pub/linux/kernel/${KERNEL_MAJOR}/sha256sums.asc
ENV KERNEL_PGP2_SIGN=https://www.kernel.org/pub/linux/kernel/${KERNEL_MAJOR}/linux-${KERNEL_VERSION}.tar.sign

# tell xz decompressor to use as much threads as cpu cores
ENV XZ_OPT="--threads=0"

# We copy the entire directory. This copies some unneeded files, but
# allows us to check for the existence /patches-${KERNEL_SERIES} to
# build kernels without patches.
COPY / /

# Download and verify kernel
# PGP keys: 589DA6B1 (greg@kroah.com) & 6092693E (autosigner@kernel.org) & 00411886 (torvalds@linux-foundation.org)
RUN set -ex \
 && curl -fsSLO ${KERNEL_SHA256_SUMS} \
 && gpg2 -q --import keys.asc \
 && gpg2 --verify sha256sums.asc \
 && KERNEL_SHA256=$(grep linux-${KERNEL_VERSION}.tar.xz sha256sums.asc | cut -d ' ' -f 1) \
 && [ -f linux-${KERNEL_VERSION}.tar.xz ] || curl -fsSLO ${KERNEL_SOURCE} \
 && echo "${KERNEL_SHA256}  linux-${KERNEL_VERSION}.tar.xz" | sha256sum -c - \
 && xz -d linux-${KERNEL_VERSION}.tar.xz \
 && curl -fsSLO ${KERNEL_PGP2_SIGN} \
 && gpg2 --verify linux-${KERNEL_VERSION}.tar.sign linux-${KERNEL_VERSION}.tar \
 && tar --absolute-names -xf linux-${KERNEL_VERSION}.tar && mv /linux-${KERNEL_VERSION} /linux

WORKDIR /linux

# Kernel config
RUN set -ex \
 && KERNEL_DEF_CONF=/linux/arch/x86/configs/x86_64_defconfig \
 && cp /config-${KERNEL_SERIES}-$(uname -m) ${KERNEL_DEF_CONF} \
 && make clean \
 && make oldconfig \
 && make scripts \
 && make defconfig \
 && make oldconfig

# Patch kernel
RUN set -ex \
 && case "$KERNEL_VERSION" in 5.0*) \
      patch -p1 < /0001-ipconfig-add-carrier_timeout-kernel-parameter.patch; \
    esac

# Kernel
RUN set -ex \
 && make -j "$(getconf _NPROCESSORS_ONLN)" KCFLAGS="-fno-pie" \
 && mv arch/x86/boot/bzImage /metal-kernel
