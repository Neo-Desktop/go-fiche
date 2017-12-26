## GoFiche!
![Logo](https://neo-desktop.github.io/go-fiche/logo.png)

### A Go fork of [Fiche](https://github.com/solusipse/fiche "fiche")

Command line pastebin for sharing terminal output.

-----

# Server-side usage

## Installation From Source

1. Clone:

    ```
    git clone https://github.com/Neo-Desktop/go-fiche.git
    ```

2. Build:

    ```
	go get github.com/ahmetb/govvv
    govvv build
    ```
    
3. Install:

    ```
    cp go-fiche /usr/local/bin
    ```

-------------------------------------------------------------------------------

## Usage

```
Usage of go-fiche:
  -B, --buffer int        This parameter defines size of the buffer used for getting data from the user. Maximum size (in bytes) of all input files is defined by this value. (default 32768)
  -d, --domain string     This will be used as a prefix for an output received by the client. Value will be prepended with http[s]. (default "localhost")
  -h, --help              Prints this help message
  -S, --https             If set, Go-Fiche returns url with https prefix instead of http.
  -l, --log string        Log file. This file has to be user-writable.
  -o, --output string     Relative or absolute path to the directory where you want to store user-posted pastes. (default "./code")
  -p, --port int          Port in which the service should listen on. (default 9999)

```

These are command line arguments. You don't have to provide any of them to run the application. Default settings will be used in such case. See section below for more info.

### Settings

-------------------------------------------------------------------------------

#### Output directory `-o`

Relative or absolute path to the directory where you want to store user-posted pastes.

```
go-fiche -o ./code
```

```
go-fiche -o /home/www/code/
```

__Default value:__ `./code`

-------------------------------------------------------------------------------

#### Domain `-d`

This will be used as a prefix for an output received by the client.
Value will be prepended with `http`.

```
go-fiche -d domain.com
```

```
go-fiche -d subdomain.domain.com
```

```
go-fiche -d subdomain.domain.com/some_directory
```

__Default value:__ `localhost`

-------------------------------------------------------------------------------

#### Slug size `-s`

This will force slugs to be of required length:

```
go-fiche -s 6
```

__Output url with default value__: `http://localhost/xxxx`,
where x is a randomized character

__Output url with example value 6__: `http://localhost/xxxx`,
where is a randomized character

__Default value:__ 4

-------------------------------------------------------------------------------

#### HTTPS `-S`

If set, fiche returns url with https prefix instead of http

```
go-fiche -S
```

__Output url with this parameter__: `https://localhost/xxxx`,
where x is a randomized character

-------------------------------------------------------------------------------

#### Buffer size `-B`

This parameter defines size of the buffer used for getting data from the user.
Maximum size (in bytes) of all input files is defined by this value.

```
go-fiche -B 2048
```

__Default value:__ 32768

-------------------------------------------------------------------------------

#### Log file `-l`

```
go-fiche -l /home/www/fiche-log.txt
```

__Default value:__ not set

__WARNING:__ this file has to be user-writable