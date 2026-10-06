import pty,os,subprocess,time,select,json,pathlib,fcntl,termios,struct
b=pathlib.Path('/tmp/mudarro-mud038-parser/runtime');e=dict(os.environ);e.update(PATH='/tmp/mudarro-mud038-parser:/home/danielsouza/sdk/go1.27.1/bin:'+e['PATH'],GOCACHE=str(b/'go-cache'),GOMODCACHE='/tmp/mudarro-mud038-parser/module-cache',GOTOOLCHAIN='local',GOPROXY='off',GOFLAGS='-buildvcs=false',TERM='xterm-256color',NO_COLOR='1');e.pop('GOWORK',None)
m,s=pty.openpty();fcntl.ioctl(s,termios.TIOCSWINSZ,struct.pack('HHHH',32,100,0,0));original=termios.tcgetattr(s);p=subprocess.Popen(['/tmp/mudarro-mud038-parser/mudarro','menu','--root',str(b/'valid')],stdin=s,stdout=s,stderr=s,env=e);start=time.monotonic();events=[];transcript=b''
def pump(seconds):
 global transcript
 end=time.monotonic()+seconds
 while time.monotonic()<end:
  if select.select([m],[],[],0.05)[0]:
   try:chunk=os.read(m,65536)
   except OSError:break
   transcript+=chunk;events.append([round(time.monotonic()-start,4),'o',chunk.decode('utf-8','replace')])
try:
 pump(.5)
 for data,delay in [(b'1\r',.4),(b'4\r',.4),(b'1\r',1),(b'\r',.4),(b'q',.25),(b'q',.25),(b'q',.5)]:os.write(m,data);pump(delay)
 p.wait(timeout=5);pump(.1);assert p.returncode==0;assert b'REAL_WORKSPACE_IMPORT_OK' in transcript;assert termios.tcgetattr(s)==original
finally:
 if p.poll() is None:p.terminate();p.wait(timeout=5)
 os.close(m);os.close(s)
with (b/'workspace-menu.cast').open('w') as f:
 f.write(json.dumps({'version':2,'width':100,'height':32,'title':'Mudarro real Go workspace menu: local cross-module import'})+'\n')
 for event in events:f.write(json.dumps(event)+'\n')
(b/'workspace-menu.pty.txt').write_bytes(transcript);print('real PTY PASS',len(events),'frames, raw lifecycle restored')
