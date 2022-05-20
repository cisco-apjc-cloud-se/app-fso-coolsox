#!/usr/bin/env bash
export COPYFILE_DISABLE=true
helm lint ./coolsox
helm package coolsox #tar -czf coolsox-0.0.1.tgz coolsox
helm repo index ./
