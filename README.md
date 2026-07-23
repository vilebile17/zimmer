# Zimmer

![Zimmer Slogan](/pics/zimmer_slogan.png)

<p align="center">
  <a href="https://skillicons.dev">
    <img src="https://skillicons.dev/icons?i=go,js,html,css,postgres,aws,linux" />
  </a>
</p>

## Table of Contents

- [About](#about)
- [Usage](#usage)
  - [Initial Steps](#initial-steps)
  - [PostgreSQL – DB setup](#postgresql---db-setup)
  - [Goose](#goose)
  - [You’re almost done!](#youre-almost-done)
- [Thank you](#thank-you)

## About

Zimmer is an [LMS](https://en.wikipedia.org/wiki/Learning_management_system) built in Golang for my [boot.dev](https://boot.dev)
[capstone project](https://www.boot.dev/courses/build-capstone-project) as a way to put all that I have learnt into one _giga-project_.

Zimmer is relatively simple in terms of features. Zimmer does only what is necessary:

- Class creation
- Assignment creation
- Assignment submissions and grading them
- Resources
- Announcements

This project was supposed to be a **backend** project, however, as I drew to the end of the backend development stage
I knew that this would need some user interface. And of course, a website makes the most sense for a platform like this.
Along the way I learnt _loads_ about frontend development and would like to mention a few _key resources_ that helped me
create the website

- [W3Schools](https://www.w3schools.com/howto/default.asp)
- [boot.dev JS](https://www.boot.dev/courses/learn-javascript)
- [Sajid](https://www.youtube.com/@whosajid) - YT channel on general css tips
- [Coding2Go](https://www.youtube.com/watch?v=wsTv9y931o8) - YT video on flex boxes

Without these resources it would be near impossible for me to make the site so I feel it's necessary to acknowledge them here :)

## Usage

The site is currently being hosted on [AWS](https://aws.amazon.com/) and you can check out the website: [zimmer.vilebile.dev](https://zimmer.vilebile.dev)
if you're just looking to use the platform.

If however, you're looking to **host the site locally**, there's a _few more step to do..._

### Initial Steps

Let's get the easy stuff done first: 

1. Install [Golang](https://go.dev/) and [PostgreSQL](https://www.postgresql.org/)
2. Clone the Repo

```bash
git clone https://github.com/vilebile17/zimmer
cd zimmer
psql --version
go version
```

### PostgreSQL - DB setup

With PostgreSQL installed we need to start the service in the background, use the correct command for your OS:

- **macOS:** `brew services start postgresql@15`
- **Linux (Debian/Ubuntu):** `sudo service postgresql start`
- **Linux (systemd):** `sudo systemctl start postgresql`

Enter the psql shell:

- macOS: `psql postgres`
- Linux: `sudo -u postgres psql`

Create the database:

```sql
CREATE DATABASE zimmer;
```

You may also set a password using the following command:

```
ALTER USER postgres WITH PASSWORD 'PASSWORD';
```

### Goose

Next we would like to actually populate the database with the correct schema. This can be
done using [Goose](https://github.com/pressly/goose)

You can install it using the command:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

It may take a while to compile the binary, if you don't want to wait, you can download a binary
from their [github releases](https://github.com/pressly/goose/releases)

Next we'll need your connection string. This is of the format `protocol://username:password@host:port/database`
So for us it'll be something like `postgres://postgres:PASSWORD@localhost:5432/zimmer`

> [!NOTE]
> If you didn't set a password just leave the `PASSWORD` parameter blank.
> I.e. `postgres://postgres:@localhost:5432/zimmer`

Finally `cd` into the `sql/schema` directory and run the `goose` migrations:

```bash
cd ./sql/schema
goose postgres postgres://postgres:PASSWORD@localhost:5432/zimmer up
```

### You're almost done!

With all that in place, you just need to make an `.env` file. You can copy `example.env`
to `.env` by running `cp example.env .env` in the root directory

Make sure to change the `DB_URL`'s password to your actual database password

> [!IMPORTANT]
> Do **not** keep the default `JWT_SECRET` value as that would be a **massive security risk.**
> You can run `openssl rand -base64 64` to generate a new one.

> [!NOTE]
> If you want to run in `https`, set the `PORT` environment variable to `443` and ensure that the `DOMAIN` variable
> points to the correct domain name. With `DOMAIN="localhost"` we use **self-signed certificates** while on everything else we use [Let's encrypt](https://letsencrypt.org/)

Once you've done that you're **good to go!** Just run `go run .` in the root directory and you should get
the success message `Hosting Zimmer at http://localhost:PORT`, head to that address and (hopefully) it
should work just fine. 

## Thank you

If you seriously went through all of those steps and successfully installed zimmer, **I greatly appreciate it!**

If you're looking to contribute to Zimmer, you're more than welcome to open a [pull request](https://github.com/vilebile17/zimmer/pulls)

Or if you find a bug and don't know how to fix it yourself, you can simply create an [issue](https://github.com/vilebile17/zimmer/issues/new)
