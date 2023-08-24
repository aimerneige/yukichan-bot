import json
import sqlite3

# wget https://github.com/FloatTech/zbpdata/raw/main/Diana/text.db
con = sqlite3.connect("./text.db")
cur = con.cursor()
res = cur.execute("SELECT data FROM text")
data = res.fetchall()
data = [val[0] for val in data]
with open("./text.json", 'w') as f:
  f.write(json.dumps(data, ensure_ascii=False))
