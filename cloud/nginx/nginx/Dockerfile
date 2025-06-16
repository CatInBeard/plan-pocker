FROM nginx:1.27.5-alpine

RUN apk add openssl

COPY nginx.conf /etc/nginx/nginx.conf
COPY index.html /admin

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]