#!/bin/sh -e

appName="helm-clear"

# shellcheck disable=SC2002
version="$(cat plugin.yaml | grep "version" | cut -d '"' -f 2)"
echo "Downloading and installing helm-clear v${version} ..."

url="https://github.com/mocyuto/helm-clear/releases/download/v${version}/${appName}"
if [ "$(uname)" = "Darwin" ]; then
  if [ "$(uname -m)" = "arm64" ]; then
    url=$url"_Darwin_arm64.tar.gz"
  else
    url=$url"_Darwin_x86_64.tar.gz"
  fi
elif [ "$(uname)" = "Linux" ] ; then
    if [ "$(uname -m)" = "aarch64" ] || [ "$(uname -m)" = "arm64" ]; then
        url=$url"_Linux_arm64.tar.gz"
    else
        url=$url"_Linux_x86_64.tar.gz"
    fi
else
    url=$url"_Windows_x86_64.zip"
fi

echo "$url"

mkdir -p "bin"
mkdir -p "config"
mkdir -p "releases/v${version}"

# Download with curl if possible.
# shellcheck disable=SC2230
if [ -x "$(which curl 2>/dev/null)" ]; then
    curl -sSL "${url}" -o "releases/v${version}.tar.gz"
else
    wget -q "${url}" -O "releases/v${version}.tar.gz"
fi
tar xzf "releases/v${version}.tar.gz" -C "releases/v${version}"
mv "releases/v${version}/clear" "bin/clear" || \
    mv "releases/v${version}/clear.exe" "bin/clear"

