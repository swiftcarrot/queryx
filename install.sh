#!/bin/sh

set -e

# Some utilities from https://github.com/client9/shlib

echoerr() {
  printf "$@\n" 1>&2
}

log_info() {
  printf "\033[38;5;61m  ==>\033[0;00m $@\n"
}

log_crit() {
  echoerr
  echoerr "  \033[38;5;125m$@\033[0;00m"
  echoerr
}

is_command() {
  command -v "$1" >/dev/null
  #type "$1" > /dev/null 2> /dev/null
}

http_download_curl() {
  local_file=$1
  source_url=$2
  header=$3
  if [ -z "$header" ]; then
    code=$(curl -w '%{http_code}' -sL -o "$local_file" "$source_url")
  else
    code=$(curl -w '%{http_code}' -sL -H "$header" -o "$local_file" "$source_url")
  fi
  if [ "$code" != "200" ]; then
    log_crit "Error downloading, got $code response from server"
    return 1
  fi
  return 0
}

http_download_wget() {
  local_file=$1
  source_url=$2
  header=$3
  if [ -z "$header" ]; then
    wget -q -O "$local_file" "$source_url"
  else
    wget -q --header "$header" -O "$local_file" "$source_url"
  fi
}

http_download() {
  if is_command curl; then
    http_download_curl "$@"
    return
  elif is_command wget; then
    http_download_wget "$@"
    return
  fi
  log_crit "http_download unable to find wget or curl"
  return 1
}

http_copy() {
  tmp=$(mktemp)
  http_download "${tmp}" "$1" "$2" || return 1
  body=$(cat "$tmp")
  rm -f "${tmp}"
  echo "$body"
}

uname_os() {
  os=$(uname -s | tr '[:upper:]' '[:lower:]')

  # fixed up for https://github.com/client9/shlib/issues/3
  case "$os" in
    msys_nt*) os="windows" ;;
    mingw*) os="windows" ;;
  esac

  # other fixups here
  echo "$os"
}

uname_os_check() {
  os=$(uname_os)
  case "$os" in
    darwin) return 0 ;;
    dragonfly) return 0 ;;
    freebsd) return 0 ;;
    linux) return 0 ;;
    android) return 0 ;;
    nacl) return 0 ;;
    netbsd) return 0 ;;
    openbsd) return 0 ;;
    plan9) return 0 ;;
    solaris) return 0 ;;
    windows) return 0 ;;
  esac
  log_crit "uname_os_check '$(uname -s)' got converted to '$os' which is not a GOOS value. Please file bug at https://github.com/client9/shlib"
  return 1
}

uname_arch() {
  arch=$(uname -m)
  case $arch in
    x86_64) arch="amd64" ;;
    x86) arch="386" ;;
    i686) arch="386" ;;
    i386) arch="386" ;;
    aarch64) arch="arm64" ;;
    armv5*) arch="armv5" ;;
    armv6*) arch="armv6" ;;
    armv7*) arch="armv7" ;;
  esac
  echo ${arch}
}

uname_arch_check() {
  arch=$(uname_arch)
  case "$arch" in
    386) return 0 ;;
    amd64) return 0 ;;
    arm64) return 0 ;;
    armv5) return 0 ;;
    armv6) return 0 ;;
    armv7) return 0 ;;
    ppc64) return 0 ;;
    ppc64le) return 0 ;;
    mips) return 0 ;;
    mipsle) return 0 ;;
    mips64) return 0 ;;
    mips64le) return 0 ;;
    s390x) return 0 ;;
    amd64p32) return 0 ;;
  esac
  log_crit "uname_arch_check '$(uname -m)' got converted to '$arch' which is not a GOARCH value.  Please file bug report at https://github.com/client9/shlib"
  return 1
}

mktmpdir() {
  test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"
  mkdir -p "${TMPDIR}"
  echo "${TMPDIR}"
}

start() {
  uname_os_check
  uname_arch_check

  pkg="github.com/swiftcarrot/queryx"
  bin="queryx.tar.gz"
  version="v0.2.11"

  prefix=${PREFIX:-"/usr/local/bin"}
  tmp="$(mktmpdir)"

  echo
  log_info "Downloading $pkg@$version"
  log_info "Downloading binary for $os $arch"

  asset=$(curl -s https://api.github.com/repos/swiftcarrot/queryx/releases/tags/$version | grep -i "browser_download_url.*$os\_$arch" | sed -E 's/.*"([^"]+)".*/\1/')
  http_download "$tmp/queryx_$version.tar.gz" "$asset"

  cd "$tmp"
  tar zxf "queryx_$version.tar.gz"

  if [ -w "$prefix" ]; then
  log_info "Installing queryx to $prefix"
    install queryx "$prefix"
  else
    log_info "Permissions required for installation to $prefix — alternatively specify a new directory with:"
    sudo install queryx "$prefix"
  fi

  log_info "Installation complete"
  echo
}

start
