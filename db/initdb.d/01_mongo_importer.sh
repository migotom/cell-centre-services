#!/bin/bash
# Import from fixtures

mongoimport --db cell-centre --collection roles --file /fixtures/00_roles.json
mongoimport --db cell-centre --collection employees --file /fixtures/01_employees.json