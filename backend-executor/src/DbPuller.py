import psycopg2

# Connect to the PostgreSQL database
conn = psycopg2.connect(
    host="localhost",
    database="bee",
    user="postgres",
    password="secret"
)

# Create a cursor object to interact with the database
cursor = conn.cursor()

# Execute the SELECT statement to pull the first row that matches the criteria
cursor.execute("""
    SELECT *
    FROM bee_schema.builds
    WHERE STATUS = 'queued'
    FOR UPDATE SKIP LOCKED
""")

# Fetch the first row from the result set
row = cursor.fetchone()

# Update the status of the fetched row to "ongoing"
if row:
    build_id = row[0]  # Assuming the first column is the primary key
    cursor.execute("""
        UPDATE bee_schema.builds
        SET STATUS = 'in_progress'
        WHERE id = %s
    """, (build_id,))
    conn.commit()

# Close the cursor and the database connection
cursor.close()
conn.close()

# Process the fetched row
if row:
    # Do something with the row
    print(row)
else:
    print("No matching rows found.")