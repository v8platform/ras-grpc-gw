# Substitute BIN for your bin directory.
# Substitute VERSION for the current released version.
# Substitute BINARY_NAME for "buf", "protoc-gen-buf-breaking", or "protoc-gen-buf-lint".
BIN="/usr/local/bin" && \
VERSION="1.0.0-rc6" && \
BINARY_NAME="buf" && \
  sudo curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v${VERSION}/${BINARY_NAME}-$(uname -s)-$(uname -m)" \
    -o "${BIN}/${BINARY_NAME}" && \
  sudo chmod +x "${BIN}/${BINARY_NAME}"