import psycopg2
import logging
from structures.BuildInfo import BuildInfo, BuildConclusion


class DbPuller:
    def __init__(self, host: str, port: int, database: str, user: str, password: str):
        # Connect to the PostgreSQL database
        self.conn = psycopg2.connect(
            host=host, port=port, database=database, user=user, password=password
        )
        self.logger = logging.getLogger(__name__)
        self.logger.info("Connected to the database")

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
                    SET STATUS = 'in_progress', updated_at = NOW()
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
            repo_name, owner_id = cursor.fetchone()
            cursor.execute(
                """
                    SELECT username
                    FROM bee_schema.users
                    WHERE id = %s
                """,
                (owner_id,),
            )
            owner_name = cursor.fetchone()[0]
            if owner_name and repo_name:
                build_info.owner_name = owner_name
                build_info.repo_name = repo_name
            self.logger.info("Got: %s", build_info)
            return build_info

        # Close the cursor
        cursor.close()
        return None

    # update build status to finished
    def update_conclusion(self, build_id: int, conclusion: BuildConclusion):
        conclusion_str = conclusion.value
        cursor = self.conn.cursor()
        cursor.execute(
            """
                UPDATE bee_schema.builds
                SET conclusion = %s, status = 'completed', updated_at = NOW()
                WHERE id = (
                    SELECT id
                    FROM   bee_schema.builds
                    WHERE  id = %s
                    FOR    UPDATE SKIP LOCKED
                )
            """,
            (conclusion_str, build_id),
        )
        self.conn.commit()
        cursor.close()
        self.logger.info(
            "Build (id: %d) conclusion updated to %s", build_id, conclusion_str
        )
