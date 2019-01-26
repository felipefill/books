# Books 

[![Build Status](https://travis-ci.org/felipefill/books.svg?branch=master)](https://travis-ci.org/felipefill/books)

Books is a serverless GO app to create, search and scrap books.

This project was written in [GO](https://golang.org/) and uses [AWS Lambda](https://aws.amazon.com/lambda/) to serve its endpoints. 
The database is also hosted by Amazon ([RDS](https://aws.amazon.com/rds/)).

I decided to use the [serverless](https://serverless.com/) framework in order to facilitate and speed up development.

I'm using PSQL here for data and also [GORM](http://gorm.io/) for handling data in the app, you should be able to easily switch to any database 
supported by GORM's [dialects](http://gorm.io/docs/dialects.html).

## Endpoints

### Create

This endpoint receives a JSON representing a book and stores it in the database. JSON should look like this:

```
{
  "title": String,
  "description": String,
  "isbn": String,
  "language": String
}
```

### Search by id

Given a specific ID (passed using path parameter) searches the database and replies with a JSON like this:

```
{
  "id": Integer,
  "title": String,
  "description": String,
  "isbn": String,
  "language": String
}
```

### Search in website

This endpoint can work in three different ways:

1. `retrieve_all` (default): retrieves all books in the database;
2. `scrap_and_store`: visits a website and scraps it looking for new books then stores then in database and return all books in the database;
3. `scrap_only`: visits a website and scraps it looking for new books and returns them.

If you want to use the default mode then no additional action is required when calling the endpoint.

In order to use `scrap_and_store` or `scrap_only` you will need to set the `mode` parameter in query string to either.

Response looks like this:

```
{
  "numberBooks": Integer,
  "books":[
    {"..."},
    {"..."},
    {"..."}
  ]
}
```

Note: when I was almost done with this project I found out that because this uses [API Gateway](https://aws.amazon.com/api-gateway/) the maximum timeout is 30 seconds. This might afect the scrapping modes but it's very unlikely that it'll run for more than that.

## Setup

### Dependencies

- GO
- Dep (go dependency manager)
- Serverless
- awscli (configured)

### Installing and configuring

You can usually install it by using a package manager, e.g.:
```
brew install go dep serverless awscli
```

You will have to clone this project inside a proper set up [GOPATH](https://golang.org/doc/code.html#GOPATH).

As for `awscli`, you must have it configured with credentials (that have access to Lambda related stuff):
```
aws configure
```

Last but not least, you will need to write the database info to a `serverless.env.yml` file. There's a sample included in this repo.

You can build, test and deploy using [make](https://en.wikipedia.org/wiki/Make_(software)):

## Build, test and deploy

I've configured these in the `Makefile`, you can run them using:

```
make build # Builds the project
make test # Run all the tests and shows code coverage
make deploy # Deploys to AWS Lambda
```

