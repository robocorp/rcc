import os
import shutil
import sys

from invoke import task

# Determine OS-specific commands
if sys.platform == "win32":
    PYTHON = "python"
    LS = "dir"
    WHICH = "where"
else:
    PYTHON = "python3"
    LS = "ls -l"
    WHICH = "which -a"


@task
def what(c):
    """Show latest HEAD with stats"""
    c.run("go version")
    c.run("git --no-pager log -2 --stat HEAD")


@task
def tooling(c):
    """Display tooling information"""
    print(f"PATH is {os.environ['PATH']}")
    print(f"GOPATH is {os.environ.get('GOPATH', 'Not set')}")
    print(f"GOROOT is {os.environ.get('GOROOT', 'Not set')}")
    print("git info:")
    c.run(f"{WHICH} git || echo NA")


@task
def noassets(c):
    """Remove asset files"""
    import glob

    patterns = [
        "blobs/assets/micromamba.*",
        "blobs/assets/*.zip",
        "blobs/assets/*.yaml",
        "blobs/assets/*.py",
        "blobs/assets/man/*.txt",
        "blobs/docs/*.md",
    ]
    for pattern in patterns:
        for file_path in glob.glob(pattern, recursive=True):
            try:
                os.remove(file_path)
                print(f"Removed: {file_path}")
            except OSError as e:
                print(f"Error removing {file_path}: {e}")


def download_link(version, platform, filename):
    return f"https://downloads.robocorp.com/micromamba/{version}/{platform}/{filename}"


@task
def micromamba(c):
    """Download micromamba files"""
    with open("assets/micromamba_version.txt", "r", encoding="utf-8") as f:
        version = f.read().strip()
    print(f"Using micromamba version {version}")

    platforms = {
        "macos64": "darwin_amd64",
        "windows64": "windows_amd64",
        "linux64": "linux_amd64",
    }

    for platform, arch in platforms.items():
        filename = "micromamba.exe" if platform == "windows64" else "micromamba"
        url = download_link(version, platform, filename)
        output = f"blobs/assets/micromamba.{arch}"
        if os.path.exists(output + ".gz"):
            print(f"Asset {output}.gz already exists, skipping")
            continue
        print(f"Downloading {url} to {output}")
        c.run(f"curl -o {output} {url}")
        print(f"Compressing {output}")
        c.run(f"gzip -f -9 {output}")


@task(pre=[micromamba])
def assets(c):
    """Prepare asset files"""
    import glob
    from zipfile import ZIP_DEFLATED, ZipFile

    # Process template directories
    for directory in glob.glob("templates/*/"):
        basename = os.path.basename(os.path.dirname(directory))
        assetname = os.path.abspath(f"blobs/assets/{basename}.zip")

        if os.path.exists(assetname):
            print(f"Asset {assetname} already exists, skipping")
            continue

        print(f"Directory {directory} => {assetname}")

        with ZipFile(assetname, "w", ZIP_DEFLATED) as zipf:
            for root, _, files in os.walk(directory):
                for file in files:
                    file_path = os.path.join(root, file)
                    arcname = os.path.relpath(file_path, directory)
                    zipf.write(file_path, arcname)

    # Copy asset files
    asset_patterns = ["assets/*.txt", "assets/*.yaml", "assets/*.py"]
    for pattern in asset_patterns:
        for file in glob.glob(pattern):
            print(f"Copying {file} to blobs/assets/")
            shutil.copy(file, "blobs/assets/")

    # Copy man pages
    os.makedirs("blobs/assets/man", exist_ok=True)
    for file in glob.glob("assets/man/*.txt"):
        print(f"Copying {file} to blobs/assets/man/")
        shutil.copy(file, "blobs/assets/man/")

    # Copy docs
    os.makedirs("blobs/docs", exist_ok=True)
    for file in glob.glob("docs/*.md"):
        print(f"Copying {file} to blobs/docs/")
        shutil.copy(file, "blobs/docs/")


@task(pre=[noassets])
def clean(c):
    """Remove build directory"""
    shutil.rmtree("build", ignore_errors=True)
    print("Removed build directory")


@task
def toc(c):
    """Update table of contents on docs/ directory"""
    c.run(f"{PYTHON} scripts/toc.py")
    print("Ran scripts/toc.py")


@task(pre=[toc])
def support(c):
    """Create necessary directories"""
    for dir in ["tmp", "build/linux64", "build/macos64", "build/windows64"]:
        os.makedirs(dir, exist_ok=True)


@task(pre=[support, assets])
def test(c, cover=False):
    """Run tests"""
    os.environ["GOARCH"] = "amd64"
    if cover:
        c.run("go test -cover -coverprofile=tmp/cover.out ./...")
        c.run("go tool cover -func=tmp/cover.out")
    else:
        c.run("go test ./...")


def version() -> str:
    import re

    with open("common/version.go", "r") as file:
        content = file.read()
        match = re.search(r"Version\s*=\s*`v([^`]+)`", content)
        if match:
            return match.group(1)
        else:
            raise ValueError("Version not found in common/version.go")


@task
def version_txt(c):
    """Create version.txt file"""
    support(c)
    target = "build/version.txt"
    v = version()
    with open(target, "w") as f:
        f.write(f"v{v}")
    print(f"Created {target} with version {v}")


@task(pre=[support, version_txt, assets])
def build(c, platform="all"):
    """Build executables"""
    from pathlib import Path

    os.environ["CGO_ENABLED"] = "0"
    os.environ["GOARCH"] = "amd64"

    build_platforms = ["linux", "darwin", "windows"]

    if platform == "all":
        platforms = build_platforms
    else:
        assert platform in build_platforms, f"Invalid platform: {platform}"
        platforms = [platform]

    for goos in platforms:
        os.environ["GOOS"] = goos
        output = f"build/{goos}64/"

        c.run(f"go build -ldflags -s -o {output} ./cmd/...")

        ext = ".exe" if goos == "windows" else ""
        f = f"{output}rcc{ext}"
        assert Path(f).exists(), f"File {f} does not exist"
        print(f"Built: {f}")


@task
def windows64(c):
    """Build windows64 executable"""
    build(c, platform="windows")


@task
def linux64(c):
    """Build linux64 executable"""
    build(c, platform="linux")


@task
def macos64(c):
    """Build macos64 executable"""
    build(c, platform="darwin")


@task
def robotsetup(c):
    """Setup build environment"""
    if not os.path.exists("robot_requirements.txt"):
        raise RuntimeError(
            f"robot_requirements.txt not found. Current directory: {os.path.abspath(os.getcwd())}"
        )
    c.run(f"{PYTHON} -m pip install --upgrade -r robot_requirements.txt")
    c.run(f"{PYTHON} -m pip freeze")


@task
def local(c, do_test=True):
    """Build local, operating system specific rcc"""
    tooling(c)
    if do_test:
        test(c)
    c.run("go build -o build/ ./cmd/...")


@task(pre=[robotsetup, assets, local])
def robot(c):
    """Run robot tests on local application"""
    print("Running robot tests...")
    c.run(f"{PYTHON} -m robot -L DEBUG -d tmp/output robot_tests")
