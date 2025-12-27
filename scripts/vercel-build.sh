#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
UI_DIR="${ROOT_DIR}/ui"

if ! command -v go >/dev/null 2>&1; then
  GO_VERSION="${GO_VERSION:-$(awk '/^go / {print $2; exit}' "${ROOT_DIR}/go.mod" 2>/dev/null || true)}"
  if [[ -z "${GO_VERSION}" ]]; then
    GO_VERSION="1.22.6"
  elif [[ "${GO_VERSION}" =~ ^[0-9]+\.[0-9]+$ ]]; then
    GO_VERSION="${GO_VERSION}.0"
  fi

  GO_INSTALL_DIR="${HOME}/.cache/go-${GO_VERSION}"
  if [[ ! -x "${GO_INSTALL_DIR}/bin/go" ]]; then
    mkdir -p "${GO_INSTALL_DIR}"
    GO_ARCHIVE="go${GO_VERSION}.linux-amd64.tar.gz"
    GO_URL="https://go.dev/dl/${GO_ARCHIVE}"
    if command -v curl >/dev/null 2>&1; then
      curl -fsSL "${GO_URL}" | tar -xz -C "${GO_INSTALL_DIR}" --strip-components=1
    elif command -v wget >/dev/null 2>&1; then
      wget -qO- "${GO_URL}" | tar -xz -C "${GO_INSTALL_DIR}" --strip-components=1
    else
      echo "curl or wget is required to install Go." >&2
      exit 1
    fi
  fi

  export PATH="${GO_INSTALL_DIR}/bin:${PATH}"
fi

make -C "${ROOT_DIR}" ui-wasm
pnpm --dir "${UI_DIR}" run build
