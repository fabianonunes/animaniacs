FROM node:10.13.0-slim

COPY . /app
WORKDIR /app
RUN yarn install

ENTRYPOINT [ "yarn" ]
CMD [ "start" ]
