# Eurostat Weekly Deaths API

Eurostat publishes data for weekly deaths statistics among EU countries (weekly aggregates calculated for various age and gender groups).

This projects downloads this data, loads it from compressed TSV file to sqlite database and exposes it via small HTTP API (written in Golang).

## API routes

### Weekly Deaths
`/api/weekly_deaths` - returns weekly deaths statistics for given country, age bucket and gender, year range. Inputs are passed as query url parameters. Example call:

`/api/weekly_deaths?country=DE&gender=T&age=TOTAL&year_from=2015&year_to=2020`.

All parameters (`country`, `gender`, `age`, `year_from`, `year_to`) are **required**.

Example response:
```json
{
  "data": [
    {
      "age": "TOTAL",
      "country": "PL",
      "gender": "T",
      "week": 53,
      "weekly_deaths": 11775,
      "year": 2020
    },
    {
      "age": "TOTAL",
      "country": "PL",
      "gender": "T",
      "week": 52,
      "weekly_deaths": 12128,
      "year": 2020
    },
    {
      "age": "TOTAL",
      "country": "PL",
      "gender": "T",
      "week": 51,
      "weekly_deaths": 11679,
      "year": 2020
    }
  ]
}
```

List of country codes:

|Country Code|Name|
|---|---|
|AD|Andorra|
|AL|Albania|
|AM|Armenia|
|AT|Austria|
|BE|Belgium|
|BG|Bulgaria|
|CH|Switzerland|
|CY|Cyprus|
|CZ|Czechia|
|DE|Germany|
|DK|Denmark|
|EE|Estonia|
|EL|Greece|
|ES|Spain|
|FI|Finland|
|FR|France|
|GE|Georgia|
|HR|Croatia|
|HU|Hungary|
|IE|Ireland|
|IS|Iceland|
|IT|Italy|
|LI|Liechtenstein|
|LT|Lithuania|
|LU|Luxembourg|
|LV|Latvia|
|ME|Montenegro|
|MT|Malta|
|NL|Netherlands|
|NO|Norway|
|PL|Poland|
|PT|Portugal|
|RO|Romania|
|RS|Serbia|
|SE|Sweden|
|SI|Slovenia|
|SK|Slovakia|
|UK|United Kingdom|

List of age codes:
|Age Code|Name|
|---|---|
|TOTAL|Total|
|UNK|Unknown|
|Y_LT5|<5|
|Y5-9|From 5 to 9|
|Y10-14|From 10 to 14|
|Y15-19|From 15 to 19|
|Y20-24|From 20 to 24|
|Y25-29|From 25 to 29|
|Y30-34|From 30 to 34|
|Y35-39|From 35 to 39|
|Y40-44|From 40 to 44|
|Y45-49|From 45 to 49|
|Y50-54|From 50 to 54|
|Y55-59|From 55 to 59|
|Y60-64|From 60 to 64|
|Y65-69|From 65 to 69|
|Y70-74|From 70 to 74|
|Y75-79|From 75 to 79|
|Y80-84|From 80 to 84|
|Y85-89|From 85 to 89|
|Y_GE90|>=90|

List of gender codes:
|Gender Code|Name|
|---|---|
|T|Total|
|F|Female|
|M|Male|


### Labels

`/api/labels` returns list of all values and their labels for the data included in the database. Value of the "value" attribute should be used when querying `/api/weekly_deaths` endpoint. Endpoint serves all three types of labels: `age`, `gender`, `country`. 

You can use for example to populate dropdowns when working on visualizing the data.


Response example:
```json
{"data": 
    [
        {"value": "TOTAL", "label": "Total", "order": 1, "type": "age"},
        {"value":"UNK","label":"Unknown","order":2,"type":"age"},
        {"value":"FR","label":"France","order":16,"type":"country"},
        {"value":"GE","label":"Georgia","order":17,"type":"country"},
        {"value":"T","label":"Total","order":1,"type":"gender"},
        {"value":"F","label":"Female","order":2,"type":"gender"},
        {"value":"M","label":"Male","order":3,"type":"gender"}
    ]
}
```

## Running project locally

First, you need to populate the database. 

```
cd cmd/build_database
go run .
```

Above command fetches latest data from Eurostat website, parses it and creates a sqlite database in the repository root (`eurostat.db`).

Now you're ready to start the webserver:

```
cd cmd/web
go run .
```
