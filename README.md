# Go Template

This is my personal stack that's based on what I've used to program when employed
to write go code. It's not supposed to be opinionated but rather a baseline to
start your projects, write less boilerplate code and be more productive overall.

## Setup

Do these steps to get it working on your machine. Alternatively, there's a
`compose.yml` that has profiles for both development and production environments,
so you only really need to have docker or podman installed.

### Tools

#### Go toolchain

- **[templ](https://templ.guide/)**: HTML template engine that's type-safe;
- **[sqlc](https://sqlc.dev/)**: Generate data mapping by only writing SQL queries.

You can run the commands below to download the tools:
```shell
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```
#### Sytem's package manager

- [Makefile](https://makefiletutorial.com/): building the application;
- Docker/Podman: scaffolding the whole service;
- Postgres: Database proven to be reliable;
- [bun](https://bun.sh/): Very fast bundler and package manager for javascript;
- [sass](https://sass-lang.com/): Tool to write and build reliable CSS;
- [UnoCSS](https://unocss.dev/): CSS utility framework similar to Tailwind;
- [htmx](https://htmx.org/): Javascript library to make HTML APIs interactive;
- [Alpine.js](https://alpinejs.dev/): Javascript framework for the client side.

#### Go libraries

The only one worth to look at their documentation is [chi](https://go-chi.io/#/README). The rest can be
looked on pkg.go.dev and encouraged to be replaced if it doesn't serve your
needs.

Keep in mind that some of them might be needed for some parts to work, so always
tests if the project still compiles when you decide to remove one.

### Steps

Run one of the commands to


#### Using docker/podman compose:
```shell
# Don't forget to edit .env when running in production
cp env.example .env

# Development with hot reload
$ docker compose --profile dev up

# Production
$ docker compose --profile prod up
```

#### Locally on your machine:

1. Clone the repository;
2. Rename the root package and all the module references (automation TBD);
3. Run `make server`;
4. Copy `env.example` as `.env`;
5. Install and initialize postgres database with the tables in `schema.sql`;
6. Export the database connection URL to the env (on linux you can run `export
$(cat .env | xargs)`);
7. Run the `wserver` executable to start the server.

## Documentation

Open the `doc.go` file inside the package, it should have everything you need to
know to get started on programming. Don't shy away from reading the source code
files, I tried to program and write comments so it's easy to understand by just
looking at the code. If you feel it's overwhelming, please open an issue.

## Very professional todo list

- Expand the blog example more to be a full CRUD;
- Maybe improve the asset manager cache system;
- Support more compression algorithms (brotli, deflate);
- TOTP/email verification;
- Write tests for some of the modules.
