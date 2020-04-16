FROM node:13.12.0-alpine AS compile-image

WORKDIR /opt/ng
COPY package.json ./
RUN npm install

ENV PATH="./node_modules/.bin:$PATH" 

COPY . ./
RUN ng build --prod

FROM nginx:1.17.9-alpine
#COPY docker/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=compile-image /opt/ng/dist/bee /usr/share/nginx/html