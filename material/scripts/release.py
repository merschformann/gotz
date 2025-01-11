#!/usr/bin/env python
import argparse
import os
import subprocess


def get_version():
    """
    Get the current version as defined in the source code.
    """
    with open(os.path.join(os.path.dirname(__file__), "..", "..", "version.go")) as f:
        return f.read().strip().split("\n")[-1].split(" ")[-1].replace('"', "")


def get_version_from_git():
    """
    Get the current version from git.
    """
    return subprocess.check_output("git tag".split()).decode().strip().split("\n")[-1]


def release(tag: str, notes: str = ""):
    """
    Release via github cli tool and show output as it occurs.
    """
    # Prepare call
    call = [
        "gh",
        "release",
        "create",
        tag,
        "--target",
        "main",
        "--title",
        f"Release {tag}",
    ]
    if notes:
        call.extend(["--notes", notes])
    else:
        call.extend(["--generate-notes"])
    # Get user confirmation
    print(" ".join(call))
    if not input("Continue? [y/N] ").lower().startswith("y"):
        print("Aborting.")
        return
    # Release
    with subprocess.Popen(
        call,
        stdout=subprocess.PIPE,
        bufsize=1,
        universal_newlines=True,
    ) as p:
        for line in p.stdout:
            print(line, end="")  # process line here
    if p.returncode != 0:
        raise subprocess.CalledProcessError(p.returncode, p.args)


def main():
    """
    Main function.
    """
    parser = argparse.ArgumentParser(description="Release a new version of gotz.")
    parser.add_argument(
        "--notes",
        "-n",
        type=str,
        help="Notes to add to the release.",
        default="",
    )
    args = parser.parse_args()
    git_version = get_version_from_git()
    go_version = get_version()
    print("Git version: {}".format(git_version))
    print("Go version: {}".format(go_version))
    if git_version == go_version:
        print("Cannot release an already released version.")
        return
    print(f"Releasing new version {go_version} ...")
    release(go_version, args.notes)


if __name__ == "__main__":
    main()
