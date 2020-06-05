FROM node:13.12.0-alpine AS compile-image

WORKDIR /opt/app
COPY package.json yarn.lock ./
RUN yarn
COPY . ./
RUN yarn build
RUN yarn test a --watchAll=false

FROM nginx:1.17.9-alpine
#COPY docker/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=compile-image /opt/app/build /usr/share/nginx/html