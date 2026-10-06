# UV/Poetry readiness, 2026-10-06

Both tools absent in checked PATH/common candidate paths; search nonexhaustive. No versions, installations, package resolution or manager runtime execution performed. readiness.json records exact checked paths. UV/Poetry folders contain dependency-free manifests and probe.py; .lock files are explicit detection-only sentinels, not generated real locks.

Future runtime checks require an already-authorized installed executable and exclusive fresh fixture copies. Do not run install or network resolution implicitly. After removing sentinel lock, inspect and use tool-specific offline options; failure due to missing cached interpreter/build backend is a real blocker, not permission to download. A minimal interpreter invocation can test runtime without installing the project:

```bash
# UV: explicit system Python, avoid creating/installing a project environment.
uv --version
uv run --offline --no-project --python /absolute/path/to/existing/python probe.py
# Poetry: isolate environment/cache in disposable fixture, disable keyring.
# poetry run alone may create an environment; authorize that temporary write and
# choose an existing interpreter first. No poetry install assumed.
poetry --version
POETRY_CACHE_DIR=/tmp/mudarro-poetry-cache POETRY_VIRTUALENVS_IN_PROJECT=true POETRY_KEYRING_ENABLED=false poetry env use /absolute/path/to/existing/python
POETRY_CACHE_DIR=/tmp/mudarro-poetry-cache POETRY_VIRTUALENVS_IN_PROJECT=true POETRY_KEYRING_ENABLED=false poetry run python probe.py
```

These are reproduction guidance, not executed or verified against installed tool versions. Mudarro-generated install actions are deliberately not invoked. Full project manager integration requires authentic lock generation, project environment and action lifecycle under separate availability/authorization. Do not describe detection-only scans as runtime success.
