FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y \
    vim \
    build-essential \
    git \
    mariadb-server \
    wget \
    munge \
    slurm-client \
    curl \
    dirmngr \
    apt-transport-https \
    lsb-release \
    ca-certificates \
    sudo \
    python3.9 \
    python3-pip \
    && useradd -m admin -s /usr/bin/bash -d /home/admin \
    && echo "admin:admin" | chpasswd \
    && adduser admin sudo \
    && echo "admin ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers \
    && curl -sL https://deb.nodesource.com/setup_12.x | bash - \
    && apt-get install -y nodejs npm \
    && apt-get install -y libopenmpi-dev \
    && npm install -g configurable-http-proxy \
    && pip3 install jupyterlab==2.1.2 \
    mpi4py \
    jupyterlab_slurm

COPY slurm.conf /etc/slurm-llnl/
COPY cgroup.conf /etc/slurm-llnl/
COPY docker-entrypoint.sh /etc/slurm-llnl/
RUN mkdir -p /home/admin \
&& chown admin /home/admin \
&& chmod 755 /etc/slurm-llnl/docker-entrypoint.sh

EXPOSE 8888

USER admin
WORKDIR /home/admin

ENTRYPOINT ["/etc/slurm-llnl/docker-entrypoint.sh"]