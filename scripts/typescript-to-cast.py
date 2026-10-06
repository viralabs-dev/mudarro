"""Convert the genuine util-linux classic timing + typescript to asciicast v2."""
from pathlib import Path
import sys
import json, codecs
p=Path(sys.argv[1]) if len(sys.argv)>1 else Path(__file__).resolve().parents[1]/"docs/samples/menu-auto/evidence"
b=(p/'visual-session.txt').read_bytes().split(b'\n',1)[1]
timing=[line.split() for line in (p/'visual-session.timing').read_text().splitlines()]
pos=0; elapsed=0.; events=[]; decoder=codecs.getincrementaldecoder('utf-8')()
for delay,length in timing:
 elapsed+=float(delay); n=int(length); chunk=b[pos:pos+n]; pos+=n
 text=decoder.decode(chunk)
 if text: events.append([round(elapsed,6),'o',text])
assert len(b)>=pos and b[pos:].startswith(b'\nScript done'), repr(b[pos:pos+50])
header={'version':2,'width':100,'height':32,'title':'Mudarro: execução real de menu automático','env':{'TERM':'xterm-256color'}}
(p/'visual-session.cast').write_text('\n'.join(json.dumps(x,ensure_ascii=False) for x in [header,*events])+'\n')
print('Events:',len(events),'duration:',elapsed,'recorded bytes:',pos)
