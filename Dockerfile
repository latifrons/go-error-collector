FROM apache/skywalking-go:0.2.0-go1.20 as builder

RUN apt-get update && apt-get install -y unzip
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v24.4/protoc-24.4-linux-x86_64.zip && \
    unzip -o protoc-24.4-linux-x86_64.zip -d /usr/local bin/protoc && \
    unzip -o protoc-24.4-linux-x86_64.zip -d /usr/local include/* && \
    rm -rf protoc-24.4-linux-x86_64.zip

RUN curl -OL https://github.com/bufbuild/protoc-gen-validate/releases/download/v1.0.2/protoc-gen-validate_1.0.2_linux_amd64.tar.gz && \
    tar -xvf protoc-gen-validate_1.0.2_linux_amd64.tar.gz && \
    mv protoc-gen-validate /usr/local/bin/ && \
    rm -rf protoc-gen-validate_1.0.2_linux_amd64.tar.gz

RUN apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /src

ADD ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN make proto

# RUN skywalking-go-agent -inject . -all
# RUN go build -toolexec="skywalking-go-agent" -a -o main .
RUN go build -a -o main .

# Copy OG into basic alpine image
FROM debian:stable-slim

RUN apt-get update && apt-get install -y curl iotop tzdata

WORKDIR /www

COPY --from=builder src/data/config ./data/config/
#COPY --from=builder src/rpc/docs ./rpc/docs/
COPY --from=builder src/main .

ENV SW_AGENT_REPORTER_GRPC_BACKEND_SERVICE=skywalking-grpc.all.internal.traefik:80
# ENV SW_AGENT_NAME="${NOMAD_JOB_NAME}"

ENTRYPOINT ["./main"]
