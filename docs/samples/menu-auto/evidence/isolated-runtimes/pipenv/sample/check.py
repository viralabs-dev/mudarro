import os,sys
assert sys.prefix != sys.base_prefix
assert os.getenv("MUDARRO_ENV_SENTINEL") is None
print("PIPENV REAL CHECK PASS",sys.version.split()[0])
