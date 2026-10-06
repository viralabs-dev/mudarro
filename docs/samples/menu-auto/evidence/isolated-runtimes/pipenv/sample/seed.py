import sqlite3
with sqlite3.connect("fixture.sqlite") as db:
 db.execute("create table if not exists seed(id integer primary key,value text)")
 db.execute("insert or ignore into seed values(?,?)",(1,chr(102)))
 assert db.execute("select count(*) from seed").fetchone()[0]==1
print("PIPENV SQLITE CUSTOM SEED PASS")
