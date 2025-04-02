#!/bin/bash
if command -v node >/dev/null 2>&1; then
    echo "Node.js is installed."
else
    echo "Node.js is not installed."
    echo "Install nodejs_2.0 ..."
    curl -fsSL https://rpm.nodesource.com/setup_20.x | sudo bash -
    dnf install -y nsolid
fi
