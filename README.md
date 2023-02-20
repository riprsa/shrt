# An URL Shortener

This project is my first attempt at a URL shortener. It's a simple web app that allows you to short any URL you want by API request. It's not perfect, but it's pretty good, I swear! It's also open source, so you can check out the code on this page.

## Usage

### Create Short URL

HTTP POST Request to <https://s.x16.me/short>

```json
{
    "url":"example.com/mypath"
}
```

Response:

```json
{
    "short":"AAAAAA"
}
```

Try it out:

```sh
curl -X POST -H "Content-Type: application/json" -d '{"url":"example.com/mypath"}' https://s.x16.me/short
```

### Use Short URL

HTTP GET Request will result in a HTTP 302 Redirect:

```sh
curl -IX GET https://s.x16.me/AAAAAA
HTTP/1.1 302 Found
Location: https://example.com/mypath
```

### Get Short URL

HTTP POST Request to <https://s.x16.me/url>

```json
{
    "short":"AAAAAA"
}
```

Response:

```json
{
    "url":"example.com/mypath"
}
```

Try it out too:

```sh
curl -X POST -H "Content-Type: application/json" -d '{"short":"AAAAAA"}' https://s.x16.me/url
```
