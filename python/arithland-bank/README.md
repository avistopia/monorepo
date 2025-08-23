# Arithland Bank

## Setup Env

1. CD to the `arithland-bank` directory and copy `.env.template` to `.env`.

## Start the Service

### Windows

1. Open **Command Prompt** (search for `cmd` in the Start menu).

2. `cd` to the python directory.
    ```bash
    cd MONOREPO_ROOT\python
    ```

3. Activate the venv.
    ```bash
    venv\Scripts\activate
    ```

4. CD to the `arithland-bank` directory.

5. **(Optionally)** clear existing data, by deleting the file `db.sqlite3` in the `arithland-bank` directory.

6. Migrate the database. (If it's the first time, or you deleted the data)
    ```bash
    python manage.py migrate
    ```

7. Create an admin. (If it's the first time, or you deleted the data)
    ```bash
    python manage.py createsuperuser
    ```

8. Run the server.
    ```bash
    python manage.py runserver 0.0.0.0:8000
    ```

9. Open the application [0.0.0.0:8000](http://0.0.0.0:8000) and login with admin credentials.

10. Find your machine's local IP, and share with other machines in the same network,
    they can access it using the address `THE_IP:8000`. (You might need to disable the
    firewall on the server machine.)

## Known Issues

### To Do List

- team card background
- add revert transaction functionality
