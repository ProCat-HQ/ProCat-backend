# ProCat-backend
Server for **"ProCat"**.

Create ```.env``` file in root directory and fill following values:
```
BIND_ADDR=
DB_USERNAME=
DB_HOST=
DB_PORT=
DB_NAME=
DB_PASSWORD=
DB_SSLMODE=
CAP_SOLVER_API_KEY=
PASSWORD_SALT=
SIGNING_KEY=
API_KEY_2GIS=
```

Migrations must be applied **inside** server docker container.
Check usage in ```Makefile```