build_test = BuildInfo(
    build_id=1,
    repo_id=1,
    commit_sha="0262a10fb0590f29471feed5ecf53b418b5b0d67",
    commit_message="Update and rename .bee-ci.sh to .bee-ci.json ",
    status="queued",
    conclusion=None,
    created_at="2022-01-01T00:00:00Z",
    updated_at="2022-01-01T00:00:00Z",
    owner_name="bee-ci-system",
    repo_name="example-using-beeci",
)

config_data_test = """
{
    "image": "alpine",
    "commands": [
        "sleep 1",
        "echo 'Hello, World!'"
    ],
    "timeout": 300
}
"""
