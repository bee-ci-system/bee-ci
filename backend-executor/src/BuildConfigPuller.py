import requests
from structures.BuildInfo import BuildInfo
import logging


logger = logging.getLogger(__name__)


class BuildConfigPuller:
    @staticmethod
    def pull_config(build_info: BuildInfo):
        repo_url = "https://raw.githubusercontent.com/{}/{}/{}/.bee-ci.json".format(
            build_info.owner_name, build_info.repo_name, build_info.commit_sha
        )
        logger.debug("Repo url: " + repo_url)
        response = requests.get(repo_url)

        if response.status_code == 200:
            config_data = response.text
            logger.debug(config_data)
            return config_data

        logger.error(
            "Failed to pull config from GitHub: code "
            + str(response.status_code)
            + " from "
            + repo_url
        )
        return None


# Example usage
build_info = BuildInfo(
    id=1,
    repo_id=1,
    commit_sha="9e3ce7d48026e8144383eb69e1f5ac87a9bcd1bf",
    commit_message="Update README.md",
    status="queued",
    conclusion=None,
    created_at="2022-01-01T00:00:00Z",
    updated_at="2022-01-01T00:00:00Z",
    owner_name="bee-ci-system",
    repo_name="example-using-beeci",
)
