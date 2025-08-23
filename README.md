# Avistopia MonoRepo

## Projects

| Project Name       | Language | Directory             |
|--------------------|----------|-----------------------|
| Arithland Bank     | Python   | python/arithland-bank |
| Arithland Telegram | Go       | go/arithland-telegram |

**Note.** Use PyCharm if working on a Python project, and use Goland if working on a Golang project.

## Quick Setup

### Windows

1. Open **Command Prompt** (search for `cmd` in the Start menu).

2. `cd` somewhere convenient.
    ```bash
    cd some\directory
    ```

3. Clone this repository. (If git is not installed, install it, and close and reopen the `cmd`.)
    ```bash
    git clone git@github.com:avistopia/monorepo.git
    cd monorepo
    ```

4. For Python
    1. CD to the Python subdirectory.
        ```bash
        cd python
        ```
    2. Create a virtual environment.
        ```bash
        python -m venv venv
        ```

    3. Activate the Virtual Environment. (After activation, your Command Prompt will show `(venv)` at the beginning
       of the line, indicating that the virtual environment is active.)
        ```bash
        venv\Scripts\activate
        ```

    4. Install requirements.
        ```bash
        pip install -r requirements.txt
        ```
5. For Golang.
    1. CD to the Golang subdirectory, and exec into Golang development container.
        ```bash
        cd go
        make up
        make exec
        ```
    2. Install dependencies.
        ```bash
        make tidy
        ```
    3. Generate autogen files.
        ```bash
        make autogen
        ```
    4. You can delete autogen files and create again if needed.
        ```bash
        make clean-autogen
        make autogen
        ```
    5. You can exit the Golang development container using `exit` command or ctrl+D. When you exit the exec shell, the container will remain running.

    6. You can also stop the Golang development container.
        ```bash
        make down
        ```

6. Check the `README.md` in any project directory you want to set up.

## Contribution

1. _(Only first time)_ setup GitHub SSH key. (Refer to the [GitHub SSH Key Doc](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account?platform=mac))

2. Make your changes. Run test and lint.
    1. If you changed a Python project. (If receiving lint issues, running `make lint-fix` can help fix some of the issues automatically.)
        ```shell
        cd python
        make install-lint # only first time
        make lint
        make test
        make lint-fix     # only if receiving lint issues, running this can help fix some of the issues automatically.
        ```
    2. If you changed a Golang project.
        ```shell
        cd go
        make up
        make exec
        ```
        And inside the container run
        ```shell
        make install-lint # only first time
        make lint
        make test
        make lint-fix     # only if receiving lint issues, running this can help fix some of the issues automatically.
        ```

3. Add and commit your changes. (Refer to
   the [GoLand Git Doc](https://www.jetbrains.com/help/go/commit-and-push-changes.html))

4. Make sure to pull before pushing your commit. Make sure to use rebase mode and not the default merge mode.
   This can be done via the IDE GUI as well as the CLI.
   ```shell
   git pull --rebase
   ```

5. Push your changes.

6. Open [Actions tab](https://github.com/avistopia/monorepo/actions/workflows/deploy.yaml) in the Monorepo GitHub
   repository, and trigger the deployment pipeline for the application you modified.

7. Monitor the pipeline. If it fails, this means either tests or lints failed, or some error happened when building and
   deploying the project. The exact error can be found in the pipeline logs.

8. If the triggered pipeline finishes successfully, your changes will be available after a few minutes on the server.

## Deployment

Only if you want to deploy from local and not using GitHub pipeline, not recommended. Copy the `.env` file into
`production.env` and put production secrets there. Then use each of the following commands for the respective actions.

```
make init-ansible
make ansible-setup-instance
make ansible-deploy app=infra/nginx
make ansible-deploy app=infra/postgres
make ansible-deploy app=python/arithland-bank
make ansible-deploy app=go/arithland-telegram
```
