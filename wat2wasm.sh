#!/bin/bash
wasm-tools parse wat/fixtures/empty.wat -o wasm/fixtures/empty.wasm
wasm-tools parse wat/fixtures/func.wat -o wasm/fixtures/func.wasm
wasm-tools parse wat/fixtures/add.wat -o wasm/fixtures/add.wasm 