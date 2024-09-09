import psycopg2
from structures.BuildInfo import BuildInfo


class DbPuller:
    def __init__(self):
        # Connect to the PostgreSQL database
        self.conn = psycopg2.connect(
            host="localhost", database="bee", user="postgres", password="secret"
        )

    def pull_from_db(self) -> BuildInfo:

        cursor = self.conn.cursor()

        # Execute the SELECT statement to pull the first row that matches the criteria
        cursor.execute(
            """
                SELECT *
                FROM bee_schema.builds
                WHERE STATUS = 'queued'
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
                    SET STATUS = 'in_progress'
                    WHERE id = %s
                """,
                (build_id,),
            )
            self.conn.commit()

            # Process the fetched row
            # Convert the row to a BuildInfo object
            build_info = BuildInfo(*row)
            # get repository and owner by repo_id
            cursor.execute(
                """
                    SELECT name, user_id
                    FROM bee_schema.repos
                    WHERE id = %s
                """,
                (build_info.repo_id,),
            )
            repo_name, owner_id  = cursor.fetchone()
            cursor.execute(
                """
                    SELECT username
                    FROM bee_schema.users
                    WHERE id = %s
                """,
                (owner_id,),
            )
            owner_name = cursor.fetchone()
            if owner_name and repo_name:
                build_info.owner_name = owner_name
                build_info.repo_name = repo_name
            print("Got:", build_info)
            # Example: Print the build_id of the fetched row
            print("Build ID:", build_info.id)
            # Do something with the row
            return build_info

        # Close the cursor
        cursor.close()
        return None
