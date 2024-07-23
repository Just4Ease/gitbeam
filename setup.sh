#!/usr/bin/env bash
git submodule update --init --recursive --remote --checkout --force --rebase --recursive || exit;

rm -rf ../gitbeam.repo.manager || exit;
git clone https://github.com/Just4Ease/gitbeam.repo.manager ../gitbeam.repo.manager
cd ../gitbeam.repo.manager && git submodule update --init --recursive --remote --checkout --force --rebase --recursive || exit;
cd - || exit;


rm -rf ../gitbeam.commit.monitor || exit;
git clone https://github.com/Just4Ease/gitbeam.commit.monitor ../gitbeam.commit.monitor
cd ../gitbeam.commit.monitor && git submodule update --init --recursive --remote --checkout --force --rebase --recursive || exit;
cd - || exit;

docker-compose up
