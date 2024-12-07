#!/bin/bash

export AWS_EXECUTION_ENV=AWS_Lambda_%s

if [ -z "$AWS_LAMBDA_EXEC_WRAPPER" ]; then
  # NOTE: Run the debugpy module
  exec /var/lang/bin/%s -Xfrozen_modules=off -m debugpy --listen 0.0.0.0:5678 /var/runtime/bootstrap.py
else
  wrapper="$AWS_LAMBDA_EXEC_WRAPPER"
  if [ ! -f "$wrapper" ]; then
    echo "$wrapper: does not exist"
    exit 127
  fi
  if [ ! -x "$wrapper" ]; then
    echo "$wrapper: is not an executable"
    exit 126
  fi
    # NOTE: Run the debugpy module
    exec -- "$wrapper" /var/lang/bin/%s -Xfrozen_modules=off -m debugpy --listen 0.0.0.0:5678 /var/runtime/bootstrap.py
fi