# Schemata

`database/schemata.js` contains all schemata for realm. An schema is defined as an JS object. Have a look at the realm doc's in order to figure out the prop's. A flow type is defined for every schema in order to represent the data. Nameingconvensions for the schema variables are "$SchemaName+Schema" e.g. "AccountBalanceSchema". The nameing convensions for the types that represent the data is: "$SchemaName+Type" e.g. "AccountBalanceType".

Current schemata are:
- ProfileSchema (information about the profile)
- AccountBalanceSchema (account balance of an specific eth address)
- MessageJobSchema (all messaging jobs)