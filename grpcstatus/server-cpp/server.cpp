#include <iostream>
#include <memory>
#include <string>
#include <thread>
#include <chrono>
#include <atomic>

#include <grpcpp/grpcpp.h>

#include "tray/status.grpc.pb.h"
#include "tray/status.pb.h"
#include "google/protobuf/empty.pb.h"

using namespace tray;

class StatusServiceImpl final : public StatusService::Service {
public:
    grpc::Status StreamStatus(
        grpc::ServerContext* context,
        grpc::ServerReaderWriter<StatusStreamMessage, StatusStreamMessage>* stream) override {

        std::cout << "✅ Client connected" << std::endl;

        std::atomic<bool> is_client_disconnected(false);

        std::thread reader_thread([stream, &is_client_disconnected]() {
            StatusStreamMessage msg;
            while (stream->Read(&msg)) {
                if (msg.has_status()) {
                    std::cout << "📬 Received message from client. Content: "  << msg.status().DebugString() << std::endl;
                } else if (msg.has_pull_request()) {
                    std::cout << "➡️ should never receive pull request from client !!! " << std::endl;
                } else {
                    std::cout << "➡️ unkown " << std::endl;
                }
            }

            is_client_disconnected = true;
        });

        while (!is_client_disconnected) {
            std::this_thread::sleep_for(std::chrono::seconds(5));

            if (is_client_disconnected) {
                break;
            }

            std::cout << "➡️ Sending pull request to client..." << std::endl;

            StatusStreamMessage pull;
            pull.mutable_pull_request();

            if (!stream->Write(pull)) {
                std::cerr << "❌ Write failed. Client might have disconnected." << std::endl;
                is_client_disconnected = true;
            }
        }

        reader_thread.join();

        std::cout << "🔌 Client disconnected" << std::endl;
        return grpc::Status::OK;
    }
};

void RunServer() {
    std::string server_address("0.0.0.0:50051");
    StatusServiceImpl service;

    grpc::ServerBuilder builder;
    builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
    builder.RegisterService(&service);

    std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
    std::cout << "🚀 Server started on " << server_address << std::endl;

    server->Wait();
}

int main(int argc, char** argv) {
    RunServer();
    return 0;
}