# 🔗✂️ Super URL shorter

This project is my first attempt at a URL shortener. It's a simple web app that allows you to short any URL you want. It's not perfect, but it's pretty good, I swear! It's also open source, so you can check out the code on this page.

## Features

It can short any URL you want from web interface and from API request. Data of URLs is stored in PostgreSQL database.

To use API go to ```hostname/api```

## How to setup database?

First, you need to set up a PostgreSQL database. You can use the following guide to do that: [PostgreSQL tutorial](https://www.postgresql.org/docs/14/tutorial-start.html).

After setup your database and user, you need to create a database that will be used by this app. You can use this SQL command to create a database:

```CREATE TABLE links (
    id SERIAL,
    url VARCHAR,
    short VARCHAR
);

```

## How to setup app?

If you want to use GitHub actions, you need to fork this repo and add to your repo's secrets this:

```HOST=\<your serves's IP or hostname>
PORT=\<your server's port>
SSHKEY=\<your ssh key>
USERNAME=\<your server's username>

```

Then, you need to setup GitHub workflow to run this app. You can do it by yourself (sorry:3)

Or you can deploy on your server by [Docker](https://www.docker.com/).
You should create a folder and create .env file in it. File will contain the following:

```DB_USERNAME=\<your DB username>
DB_PASSWORD=\<your DB password>
DB_HOSTNAME=\<your DB hostname>
DB_NAME=\<your DB name>
DB_MODE=disable # or enable
PORT=\<port of app>

```

Then, you can run this in folder: (but I am not sure that here is right way to run Docker image)

```docker run -p 80:80 -d haraldka/shorter```

## Why I wrote this?

I needed some expirience with web programming and I wanted to try something new. I wanted to make a project that I can use in my other own projects. I think, I did it, but not that clear that I imagine. Thank you for attention!
