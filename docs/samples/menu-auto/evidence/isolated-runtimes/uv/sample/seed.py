import sqlite3
c=sqlite3.connect("seed.sqlite3")
c.execute("CREATE TABLE IF NOT EXISTS sample(id INTEGER PRIMARY KEY, name TEXT NOT NULL)")
c.execute("INSERT OR IGNORE INTO sample VALUES(1,'mudarro')")
c.commit();assert c.execute("SELECT COUNT(*) FROM sample").fetchone()[0]==1
print("UV SQLite custom seed idempotent PASS")
