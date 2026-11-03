#!/bin/sh
# cgo 交叉编译用的 CC 包装器。
#
# 为什么需要包装器：Go 1.17 的 cmd/cgo 用 strings.Fields 拆分 $CC（gcc.go 的 gccBaseCmd），
# 对含空格的 CC（如 "zig cc -target x86_64-linux-musl"）在 -E 预处理那一步会丢失 "cc -target ..."，
# 导致 zig 把 -E 当子命令报 "unknown command: -E"。让 CC 指向本脚本（单个无空格 token）即可规避。
#
# 目标三元组经环境变量 ZIG_TARGET 透传（如 x86_64-linux-musl / x86_64-windows-gnu）。
# 要求 zig 已安装且在 PATH 上。
exec zig cc -target "${ZIG_TARGET:?ZIG_TARGET 未设置}" "$@"
