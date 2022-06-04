# homework
[TOC]

### Clone repo

```shell=
git clone https://github.com/fredyyu/homework.git
```

#### If it's the first time clone
### Modify environment variables

1. Copy `.env.example` file to `.env` in repository
2. Modify the `.env` content before start the service
3. You can use `docker-compose config` to check your configuration

### Boot up containers

==Running in foreground==
```shell=
docker-compose up
```

==Running in background==
```shell=
docker-compose up -d
```

## Check Services

```shell=
docker-compose ps
```

get in to database container

```shell=
docker-compose exec -ti <SERVICE_NAME> bash
```

### Stop service and ==KEEP DATA==

```shell=
docker-compose down
```

### Clean whole data

```shell=
docker-compose down -v
```

### Build the service image

At the homework root folder run the command

```shell=
$ docker build -t homework -f ./build/Dockerfile . --no-cache
```