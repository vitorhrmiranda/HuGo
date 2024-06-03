#!/bin/sh

git checkout main && \
git push origin main && \
git checkout pages && \
git reset --hard origin/pages && \
git merge --no-ff --no-edit origin/main && \
git push origin pages
