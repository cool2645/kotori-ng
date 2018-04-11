# kotori-ng

Kotori-ng is the new generation of blog system [kotori](https://github.com/rikakomoe/kotori).

To develop a plugin for kotori-ng, refer to [kotori-ng-sampleplugin](https://github.com/cool2645/kotori-ng-sampleplugin).

API documentations are available at [here](https://app.swaggerhub.com/apis/cool2645/kotori-ng/1.0.0).

To run this under the same domain of your front-end (assumed as an single page application), maybe this will help:

```nginx
        location /api {
                rewrite ^/api/(.*) /$1 break;
                proxy_pass  http://127.0.0.1:2233;
                proxy_redirect / /api/;
        }

        location / {
                root /var/www/front-end-of-my-blog;
                index index.html index.htm;
                try_files $uri $uri/ /index.html;
        }
```
