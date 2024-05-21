import psycopg2
import json
from dotenv import load_dotenv
from os import getenv

load_dotenv()
dbname = getenv("DB_NAME")
dbpass = getenv("DB_PASS")
dbuser = getenv("DB_USER")
dbhost = getenv("DB_HOST")

conn = psycopg2.connect(f"dbname={dbname} user={dbuser} host={dbhost} password={dbpass}")

with open("./two-sum.json") as f:
    # load json data from file into python dict
    data = json.load(f)
  
try:
    with conn.cursor() as cursor:

        # extract data fields
        name = data["name"]
        prompt = data["prompt"]
        starter_code = data["starter_code"]
        difficulty = data["difficulty"]
        function_name = data["function_name"]
        
        query = """
            INSERT INTO problem (problem_name, prompt, starter_code, difficulty, function_name)
            VALUES (%s, %s, %s, %s, %s)
        """
        values = (name, prompt, starter_code, difficulty, function_name)
        cursor.execute(query, values)

        cursor.execute("SELECT problem_id FROM problem WHERE problem_name=%s", (name,)) # this must be a tuple, adding a comma converts it to single element tuple
        id = cursor.fetchone()
        print(id)

        for test in data["tests"]:
            sanity = test["sanity"]
            io = json.dumps(test["io"])

            query = """
            INSERT INTO testcase (problem_id, is_sanity_check, io)
            VALUES (%s, %s, %s::jsonb)
            """
            values = (id, sanity, io)
            cursor.execute(query, values)
        

        cursor.execute("SELECT * FROM problem")
        record = cursor.fetchall()
        print(record)

        cursor.execute("SELECT * FROM testcase")
        record = cursor.fetchall()
        print(record)
        conn.commit()

except Exception as e:
    conn.rollback() # this should happen implicitly, but i want to make it explicit that it's happening for myself later
    print("Error:", e)

finally:
    conn.close()
