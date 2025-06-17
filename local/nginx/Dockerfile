FROM nginx:1.27.5-alpine

RUN apk add openssl

COPY nginx.conf /etc/nginx/nginx.conf
COPY htpasswd.sh /usr/local/bin/htpasswd.sh
RUN chmod +x /usr/local/bin/htpasswd.sh

EXPOSE 80

ENTRYPOINT ["/bin/sh", "-c"]

CMD ["/usr/local/bin/htpasswd.sh && nginx -g 'daemon off;'"]