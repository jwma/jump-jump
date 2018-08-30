FROM node:alpine

WORKDIR /home/node/app
ADD $PWD/package.json ./package.json
RUN npm --registry https://registry.npm.taobao.org install