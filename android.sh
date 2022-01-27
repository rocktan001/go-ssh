#!/bin/bash
# set -x
# python $NDK_HOME/build/tools/make_standalone_toolchain.py  --arch=arm64 --install-dir="D:/tanzhongqiang/go_win_workspace/toolchain-arm64-ndk-r21"

export GOARCH=arm64

export GOOS=android
export CGO_ENABLED=1
export CXX="D:/tanzhongqiang/go_win_workspace/toolchain-arm64-ndk-r21\bin\clang++"
export CC="D:/tanzhongqiang/go_win_workspace/toolchain-arm64-ndk-r21\bin\clang"
