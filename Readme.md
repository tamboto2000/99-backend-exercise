# 99 Backend Excercise
A technical test for Backend Developer

## How to run
This is a step-by-step sfdsf to run the services stack.

It is required that you run this stack on Unix or Linux based system with Docker and Docker Compose installed.

### 1. Set the configs
To set the configs, run this command
```sh
make prepare-configs
```

This command will generate 2 set of configs: service configs and deployment configs.

`configs/` directory contains the service configs, which is `user_svc.yaml` config for User Service, and `public_api_svc.yaml` for Public API Service. You can leave this configs as-is or you can modify it.

`deployment/` contains config for deployment, one file that you must pay attention to is `docker-compose.env`. In that file, you must set the `LISTING_GIT_USERNAME` and `LISTING_GIT_TOKEN` to your github username and access token for [99-backend-exercise repository](https://github.com/team99-exercise/99-backend-exercise)

The rest of the configs you can leave it as-is, or you can modify it

### 2. Deploy!
Run this command to delpoy the whole stack:
```sh
make deploy
```

If you don't change the Public API Service HTTP port, you can access it from `localhost:8001`, or access it with any other port that you set. The endpoints for the public APIs are as required in the [99-backend-exercise repository](https://github.com/team99-exercise/99-backend-exercise)


To stop you can run this command:
```sh
make shutdown
```

## Libraries and tools
The modules or libraries that lives inside `pkg` is pulled out of my side project or from other technical tests that I have done, only small amount of direct third-party libraries is used for this technical test
