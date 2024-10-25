import logging
import requests
from structures.BuildInfo import BuildInfo


logger = logging.getLogger(__name__)


class BuildConfigPuller:
    @staticmethod
    def pull_config(build_info: BuildInfo) -> str:
        repo_url = "https://raw.githubusercontent.com/{}/{}/{}/.bee-ci.json".format(
            build_info.owner_name, build_info.repo_name, build_info.commit_sha
        )
        logger.debug("Repo url: %s", repo_url)
        response = requests.get(repo_url, timeout=10)

        if response.status_code == 200:
            config_data = response.text
            logger.debug(config_data)
            return config_data

        logger.error(
            "Failed to pull config from GitHub: code %s from %s",
            response.status_code,
            repo_url,
        )
        return None
