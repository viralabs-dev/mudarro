#!/usr/bin/env python3
"""Real CLI PTY regression/capture. No desktop or browser. Only owned fixtures."""
import codecs,errno,fcntl,hashlib,json,os,pathlib,pty,re,select,struct,subprocess,sys,termios,time,unicodedata
BIN=pathlib.Path(sys.argv[1]).resolve(); BASE=pathlib.Path(sys.argv[2]).resolve(); BASE.mkdir(parents=True,exist_ok=True)
class Screen:
 def __init__(self,w,h):self.w=w;self.h=h;self.cells=[[' ']*w for _ in range(h)];self.row=0;self.col=0;self.pending='';self.bad=[]
 def feed(self,text):
  self.pending+=text
  while self.pending:
   s=self.pending
   if s[0]=='\x1b':
    if len(s)<2:return
    if s[1]=='[':
     m=re.match(r'\x1b\[([0-9;?]*)([ -/]*)([@-~])',s)
     if not m:return
     arg,_,key=m.groups();self.pending=s[m.end():]; nums=[int(x or '0') for x in arg.lstrip('?').split(';')]
     if key in 'Hf':self.row=(nums[0] or 1)-1;self.col=(nums[1] if len(nums)>1 else 1)-1
     elif key=='K':
      if 0<=self.row<self.h:
       if nums[0]==2:self.cells[self.row]=[' ']*self.w
       else:self.cells[self.row][self.col:]=[' ']*(self.w-self.col)
     elif key=='J' and nums[0]==2:self.cells=[[' ']*self.w for _ in range(self.h)]
     continue
    self.pending=s[2:];continue
   self.pending=s[1:];ch=s[0]
   if ch=='\r':self.col=0;continue
   if ch=='\n':self.row+=1;continue
   if ord(ch)<32:continue
   width=0 if unicodedata.combining(ch) else (2 if unicodedata.east_asian_width(ch) in 'WF' or ord(ch)>=0x1f000 else 1)
   if 0<=self.row<self.h and 0<=self.col<self.w:self.cells[self.row][self.col]=ch
   else:self.bad.append([self.row,self.col,ch])
   self.col+=width
 def text(self):return [''.join(x) for x in self.cells]
