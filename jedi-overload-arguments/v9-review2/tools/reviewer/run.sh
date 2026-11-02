#!/bin/bash
docker run --rm --network none -v /tmp/claude-1000/-home-tes-Desktop-ADK-olympus-project/f7f63570-6352-45a7-966e-f738c710dbdf/scratchpad/reviewer:/w jedi-v7:base bash -c "cd /app && git -c safe.directory='*' apply /w/solution.patch && cd /w && PYTHONPATH=/app:/w python $1"
