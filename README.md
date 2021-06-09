# Checkout Project
This is a simple project for checkout items

## Description
User will send the requested items, system will calculate the price and list of items that user will received

## Pre-requisites
- golang

## Usage

`$ make run-test` // for run test only

`$ make build` // build and generate binary

`$ make run` // run the test, generate binary then execute

## Steps

- clone this repo 
- open terminal in the project folder
- run `$ make run` on the terminal
- send request in Postman with URL `http://localhost:10000/checkout` and method `POST`

## Request Body Example
```json
{
    "items": [
        {
            "sku": "120P90",
            "qty": 1
        }
    ]
}
```

## Response Body Example
```json
{
    "total": 49.99,
    "items": [
        {
            "sku": "120P90",
            "qty": 1
        }
    ]
}
```

## Available SKUs
- GOOGLE_HOME    = `120P90`
- MACBOOK_PRO    = `43N23P`
- ALEXA_SPEAKER  = `A304SD`
- RASPBERRY_PI_B = `234234`
