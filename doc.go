package main

/*
Toggl is a service that implements a RESTful API for questions

The API includes the following endpoints:
- GET /questions 			returns all questions in the database
- GET /questions/{id}       seek pagination, where limit is provided in requests body
- POST /questions/          create new question
- PUT /questions/{id}       update existing question
- DELETE /questions/{id}    delete existing question

All endpoints require authorization over jwt
JWT uses HS256 for generating token with secret key: toggl

Neccessery header:
Authorization: <token>

RUN INSTRUCTIONS:
 1. sqlite3 toggl.db < sql/startup.sql
 2. docker build -t toggl .
 3. docker-compose up -d

 Service will be run on env "PORT" variable, if not defined default port is 3000
*/
