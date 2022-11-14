# build-stage
FROM node:12.14-alpine as build-stage

ENV PATH /app/node_modules/.bin:$PATH

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY package.json ./
COPY yarn.lock ./

RUN yarn install

COPY . ./

RUN yarn build


# deploy stage
FROM nginx:stable-alpine

WORKDIR /usr/share/nginx/html

COPY --from=build-stage /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
COPY CONTRIBUTING.* DCO LICENSE README.* ./
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

EXPOSE 80

HEALTHCHECK --interval=30s --start-period=10s --timeout=30s \
  CMD wget --quiet --tries=1 --spider "http://127.0.0.1:80/config.js" || exit 1

ENTRYPOINT ["/entrypoint.sh"]
