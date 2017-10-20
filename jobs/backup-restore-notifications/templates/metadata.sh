#!/usr/bin/env bash
echo "---
backup_should_be_locked_before:
- job_name: cloud_controller_ng
  release: capi
- job_name: uaa
  release: uaa
restore_should_be_locked_before:
- job_name: cloud_controller_ng
  release: capi
- job_name: uaa
  release: uaa"
