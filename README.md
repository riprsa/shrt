# An URL Shortener

This project is my first attempt at a URL shortener. It's a simple web app that allows you to short any URL you want by API request. It's not perfect, but it's pretty good, I swear! It's also open source, so you can check out the code on this page.

## Usage

### Create a short URL

```sh
curl -X POST -H "Content-Type: application/json" -d '{"url":"example.com/mypath"}' https://s.x16.me/api/short
```

### Get a short URL

```sh
curl -X POST -H "Content-Type: application/json" -d '{"short":"AAAAAA"}' https://s.x16.me/api/url
```

### Redirect from a short URL

```sh
curl -IX GET https://s.x16.me/AAAAAA
```
