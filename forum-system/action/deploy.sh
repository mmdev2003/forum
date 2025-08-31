#!/bin/bash

apt update -y && apt upgrade -y

cd ..
git clone git@github.com:forum-mmdev/forum-authentication.git
git clone git@github.com:forum-mmdev/forum-authorization.git
git clone git@github.com:forum-mmdev/forum-dialog.git
git clone git@github.com:forum-mmdev/forum-thread.git
git clone git@github.com:forum-mmdev/forum-user.git
git clone git@github.com:forum-mmdev/forum-frame.git
git clone git@github.com:forum-mmdev/forum-status.git
git clone git@github.com:forum-mmdev/forum-notification.git
git clone git@github.com:forum-mmdev/forum-support.git
git clone git@github.com:forum-mmdev/forum-payment.git
git clone git@github.com:forum-mmdev/forum-admin.git
cd forum-system

./tools/docker/install.sh
./tools/nginx/install.sh
