from pathlib import Path
import os,pty,select,time,json,re,struct,fcntl,termios
b=Path('/tmp/mudarro-runtimes.K7F4J1');root=b/'fixtures/uv-final';logs=b/'logs/uv';env=os.environ.copy();env.update(PATH=str(b/'uv/tool-venv/bin')+':/tmp/mudarro-timeout-checkpoint:'+env['PATH'],UV_CACHE_DIR=str(b/'cache/uv/runtime'),UV_PYTHON='/home/linuxbrew/.linuxbrew/bin/python3',UV_PYTHON_DOWNLOADS='never',UV_OFFLINE='true',UV_NO_CONFIG='true',TERM='dumb')
pid,fd=pty.fork()
if pid==0:os.chdir(root);os.execve('/bin/bash',['bash',str(root/'menu.sh')],env)
fcntl.ioctl(fd,termios.TIOCSWINSZ,struct.pack('HHHH',40,100,0,0));frames=[];raw='';events=[];start=time.monotonic()
def pump(seconds=.5):
 global raw
 text='';end=time.monotonic()+seconds
 while time.monotonic()<end:
  if not select.select([fd],[],[],.05)[0]:continue
  try:data=os.read(fd,65536)
  except OSError:break
  if not data:break
  txt=data.decode(errors='replace');raw+=txt;text+=txt;frames.append([round(time.monotonic()-start,6),'o',txt])
 return text
def key(label,value,seconds=.5):events.append({'time':round(time.monotonic()-start,6),'label':label,'input':value});os.write(fd,value.encode());return pump(seconds)
try:
 pump();groups=key('select fixture service','1\n');m=re.search(r'(\d+)\.\s+quality',groups,re.I);assert m,groups
 actions=key('select quality group',m[1]+'\n');m=re.search(r'(\d+)\.\s+test\b',actions,re.I);assert m,actions
 body=key('execute real uv test',m[1]+'\n',1);assert 'UV real check PASS' in body,body
 key('return from execution','\n');key('actions back','0\n');key('groups back','0\n');key('exit','0\n');
 deadline=time.monotonic()+3
 while time.monotonic()<deadline:
  done,status=os.waitpid(pid,os.WNOHANG)
  if done:break
  pump(.1)
 else:
  os.kill(pid,15);os.waitpid(pid,0);raise AssertionError('owned recording child did not exit')
 assert os.waitstatus_to_exitcode(status)==0
finally:os.close(fd)
with (logs/'runtime-menu.cast').open('w') as f:
 f.write(json.dumps({'version':2,'width':100,'height':40,'title':'uv 0.12.23 actual generated menu and test','env':{'TERM':'dumb'}})+'\n')
 for frame in frames:f.write(json.dumps(frame)+'\n')
(logs/'runtime-menu.txt').write_text(raw);(logs/'runtime-menu-events.json').write_text(json.dumps(events,indent=2)+'\n');print('Genuine UV menu test recording PASS')
