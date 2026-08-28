#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
artifact_dir="${C311_ARTIFACT_DIR:-$(mktemp -d "${TMPDIR:-/tmp}/c311-packages.XXXXXX")}"

mkdir -p "$artifact_dir"

build_and_pack () {
  local package_dir="$1"
  local artifact_name="$2"

  corepack yarn --cwd "$repo_root/$package_dir" install --frozen-lockfile
  corepack yarn --cwd "$repo_root/$package_dir" build
  corepack yarn --cwd "$repo_root/$package_dir" pack --filename "$artifact_dir/$artifact_name.tgz"
}

install_artifact () {
  local app_dir="$1"
  local package_name="$2"
  local artifact="$3"
  local target="$repo_root/$app_dir/node_modules/@cortezaproject/$package_name"

  case "$target" in
    "$repo_root"/*/node_modules/@cortezaproject/*) ;;
    *) echo "Refusing to replace an unscoped package path: $target" >&2; exit 1 ;;
  esac

  rm -rf "$target"
  mkdir -p "$target"
  tar -xzf "$artifact" --strip-components=1 -C "$target"
}

build_and_pack lib/js corteza-js

corepack yarn --cwd "$repo_root/lib/vue" install --frozen-lockfile
install_artifact lib/vue corteza-js "$artifact_dir/corteza-js.tgz"
corepack yarn --cwd "$repo_root/lib/vue" build
corepack yarn --cwd "$repo_root/lib/vue" pack --filename "$artifact_dir/corteza-vue.tgz"

for app in client/web/compose client/web/admin; do
  corepack yarn --cwd "$repo_root/$app" install --frozen-lockfile
  install_artifact "$app" corteza-js "$artifact_dir/corteza-js.tgz"
  install_artifact "$app" corteza-vue "$artifact_dir/corteza-vue.tgz"

  APP_ROOT="$repo_root/$app" node <<'NODE'
const { createRequire } = require('node:module')
const { readFileSync } = require('node:fs')

const appRoot = process.env.APP_ROOT
const appRequire = createRequire(`${appRoot}/package.json`)
const jsEntry = appRequire.resolve('@cortezaproject/corteza-js')
const vueEntry = appRequire.resolve('@cortezaproject/corteza-vue')
if (!jsEntry.startsWith(`${appRoot}/node_modules/@cortezaproject/corteza-js/`)) {
  throw new Error(`C311 JS package resolved outside the app package: ${jsEntry}`)
}
if (!vueEntry.startsWith(`${appRoot}/node_modules/@cortezaproject/corteza-vue/`)) {
  throw new Error(`C311 Vue package resolved outside the app package: ${vueEntry}`)
}
const js = appRequire('@cortezaproject/corteza-js')
const vueBundle = readFileSync(vueEntry, 'utf8')
if (typeof js.MockC311Provider !== 'function' || !vueBundle.includes('C311AppShell')) {
  throw new Error('Generated C311 exports are missing from the installed local packages')
}
console.log(`C311 packages verified for ${appRoot}: ${jsEntry}; ${vueEntry}`)
NODE
done

echo "C311 package artifacts: $artifact_dir"
