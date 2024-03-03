#!/bin/bash

# Executes a goose migration to a postgres db 
# take a command line argument of up or down and the .env file for the environment
# then execute the migrations 


migration_type=$1
env_file=$2

export $(grep -v '^#' $env_file | xargs)

DB_URL=$(echo $DB_URL | cut -d'?' -f 1)

echo "$migration_type migration selected"

echo "Moving to schema directory"

cd ./sql/schema

echo "executing $migration_type migration"

goose postgres $DB_URL $migration_type