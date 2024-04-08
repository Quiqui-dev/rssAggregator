#!/bin/bash

DB_URL=$DB_URL

goose postgres "$DB_URL" up