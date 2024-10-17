import json
import logging
import subprocess
import sys

log = logging.getLogger(__name__)


def fix_command(command):
    if sys.platform == "win32":
        command = command.replace("build/rcc", ".\\build\\rcc.exe", 1)
    return command


def get_cwd():
    import os

    cwd = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    detail = (
        "(rcc doesn't seem to be built, please run `inv local` before running tests)"
    )
    assert "build" in os.listdir(cwd), f"Missing build directory in: {cwd!r} {detail}"

    build_dir = os.path.join(cwd, "build")
    if sys.platform == "win32":
        assert "rcc.exe" in os.listdir(
            build_dir
        ), f"Missing rcc.exe in: {build_dir!r} {detail}"
    else:
        assert "rcc" in os.listdir(build_dir), f"Missing rcc in: {build_dir!r} {detail}"
    return cwd


def log_command(command: str, cwd: str):
    msg = f"Running command: {command!r} cwd: {cwd!r}"
    log.info(msg)


def capture_flat_output(command):
    command = fix_command(command)
    cwd = get_cwd()
    log_command(command, cwd)

    task = subprocess.Popen(
        command,
        shell=True,
        stderr=subprocess.PIPE,
        stdout=subprocess.PIPE,
        cwd=cwd,
    )
    out, _ = task.communicate()
    assert (
        task.returncode == 0
    ), f"Unexpected exit code {task.returncode} from {command!r}"
    return out.decode().strip()


def run_and_return_code_output_error(command):
    command = fix_command(command)
    cwd = get_cwd()
    log_command(command, cwd)

    task = subprocess.Popen(
        command,
        shell=True,
        stderr=subprocess.PIPE,
        stdout=subprocess.PIPE,
        cwd=cwd,
    )
    out, err = task.communicate()
    return task.returncode, out.decode(), err.decode()


def parse_json(content):
    parsed = json.loads(content)
    assert isinstance(parsed, (list, dict)), f"Expecting list or dict; got {parsed!r}"
    return parsed
