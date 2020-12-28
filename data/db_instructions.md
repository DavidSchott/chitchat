# Database Instructions

## Creating the database
Assuming default Postgress installation, you can create a DB called "chitchat" via:
```
createdb -h localhost -p 5432 -U postgres chitchat
```

## Creating the database tables
To create the database tables used by ChitChat, use [create_tables.sql](./create_tables.sql):
```
psql -f create_tables.sql -d chitchat
```