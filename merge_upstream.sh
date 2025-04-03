#!/bin/bash
# display remote : git remote -v
# git remote add upstream https://github.com/cockpit-project/starter-kit.git
git fetch upstream
git checkout main
git merge upstream/main
git push origin main