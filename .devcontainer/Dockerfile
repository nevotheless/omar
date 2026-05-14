FROM mcr.microsoft.com/devcontainers/go:1.24

RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    python3-argcomplete \
    python3-venv \
    btrfs-progs \
    dosfstools \
    e2fsprogs \
    cryptsetup \
    systemd \
    systemd-boot \
    kmod \
    xz-utils \
    zstd \
    gnupg \
    openssl \
    qemu-system-x86 \
    qemu-utils \
    skopeo \
    && rm -rf /var/lib/apt/lists/*

RUN pip3 install mkosi

RUN mkosi --version

RUN echo 'alias make="make -C /workspaces/omar"' >> /home/vscode/.bashrc
RUN echo 'alias build="cd /workspaces/omar && go build -ldflags=\"-X github.com/nevotheless/omar/internal/version.Version=2026.5.0-dev\" -o bin/omar ./cmd/omar"' >> /home/vscode/.bashrc

USER vscode
RUN go install github.com/spf13/cobra-cli@latest

WORKDIR /workspaces/omar
