# Demo app for Trivago Case Study
This repository contains the demo app used for Trivago case study
Is a simple golang app builded from scratch to expose prometheus metrics and a connection to MySQL

## Installation
- Follow the instructions on https://github.com/rferrerb/trivago-case-study to build de needed infraestructure.
- Create the DB and populate it, you can use the */trivago-app/db/database_deployment.sql*
- Build the docker image with the provided Dockerfile

## Test the app
The app is listening by default on port 8080.
The metrics are exposed on "/metrics" path
To test the DB connection access to "/employees"
