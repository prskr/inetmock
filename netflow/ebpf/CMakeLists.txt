cmake_minimum_required(VERSION 3.21)
project(ebpf C)

set(CMAKE_C_STANDARD 11)

include_directories(.)

add_executable(nat
        common.h
        helpers.h
        types.h
        nat.c)

add_executable(firewall
        common.h
        helpers.h
        types.h
        firewall.c)

add_executable(tests
        common.h
        helpers.h
        tests.c
)