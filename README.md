# user-crud

Basic CRUD for user data management

# Building
This service is built using Go 1.20. It uses gin for routing and Mongo Db as database.


```
make build
```

# Environment Variables

The service requires the following environment variables in order to connect to MongDB and set it's api key for authentication.

```
API_KEY
MONGO_USERNAME
MONGO_PASSWORD
MONGO_URI
```

# Development instructions

To run the app in the development mode, init the proper containers using docker compose.

```
docker compose up -d
```

This will start a mongo db and mongo express containers with testing credentials.

You may then start the app using the command

```
make run
```

To execute the tests use the unit-test command

```
make unit-test
```

# API

This section documents the API endpoints.

## Notes

* Password field is provided in plaintext and only the hash is returned in the responses. The service takes care of the hashing before storage
* Dates use the format "YYYY-MM-DD"
* E-mail is unique for each user and the format is validated
* UUID must be compliant
* Authentication is done via API KEY. X-API-KEY should be added to the request header

___


## Create User

Get user details

**URL** : `/api/user`

**Method** : `POST`

**Params**

```json
{
    "name": "John Doe",
    "birth_date": "1970-01-02",
    "email": "joe25@mailprovider.com",
    "address": "3197 Woodrow Way",
    "password": "my secret password"
}
```


## Response

**Code** : `201 Created`, `400 Bad request`

**Content examples**

Succesful response example.

```json
{
    "name": "John Doe",
    "uuid": "d035e79d-ffe9-4ebf-b665-747353b3ea40",
    "birth_date": "1970-01-02",
    "email": "joe25@mailprovider.com",
    "address": "3197 Woodrow Way",
    "password": "6d795f5365637265745f70617373e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
}
```


___


## Get User

Get user details

**URL** : `/api/user/{uuid}`

**Method** : `GET`

## Response

**Code** : `200 OK`, `400 Bad request`, `404 Not found`

**Content examples**

Sucessful Response example.

```json
{
    "name": "John Doe",
    "uuid": "d035e79d-ffe9-4ebf-b665-747353b3ea40",
    "birth_date": "1970-01-02",
    "email": "joe25@mailprovider.com",
    "address": "3197 Woodrow Way",
    "password": "6d795f5365637265745f70617373e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
}
```

___


## Update User

Get user details

**URL** : `/api/user/{uuid}`

**Method** : `PATCH`

**Params**

```json
{
    "name": "John Doe",
    "birth_date": "1970-01-02",
    "email": "joe25@mailprovider.com",
    "address": "3197 Woodrow Way",
    "password": "6d795f5365637265745f70617373e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
}
```


## Response

**Code** : `200 OK`, `400 Bad request`, `404 Not found`

**Content examples**

Succesful response example.

```json
{
    "name": "John Doe",
    "uuid": "d035e79d-ffe9-4ebf-b665-747353b3ea40",
    "birth_date": "1970-01-02",
    "email": "joe25@mailprovider.com",
    "address": "3197 Woodrow Way",
    "password": "6d795f5365637265745f70617373e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
}
```

___

## Delete User

Get user details

**URL** : `/api/user/{uuid}`

**Method** : `DELETE`

## Response

**Code** : `200 OK`, `400 Bad request`, `404 Not found`
