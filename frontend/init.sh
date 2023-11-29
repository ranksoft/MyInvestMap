#!/bin/sh

if [ ! -d "node_modules" ]; then
    npx create-react-app .
    npm install
else
    npm install
fi

npm start
