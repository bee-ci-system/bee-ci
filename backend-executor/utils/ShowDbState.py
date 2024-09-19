import psycopg2

# Connect to the PostgreSQL database
conn = psycopg2.connect(
    host="localhost", database="bee", user="postgres", password="secret"
)

# Create a cursor object to interact with the database
cursor = conn.cursor()

# Execute the SELECT statement to pull all rows that match the criteria
cursor.execute(
    """
    SELECT *
    FROM bee_schema.builds
    WHERE STATUS IN ('in_progress', 'completed')
    FOR UPDATE SKIP LOCKED
"""
)

# Fetch all rows from the result set
rows = cursor.fetchall()

# Print each row
for row in rows:
    print(row)

# Close the cursor
cursor.close()

# Close the connection
conn.close()
