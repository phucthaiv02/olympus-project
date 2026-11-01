import sys
sys.path.insert(0, sys.argv[1])
exec(open(sys.argv[0].replace('one.py', 'probe_review.py')).read().split('probes = [')[0].split("FILTER =")[0] + "HEADER = " + open(sys.argv[0].replace('one.py','probe_review.py')).read().split('HEADER = ')[1].split('probes = [')[0])
for code in sys.argv[2:]:
    print(repr(code[-40:]), infer(eval(code)))
