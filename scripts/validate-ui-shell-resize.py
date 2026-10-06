import importlib.util,sys,os,pty,fcntl,termios,struct,subprocess,time,select,pathlib,json,signal,re
# Reuse the actual-output screen decoder without invoking its recordings.
script_dir=pathlib.Path(__file__).resolve().parent
binary=pathlib.Path(sys.argv[1]).resolve();output=pathlib.Path(sys.argv[2]).resolve()
src=(script_dir/'record-ui-shell.py').read_text().split("record('mudarro-shell-dark-en'")[0]
sys.argv=['validator',str(binary),str(output)];scope={};exec(src,scope);Screen=scope['Screen']
root=output;root.mkdir(exist_ok=True)
(root/'mudarro.json').write_text(json.dumps({'version':1,'name':'mudarro','ui':{'lettering':'text','preview':{'enabled':True,'mouse':'on'}},'services':[{'id':'mudarro','dir':'.','language':'custom','infrastructure':{'kind':'custom'},'commands':{'test':{'args':['true'],'group':'quality'}}}]}))
results=[]
for no_color in [False,True]:
 m,s=pty.openpty();before=termios.tcgetattr(s)
 fcntl.ioctl(s,termios.TIOCSWINSZ,struct.pack('HHHH',32,100,0,0))
 env=os.environ.copy();env['TERM']='xterm-256color'
 if no_color:env['NO_COLOR']='1'
 else:env.pop('NO_COLOR',None)
 p=subprocess.Popen([str(binary),'menu','--root',str(root)],stdin=s,stdout=s,stderr=s,env=env)
 def read(t=.3):
  out=b'';end=time.monotonic()+t
  while time.monotonic()<end:
   ready,_,_=select.select([m],[],[],.02)
   if ready:out+=os.read(m,65536)
  return out
 initial=read();initial_screen=Screen(100,32);initial_screen.feed(initial.decode());assert any('mudarro' in line for line in initial_screen.text()[:10]),'Config.Name missing from project header'
 assert not re.search(rb'no\s+ai|sem\s+ia',initial,re.IGNORECASE),'Removed UI slogan reappeared'
 os.write(m,b'1\r');read();os.write(m,b'1\r');read();os.write(m,b'\x1b[B');read()
 for cols,rows in [(12,6),(8,2),(100,32)]:
  fcntl.ioctl(s,termios.TIOCSWINSZ,struct.pack('HHHH',rows,cols,0,0));data=read();sc=Screen(cols,rows);sc.feed(data.decode());
  # DECAWM disables wrap; verify row addresses never leave actual dimensions.
  import re
  # Queued output from the preceding dimensions may arrive before the resize reset.
  steady=data[data.rfind(b'\x1b[r'):] if b'\x1b[r' in data else data
  addresses=[int(x) for x in re.findall(rb'\x1b\[(\d+);\d+H',steady)]
  assert addresses and all(1<=x<=rows for x in addresses),(cols,rows,addresses)
  if rows==2:
   os.write(m,b'\r');data=read();assert b'Execution' not in data,'invisible action executed'
  results.append({'no_color':no_color,'cols':cols,'rows':rows,'cursor_rows_within_actual_height':True})
 p.send_signal(signal.SIGTERM);cleanup=read();p.wait(timeout=2)
 assert termios.tcgetattr(s)==before
 assert b'\x1b[?1049l' in cleanup and b'\x1b[r' in cleanup
 if no_color:assert not re.search(rb'\x1b\[[0-9;]+m',initial)
 os.close(s);os.close(m)
print(json.dumps({'resize_checks':results,'SIGTERM_cleanup':True,'NO_COLOR':True},indent=2))
