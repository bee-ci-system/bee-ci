import requests
from structures.BuildInfo import BuildInfo


def pull_script(build_info: BuildInfo):
    repo_url = "https://raw.githubusercontent.com/{}/{}/{}/.bee-ci.sh".format(
        build_info.owner_name, build_info.repo_name, build_info.commit_sha
    )
    print(repo_url)
    response = requests.get(repo_url)

    if response.status_code == 200:
        script_content = response.text
        print(script_content)
        # Do something with the script content
    else:
        print("Failed to pull script from GitHub"+str(response.status_code))


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
pull_script(build_info)