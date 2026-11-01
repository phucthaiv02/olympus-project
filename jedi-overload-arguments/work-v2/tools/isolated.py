import subprocess, sys
from concurrent.futures import ThreadPoolExecutor
ids = [l.strip() for l in open('../ids.txt') if l.strip()]
def run(args):
    seed, node = args
    r = subprocess.run(['../venv/bin/python', '-m', 'pytest', '-p', 'no:cacheprovider', '-q', node],
                       env={'PYTHONHASHSEED': str(seed), 'PATH': '/usr/bin:/bin'}, capture_output=True, text=True)
    return None if r.returncode == 0 else 'FAIL seed=%s %s' % (seed, node)
with ThreadPoolExecutor(8) as ex:
    res = [x for x in ex.map(run, [(s, n) for s in range(4) for n in ids]) if x]
print('\n'.join(res) or 'ALL PASS', len(ids))
