import enum

# CREATE TYPE build_status AS ENUM ('queued', 'in_progress', 'completed');
# CREATE TYPE build_conclusion AS ENUM ('canceled', 'failure', 'success', 'timed_out');


class BuildStatus(enum.Enum):
    QUEUED = "queued"
    IN_PROGRESS = "in_progress"
    COMPLETED = "completed"


class BuildConclusion(enum.Enum):
    CANCELED = "canceled"
    FAILURE = "failure"
    SUCCESS = "success"
    TIMED_OUT = "timed_out"


class BuildInfo:
    def __init__(
        self,
        build_id: int,
        repo_id: int,
        commit_sha: str,
        commit_message: str,
        status: BuildStatus,
        conclusion: BuildConclusion,
        created_at: str,
        updated_at: str,
        owner_name: str = None,
        repo_name: str = None,
    ):
        self.build_id = build_id
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
        return f"{self.build_id}, {self.repo_id}, {self.commit_sha}, {self.commit_message}, {self.status}, {self.conclusion}, {self.created_at}, {self.updated_at}, {self.owner_name}, {self.repo_name}"
