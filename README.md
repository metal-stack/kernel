# Kernel

This project is used to build a small linux kernels with the `kexec` feature enabled.

In [Linuxkit](https://github.com/linuxkit/linuxkit) this feature is disabled for `x86_64`.

The linuxkit kernel config is used as base for our compilation: s. https://github.com/linuxkit/linuxkit/tree/master/kernel

## Build

Execute `make`, optionally define the variables `KERNEL_SERIES`, `KERNEL_VERSION` and `KERNEL_MAJOR` in order to diverge from the defaults.

## Manage new major kernel config upgrades.

- download the new kernel sources (e.g. 4.20)
- copy the existing config to the sources directory *.config*
- run `make oldconfig` and you can answer all questions with the given default
- the create *.config* will then contain all new options, and all old options.
