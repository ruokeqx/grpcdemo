cmake_minimum_required(VERSION 3.15)

project(grpc_server_cpp CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Threads REQUIRED)
find_package(Protobuf REQUIRED)
find_package(gRPC CONFIG REQUIRED)


set(PROTO_PATH "${CMAKE_CURRENT_SOURCE_DIR}/../proto")
set(PROTO_FILE "${PROTO_PATH}/tray/status.proto")

set(PROTO_BINARY_DIR "${CMAKE_CURRENT_BINARY_DIR}/generated")
file(MAKE_DIRECTORY ${PROTO_BINARY_DIR})

set(PROTO_SRCS "${PROTO_BINARY_DIR}/tray/status.pb.cc")
set(PROTO_HDRS "${PROTO_BINARY_DIR}/tray/status.pb.h")
set(GRPC_SRCS "${PROTO_BINARY_DIR}/tray/status.grpc.pb.cc")
set(GRPC_HDRS "${PROTO_BINARY_DIR}/tray/status.grpc.pb.h")

find_program(_PROTOBUF_PROTOC protoc)
find_program(_GRPC_CPP_PLUGIN_EXECUTABLE grpc_cpp_plugin)

add_custom_command(
    OUTPUT "${PROTO_SRCS}" "${PROTO_HDRS}" "${GRPC_SRCS}" "${GRPC_HDRS}"
    COMMAND ${_PROTOBUF_PROTOC}
    ARGS --grpc_out "${PROTO_BINARY_DIR}"
         --cpp_out "${PROTO_BINARY_DIR}"
         -I "${PROTO_PATH}"
         --plugin=protoc-gen-grpc="${_GRPC_CPP_PLUGIN_EXECUTABLE}"
         "${PROTO_FILE}"
    DEPENDS "${PROTO_FILE}"
    COMMENT "Generating protobuf and gRPC files"
)


add_executable(
        server-cpp
        server.cpp
        ${PROTO_SRCS}
        ${GRPC_SRCS}
)

target_link_libraries(
        server-cpp
        PRIVATE
        gRPC::grpc++
        protobuf::libprotobuf
)

target_include_directories(
        server-cpp
        PRIVATE
        ${PROTO_BINARY_DIR}
)