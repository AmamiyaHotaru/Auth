#!/bin/bash
set -o allexport
source .env
set +o allexport

wails build -clean -ldflags="-X auth/utils.embeddedMasterPassword=$MASTER_PASSWORD"
