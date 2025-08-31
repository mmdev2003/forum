#!/bin/bash
cd ..
all_tests_passed=true

cd forum-authentication && docker compose -f .github/test/unit/docker-compose.yaml up --build --abort-on-container-exit
cd ..
if [ $? -ne 0 ]; then
    all_tests_passed=false
fi

cd forum-authorization && docker compose -f .github/test/unit/docker-compose.yaml up --build --abort-on-container-exit
cd ..
if [ $? -ne 0 ]; then
    all_tests_passed=false
fi

cd forum-dialog && docker compose -f .github/test/unit/docker-compose.yaml up --build --abort-on-container-exit
cd ..
if [ $? -ne 0 ]; then
    all_tests_passed=false
fi

cd forum-thread && docker compose -f .github/test/unit/docker-compose.yaml up --build --abort-on-container-exit
cd ..
if [ $? -ne 0 ]; then
    all_tests_passed=false
fi

cd forum-user && docker compose -f .github/test/unit/docker-compose.yaml up --build --abort-on-container-exit
cd ..
if [ $? -ne 0 ]; then
    all_tests_passed=false
fi

cd forum-frame && docker compose -f .github/test/unit/docker-compose.yaml up --build --abort-on-container-exit
cd ..
if [ $? -ne 0 ]; then
    all_tests_passed=false
fi

cd forum-status && docker compose -f .github/test/unit/docker-compose.yaml up --build --abort-on-container-exit
cd ..
if [ $? -ne 0 ]; then
    all_tests_passed=false
fi

# cd forum-notification && docker compose -f .github/test/unit/docker-compose.yaml up --build --abort-on-container-exit
# if [ $? -ne 0 ]; then
#     all_tests_passed=false
# fi
#
# cd forum-support && docker compose -f .github/test/unit/docker-compose.yaml up --build --abort-on-container-exit
# if [ $? -ne 0 ]; then
#     all_tests_passed=false
# fi

cd forum-payment && docker compose -f .github/test/unit/docker-compose.yaml up --build --abort-on-container-exit
if [ $? -ne 0 ]; then
    all_tests_passed=false
fi

if [ "$all_tests_passed" = false ]; then
    echo "Some tests failed. Exiting with error."
    exit 1
else
    echo "All tests passed successfully."
    exit 0
fi