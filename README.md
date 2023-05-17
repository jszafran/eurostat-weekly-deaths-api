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

```

List of available countries:

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
### Labels

`/api/labels` returns list of all values and their labels for the data included in the database. Value of the "value" attribute should be used when querying `/api/weekly_deaths` endpoint. Endpoint serves all three types of labels: `age`, `gender`, `country`. 

You can use for example to populate dropdowns when working on visualizing the data.


Response example:
```json
{"data": [
    {"value": "TOTAL", "label": "Total", "order": 1, "type": "age"},
    {"value":"UNK","label":"Unknown","order":2,"type":"age"},
    ...
    {"value":"FR","label":"France","order":16,"type":"country"},
    {"value":"GE","label":"Georgia","order":17,"type":"country"},
    ...
    {"value":"T","label":"Total","order":1,"type":"gender"},
    {"value":"F","label":"Female","order":2,"type":"gender"},
    {"value":"M","label":"Male","order":3,"type":"gender"}
]}
```