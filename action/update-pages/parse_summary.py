import re

text = open('/tmp/drift-output.txt').read()
m = re.search(
    r'\*\*Total:\s*(\d+)\*\*\s*\|\s*Breaking:\s*(\d+)\s*\|\s*Non-Breaking:\s*(\d+)\s*\|\s*Info:\s*(\d+)',
    text,
)
total, breaking, non_breaking, info = m.groups() if m else ('0', '0', '0', '0')
open('/tmp/summary.txt', 'w').write(f"{total}\n{breaking}\n{non_breaking}\n{info}\n")
