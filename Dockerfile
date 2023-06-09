FROM gradle:jdk17 AS base
ENV PATH=$PATH:/usr/local/go/bin
RUN apt-get update \
    && apt-get install -y python3-pip python-is-python3 graphviz graphviz-dev
RUN mkdir "upload" \
    && wget "https://github.com/NationalSecurityAgency/ghidra/releases/download/Ghidra_10.2.3_build/ghidra_10.2.3_PUBLIC_20230208.zip" -O ghidra.zip \
    && unzip ghidra.zip \
    && wget "https://github.com/mandiant/Ghidrathon/archive/refs/tags/v2.0.1.zip" \
    && unzip v2.0.1.zip \
    && cd "Ghidrathon-2.0.1/" \
    && gradle --no-daemon -PGHIDRA_INSTALL_DIR="/home/gradle/ghidra_10.2.3_PUBLIC" \
    && cd "/home/gradle/ghidra_10.2.3_PUBLIC/Ghidra/Extensions" \
    && unzip "/home/gradle/Ghidrathon-2.0.1/dist/ghidra_10.2.3_PUBLIC_$(date '+%Y%m%d')_Ghidrathon-2.0.1"
COPY ./requirements.txt ./src/
RUN pip3 install -r ./src/requirements.txt

FROM golang:1.20.2-alpine3.17 AS builder
WORKDIR ./src
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY . /go/src
RUN go build ./main.go

FROM base
COPY --from=builder /go/src ./src
WORKDIR ./src
ENTRYPOINT ["./main"]
EXPOSE 8000