FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm

RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/axboe/liburing.git /usr/src/liburing

RUN cd /usr/src/liburing && \
    ./configure --cc=gcc --cxx=g++ && \
    make -j$(nproc) && \
    make install && \
    rm -rf /usr/src/liburing
