import json, subprocess, sys
d = json.loads(subprocess.run(['olympus','checks','view','kh7d8qtdcfr0k5mt1rvftcvebn8gb2va','--json'],capture_output=True,text=True).stdout)
for c in d['readiness']['criteria']:
    print('R', c['id'], c.get('status'), c.get('verdict',''), c.get('detail',''), 'STALE' if c.get('stale') else '')
for st in d['prechecks'] if isinstance(d['prechecks'], list) else d['prechecks'].get('stages', []):
    print('P', st.get('stageId'), st.get('status'))
    for ch in st.get('checks', []):
        if ch.get('status') != 'pass': print('   ', ch.get('status'), ch.get('name'), str(ch.get('message',''))[:300])
for k, v in d['dynamicChecks'].items():
    if k.startswith('_') or not isinstance(v, dict): continue
    o = v.get('output') or {}
    print('Q', k, v.get('status'), o.get('verdict'), v.get('updatedAt', v.get('completedAt','')))
