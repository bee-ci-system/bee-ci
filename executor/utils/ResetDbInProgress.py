import psycopg2

# Connect to the PostgreSQL database
conn = psycopg2.connect(
    host="localhost", database="bee", user="postgres", password="secret"
)

# Create a cursor object to interact with the database
while True:
    cursor = conn.cursor()

    # Execute the SELECT statement to pull the first row that matches the criteria
    cursor.execute(
        """
        SELECT *
        FROM bee_schema.builds
        WHERE STATUS IN ('in_progress', 'completed')
        FOR UPDATE SKIP LOCKED
    """
    )

    # Fetch the first row from the result set
    row = cursor.fetchone()

    # Update the status of the fetched row to "in_progress"
    if row:
        build_id = row[0]  # Assuming the first column is the primary key
        cursor.execute(
            """
            UPDATE bee_schema.builds
            SET STATUS = 'queued', CONCLUSION = NULL
            WHERE id = %s
        """,
            (build_id,),
        )
        conn.commit()

        # Process the fetched row
        print("Modify:", row)
    else:
        break

    # Close the cursor
    cursor.close()
