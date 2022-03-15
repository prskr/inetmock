cmake_minimum_required(VERSION 3.21)
project(bpf C)

set(CMAKE_C_STANDARD 11)

add_executable(bpf
        netflow.c
        types.h
        helpers.h)
