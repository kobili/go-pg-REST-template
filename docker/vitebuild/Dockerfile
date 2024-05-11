FROM node:21

WORKDIR /app

RUN mkdir frontend

WORKDIR frontend

COPY frontend .

RUN yarn install

CMD ["yarn", "build"]
