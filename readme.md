# Authme

Have you ever felt that your authentication logic was in the way of you
application? I did.

Authme is an application in 200 lines of code which allows the
registration, log in and authentication of users.
It uses Argon2 for password encryption which is totally rad nowadays.

## Usage

Authme is supposed to be used as middleware for Nginx. Configure Nginx like so:

```
server {
	listen 80; # You of course use TLS, but let's not complicate matters here.
	server_name localhost;

	error_page 401 = @error401;
	location @error401 {
		return 301 /login.html;
	}

	location /secure/ {
		auth_request /authenticated;
	}

	location = /authenticated {
		internal;
		proxy_pass http://localhost:8080;
	}

	location = /login {
		proxy_pass http://localhost:8080;
	}

	location = /register {
		proxy_pass http://localhost:8080;
	}
}
```

Note that there's a difference between `login` and `login.html`. The former is
the POST target for the login form. The login form should be found at
`login.html`. The same is true for `register` and `register.html`.

We've used the `internal` directive in the `authenticated` location. That means
this endpoint is only meant for internal Nginx usage.

As we don't want our authentication middleware to concern itself with how
authentication failures are handled we use the named `@error401` to turn a 401
into a 301.

## Storage

SQlite3 is used as a storage backend. That should be good for 90% of the use
cases. And if it isn't; the Sqlite3 driver is `database/sql` compatible so
you're one import line away from any other database.

Initializing databases was never so sexy:

```
$ sqlite3 users.db < schema.sql
```
