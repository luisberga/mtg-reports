Magic The Gathering Cards Management Reports
=======================================

Overview
--------

This project allows users to manage a collection of Magic The Gathering (MTG) cards. The project consists of three applications:

1.  An API to manage the cards.
2.  A conciliation application called `conciliateJob`, which updates card prices from the Scryfall API.
3.  A reporting application called `reportJob`, which generates a report of the top 100 cards that most changed price and send it by email.

API Usage
---------

The API exposes several endpoints for card management. Here are the routes provided:

-   POST `/insert-card`: Inserts a single card into the database.
-   POST `/insert-cards`: Inserts multiple cards into the database in bulk.
-   GET `/card/{id}`: Retrieves a card by its ID.
-   GET `/cards`: Retrieves cards filtered by set name, card name, or collector number.
-   DELETE `/delete-card/{id}`: Deletes a card by its ID.
-   GET `/card-history/{id}`: Retrieves the price history of a card by its ID.
-   PATCH `/update-card/{id}`: Updates a card by its ID.

### Request and Response Formats

Details can be found in swagger file in `/docs/swagger.yaml`

[Click here to access the API documentation](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/luisberga/mtg-reports/main/docs/swagger.yaml)

### Bulk insert Cards

The `insert-cards` endpoint expects a POST request with a file attached. The file must be named cards.txt and should contain multiple entries, each in the following format:

`name: card name, set_name: set name, collector_id: collector number, foil: boolean`

Example: 

`name: Samwise the Stouthearted, set_name: ltr, collector_id: 449, foil: true`

The response includes the count of processed and unprocessed entries:

```YAML
{
  "processed": processed count,
  "not_processed": not processed count
}
```

Errors
------

In case of an error, the API will return a response with an appropriate HTTP status code and a JSON body with an error message. The jobs will log the error.

Running the Applications
------------------------

To run the applications, follow the steps below:

1.  Generate the `config.yaml` file:

    `make generate-config`

    This command will generate the `config.yaml` file required for the applications to work correctly.

2.  Build and start the API and jobs containers:

    `make build-up`

    This command will build the Docker containers for the API and the two jobs (`conciliateJob` and `reportJob`) and start them in the background.

3.  Access the API and manage the cards: The API will be accessible at `http://localhost:8080`.

4.  Run the `conciliateJob` to update card prices:

    `make conciliate-cards`

    This command will run the `conciliateJob`, which will update the card prices in the database by fetching data from the Scryfall API.

5.  Run the `reportJob` to generate the top 20 most expensive cards report:

    `make report-top-cards`

    This command will run the `reportJob`, which will generate a report with the top 20 most expensive cards and display the results.

6.  Stop and remove the containers (when finished):

    `make down`

    This command will stop the Docker containers for the API and the jobs.


Project Status
--------------

The project is ongoing. Future updates may include adding a queuing system for the `insert-cards` endpoint to handle large volumes of data.

SMTP Email and Exchange Rate
----------------------------

The job for sending emails via SMTP requires you to have an SMTP server account. Please make sure to set up your SMTP server credentials in the `config.yaml` file.

Additionally, the application utilizes the `exchangerate-api` to get the exchange rate for the value of the dollar to the Brazilian Real (BRL). By default, the exchange rate is set to 5 BRL (Brazilian Real) to 1 USD (US Dollar). If you prefer not to use the `exchangerate-api`, you can modify this value as a constant within the code.

Note for ARM Architecture Users
-------------------------------

If you are using a computer with ARM architecture, make sure to use the following image for the database in your `docker-compose.yaml` file:

```YAML
db:
  image: yobasystems/alpine-mariadb
  # other configurations
```

Probably you will need to start the database schema manually. Is is located in `migrations/ddl`.

Contributions
-------------

Contributions, issues, and feature requests are welcome. Feel free to check the issues page if you want to contribute.

Author
------

Luis Felipe de Oliveira Bergamim

License
-------

MIT License