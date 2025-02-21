# Go Template

This is my personal stack that's based on what I've used to program when employed
to write go code. It's not supposed to be opinionated but rather a baseline to
start your projects, write less boilerplate code and be more productive overall.

## Setup

Install all the necessary tools in bold using either the go toolchain or the
package manager of your choice.

### Tools

- **Makefile**: building the application;
- **Docker/Podman**: scaffolding the whole service;
- Postgres: Database proven to be reliable;
- **templ**: HTML template engine that's type-safe;
- **bun**: Very fast bundler and package manager for javascript;
- **sass**: Tool to write and build reliable CSS;
- picocss: CSS framework that provides sane defaults;
- htmx: Javascript library to make HTML APIs interactive;
- Alpine.js: Javascript framework for the client side;
- **sqlc**: Generate data mapping by only writing SQL queries.

### Steps

1. Clone the repository;
2. Rename the root package and all the module references (automation TBD);
3. Run `make server`;
3. Copy `env.example` as `.env`;
4. Start the compose service using either docker or podman;
5. Run `schema.sql` to initialize the database;
6. Export the database connection URL to the env (on linux you can run `export
$(cat .env | xargs)`);
7. Run the `wserver` executable to start the server.

## Documentation

Open the `doc.go` file inside the package, it should have everything you need to
know to get started on programming. Don't shy away from reading the source code
files, I tried to program and write comments so it's easy to understand by just
looking at the code. If you feel it's overwhelming, please open an issue.

## Very professional todo list

- Decide authentication library;
- Expand the blog example more to be a full CRUD;
- Automate database table creation (aka. remove step 5);
- Implement live reload mechanism (very hard);
- Write a Dockerfile to be fully containerized;
- Maybe improve the cache system;
- Support more compression algorithms (brotli, deflate);
- Write tests for some of the modules.
