# Authme

Have you ever felt that your authentication logic was in the way of you
application? I did.

Authme is an application in 200 lines of code which allows the
registration, log in and authentication of users.
It uses Argon2 for password encryption which is totally rad nowadays.

## Usage

Authme is supposed to be used as middleware for Nginx:


```
server {
	listen 80; # You of course use TLS, but let's not complicate matters here.
	server_name localhost;

    error_page 401 = @error401;
    location @error401 {
            return 301 /login;
    }

    location /secure/ {
            auth_request /authenticated;
    }

    location = /authenticated {
            internal;
            proxy_pass http://localhost:8080;
    }

    location = /login {
            internal;
            proxy_pass http://localhost:8080;
    }

    location = /register {
            internal;
            proxy_pass http://localhost:8080;
	}
}
```

