import os
import subprocess
import sys

use_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

assert len(sys.argv) >= 2, "No task provided when calling `call_invoke.py`"
task = sys.argv[1]
exit(subprocess.run(("invoke", task), cwd=use_dir).returncode)
