for who in reference rd749ksxw2r7g5yscck7sh4k7n8gymkw rd74wc0e3hx2ss0wjp0t6vc75s8h97y2; do
  rm -rf /tmp/a && cp -a /app /tmp/a && cd /tmp/a
  if [ $who = reference ]; then git -c safe.directory='*' apply /p/solution.patch; else git -c safe.directory='*' apply --include='jedi/*' /p/runs/$who/$who.patch; fi
  cd /w
  for p in probe.py probe2.py probe3.py probe4.py probe5.py probe6.py probe7.py probe8.py probe9.py probe10.py probe11.py proto.py; do
    PYTHONPATH=/tmp/a timeout 240 python $p > /o/old-$who-$p.txt 2>&1
  done
  echo "$who done"
done
