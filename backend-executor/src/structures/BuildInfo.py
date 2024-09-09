class BuildInfo:
    def __init__(
        self,
        id,
        repo_id,
        commit_sha,
        commit_message,
        status,
        conclusion,
        created_at,
        updated_at,
        owner_name = None,
        repo_name = None
    ):
        self.id = id
        self.repo_id = repo_id
        self.commit_sha = commit_sha
        self.commit_message = commit_message
        self.status = status
        self.conclusion = conclusion
        self.created_at = created_at
        self.updated_at = updated_at
        self.owner_name = owner_name
        self.repo_name = repo_name

    def __str__(self):
        return f"{self.id}, {self.repo_id}, {self.commit_sha}, {self.commit_message}, {self.status}, {self.conclusion}, {self.created_at}, {self.updated_at}, {self.owner_name}, {self.repo_name}"
