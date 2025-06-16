# Этап сборки
FROM node:24-alpine3.21 AS builder

WORKDIR /app

COPY frontend/package*.json ./

RUN npm install

COPY frontend/ .

RUN npm run build

FROM nginx:1.27.5-alpine

COPY --from=builder /app/build /usr/share/nginx/html

COPY ./cloud/frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
