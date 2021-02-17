#!/bin/bash

yarn install
yarn typeorm migration:run
yarn console fixtures
yarn start:dev