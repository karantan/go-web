# go-web
My personal template for building web apps with Go and Elm

## Development

Start the Postgres DB and connect to it run `make db` or `docker compose up db`. This
will run a docker container with a running postgres service. To connect to it via
terminal type:

```
$ pgcli postgres://postgres:secret@localhost:5432/goweb
```

## Useful resources

- [Style guideline for Go packages](https://rakyll.org/style-packages/)
- https://github.com/gothinkster/golang-gin-realworld-example-app/
- https://gorm.io/docs/index.html
- https://github.dev/demo-apps/go-gin-app


https://hoohoo.top/blog/20210530112304-golang-tutorial-introduction-gin-html-template-and-how-integration-with-bootstrap/