def record(profile,locale,theme):
 root=BASE/profile;root.mkdir(exist_ok=True)
 (root/'long.sh').write_text('''#!/usr/bin/env bash
set -euo pipefail
for i in $(seq 1 90); do printf 'Real central output line %03d\\n' "$i"; sleep .01; done
printf '\\033[2J\\033[H\\033]0;fixture-title\\007CONTROL_SAFE\\n'
printf 'owned stderr\\n' >&2
printf 'output-finished\\n'
''')
 (root/'cancel.sh').write_text("#!/usr/bin/env bash\nprintf 'cancel-ready\\n'\ntrap '' TERM\nsleep 30 &\nprintf '%s\\n' $! > owned-child.pid\nwait\n")
 (root/'input.sh').write_text('read answer; printf \"input:%s\\n\" \"$answer\"\n')
 (root/'confirm.sh').write_text('read answer; printf \"confirmed:%s\\n\" \"$answer\"\n')
 commands={'a-output':{'args':['bash','long.sh']},'b-error':{'args':['sh','-c','printf "real-failure\\n"; exit 7']},'c-input':{'args':['sh','input.sh']},'d-cancel':{'args':['bash','cancel.sh']},'e-confirm':{'args':['sh','confirm.sh'],'destructive':True}}
 for c in commands.values():c['group']='quality'
 conf={'version':1,'name':'Aurora','ui':{'locale':locale,'theme':theme,'preview':{'enabled':True,'mouse':'on'}},'services':[{'id':'aurora','dir':'.','language':'custom','infrastructure':{'kind':'custom'},'commands':commands}]}
 (root/'mudarro.json').write_text(json.dumps(conf,indent=2)+'\n')
 master,slave=pty.openpty(); original=termios.tcgetattr(slave); w,h=100,32;fcntl.ioctl(slave,termios.TIOCSWINSZ,struct.pack('HHHH',h,w,0,0))
 env=os.environ.copy();env.pop('NO_COLOR',None);env.update(TERM='xterm-256color')
 p=subprocess.Popen([str(BIN),'menu','--root',str(root)],stdin=slave,stdout=slave,stderr=slave,env=env)
 start=time.monotonic();raw=bytearray();events=[];dec=codecs.getincrementaldecoder('utf8')('replace');screen=Screen(w,h);chrome=None;checks=0
 cast=(BASE/(profile+'.cast')).open('w');cast.write(json.dumps({'version':2,'width':w,'height':h,'timestamp':int(time.time()),'title':'Mudarro persistent shell — real PTY '+profile})+'\n')
 def pump(t=.4,check=True):
  nonlocal chrome,checks
  end=time.monotonic()+t
  while time.monotonic()<end:
   ready,_,_=select.select([master],[],[],min(.02,max(0,end-time.monotonic())))
   if not ready:continue
   try:data=os.read(master,65536)
   except OSError as e:
    if e.errno==errno.EIO:break
    raise
   if not data:break
   raw.extend(data);txt=dec.decode(data);screen.feed(txt);cast.write(json.dumps([round(time.monotonic()-start,6),'o',txt],ensure_ascii=False)+'\n');cast.flush()
  if check:
   rows=screen.text(); current=rows[:10]+rows[-2:]
   if chrome is None:chrome=current
   elif current!=chrome:raise AssertionError('Header/footer changed outside resize: '+repr(current))
   if screen.bad:raise AssertionError('paint outside screen '+repr(screen.bad[:4]))
   checks+=1
 def key(label,data,t=.4,check=True):events.append({'time':round(time.monotonic()-start,6),'label':label,'hex':data.hex()});os.write(master,data);pump(t,check)
 try:
  pump(.5);key('service',b'1\r');key('group',b'1\r');key('preview',b'\x1b[B');key('preview page',b'\x1b[6~');key('wheel in preview',b'\x1b[<65;98;20M');key('outside header mouse',b'\x1b[<65;98;2M');key('collapse',b'\x1b[A')
  key('execute long/ANSI/stdout+stderr',b'1\r',1.7)
  if b'output-finished' not in raw:raise AssertionError('long action missing')
  key('return',b'\r');key('execute error',b'2\r');key('return error',b'\r');key('execute stdin',b'3\r');key('owned stdin',b'fixture-input\r')
  if b'input:fixture-input' not in raw:raise AssertionError('stdin missing')
  key('return stdin',b'\r');key('execute cancel',b'4\r');key('Ctrl+C cancel process group',b'\x03',.7)
  pid=int((root/'owned-child.pid').read_text());
  try:
   state=pathlib.Path('/proc')/str(pid)/'stat'
   if state.exists() and state.read_text().split(') ')[1][0]!='Z':raise AssertionError('owned child survives')
  except FileNotFoundError:pass
  key('execute after cancel',b'2\r');key('return repeat',b'\r');key('destructive prompt',b'5\r');key('explicit confirmation',b'APAGAR aurora\rqueued-child-input\r')
  if b'confirmed:queued-child-input' not in raw:raise AssertionError('confirmation missing')
  key('return confirmation',b'\r');key('actions back',b'q');key('groups back',b'q');key('exit',b'q',.2,False);p.wait(timeout=3)
  restored=termios.tcgetattr(slave)==original
  if p.returncode or not restored:raise AssertionError('exit/termios failure')
  (BASE/(profile+'.pty.txt')).write_bytes(raw)
  meta={'profile':profile,'locale':locale,'theme':theme,'exit_code':p.returncode,'termios_restored':restored,'chrome_checks':checks,'duration_seconds':round(time.monotonic()-start,3),'binary_sha256':hashlib.sha256(BIN.read_bytes()).hexdigest(),'events':events,'capture':'genuine CLI PTY output; injected keys/SGR mouse bytes; no GUI or browser','cases':['preview','mouse bounds','long output','ANSI containment','stderr','error','stdin','cancel descendant','repeat after cancel','destructive confirmation','cleanup']}
  (BASE/(profile+'.metadata.json')).write_text(json.dumps(meta,indent=2)+'\n');print(json.dumps(meta|{'events':'see metadata'},ensure_ascii=False))
 finally:
  cast.close()
  if p.poll() is None:p.terminate();p.wait(timeout=3)
  os.close(slave);os.close(master)
record('aurora-shell-dark-en','en','dark');record('aurora-shell-light-ptbr','pt-BR','light')
