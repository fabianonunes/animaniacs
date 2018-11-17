FROM node:10.13.0-slim

WORKDIR /animaniacs

COPY package.json yarn.lock ./
RUN yarn install

COPY app app

USER node
ENTRYPOINT [ "yarn" ]
CMD [ "start" ]
