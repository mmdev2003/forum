#!/bin/bash

apt install nginx -y

rm /etc/nginx/nginx.conf
cat <<EOF >>/etc/nginx/nginx.conf
user www-data;
worker_processes auto;
pid /run/nginx.pid;
error_log /var/log/nginx/error.log;

events {
    worker_connections 768;
}


http {
    server_names_hash_bucket_size 128;
    sendfile on;
    tcp_nopush on;
    client_max_body_size 2048M;
    types_hash_max_size 2048;
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    ssl_protocols TLSv1 TLSv1.1 TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;
    access_log /var/log/nginx/access.log;
    gzip on;

    server {
        server_name _;
        listen 80;

        location /api/authentication {
            proxy_pass http://127.0.0.1:8000;
        }

        location /api/authorization {
            proxy_pass http://127.0.0.1:8001;
        }

        location /api/dialog {
            proxy_pass http://127.0.0.1:8002;
        }

        location /api/thread {
            proxy_pass http://127.0.0.1:8003;
        }

        location /api/frame {
            proxy_pass http://127.0.0.1:8004;
        }

        location /api/status {
            proxy_pass http://127.0.0.1:8005;
        }

        location /api/user {
            proxy_pass http://127.0.0.1:8006;
        }

        location /api/notification {
            proxy_pass http://127.0.0.1:8007;
        }

        location /api/support {
            proxy_pass http://127.0.0.1:8008;
        }

        location /api/payment {
            proxy_pass http://127.0.0.1:8009;
        }
    }
}
EOF
systemctl restart nginx

#apt install snapd
#snap install --classic certbot
#ln -s /snap/bin/certbot /usr/bin/certbot
#certbot --nginx