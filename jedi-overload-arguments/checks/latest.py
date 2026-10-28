import json, subprocess, sys, time
since = float(sys.argv[1]) if len(sys.argv) > 1 else 0
d = json.loads(subprocess.run(['olympus','checks','view','kh7d8qtdcfr0k5mt1rvftcvebn8gb2va','--json'],capture_output=True,text=True).stdout)
latest = {}
for st in d['prechecks']:
    k = st['stageId']
    if k not in latest or float(st['_creationTime']) > float(latest[k]['_creationTime']): latest[k] = st
now = time.time()*1000
for k, st in latest.items():
    age = (now - float(st['_creationTime']))/60000
    print('P %-18s %-8s %5.0f min ago %s' % (k, st['status'], age, 'NEW' if float(st['_creationTime']) > since else 'old'))
    for ch in st.get('checks', []):
        if ch.get('status') != 'pass': print('     ', ch.get('status'), ch.get('name'), '|', str(ch.get('message',''))[:400].replace('\n',' '))
for c in d['readiness']['criteria'][:3]:
    print('R', c['id'], c.get('status'), c.get('verdict',''), c.get('detail',''))
for k, v in d['dynamicChecks'].items():
    if k.startswith('_') or not isinstance(v, dict): continue
    o = v.get('output') or {}
    t = float(v.get('_creationTime', v.get('updatedAt', 0)) or 0)
    print('Q %-22s %-10s %-11s %s' % (k, v.get('status'), o.get('verdict'), 'NEW' if t > since else 'old'))
