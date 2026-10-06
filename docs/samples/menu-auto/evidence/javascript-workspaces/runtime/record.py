import pty,os,subprocess,time,select,json,pathlib,fcntl,termios,struct
b=pathlib.Path('/tmp/mudarro-mud036-runtime');root=b/'pnpm';e=dict(os.environ);e.update(json.loads((b/'pnpm-environment.json').read_text()));e.update(PATH='/tmp/mudarro-mud036-root:'+e['PATH'],TERM='xterm-256color',NO_COLOR='1',pnpm_config_store_dir=str(root/'private-pnpm-store'))
marker=root/'packages/a/OWN_CWD_MARKER';marker.unlink(missing_ok=True)
m,s=pty.openpty();fcntl.ioctl(s,termios.TIOCSWINSZ,struct.pack('HHHH',32,100,0,0));original=termios.tcgetattr(s);p=subprocess.Popen(['/tmp/mudarro-mud036-root/mudarro','menu','--root',str(root)],stdin=s,stdout=s,stderr=s,env=e);start=time.monotonic();events=[];transcript=b''
def pump(seconds):
 global transcript
 end=time.monotonic()+seconds
 while time.monotonic()<end:
  if select.select([m],[],[],.05)[0]:
   try:chunk=os.read(m,65536)
   except OSError:break
   transcript+=chunk;events.append([round(time.monotonic()-start,4),'o',chunk.decode('utf-8','replace')])
try:
 pump(.5)
 for data,delay in [(b'2\r',.35),(b'3\r',.35),(b'1\r',1),(b'\r',.4),(b'q',.25),(b'q',.25),(b'q',.4)]:os.write(m,data);pump(delay)
 p.wait(timeout=5);pump(.1);assert p.returncode==0;assert b'WORKSPACE_CHILD_a' in transcript;assert marker.read_text()==str(root/'packages/a');assert termios.tcgetattr(s)==original
finally:
 if p.poll() is None:p.terminate();p.wait(timeout=5)
 os.close(m);os.close(s)
with (b/'workspace-child-menu.cast').open('w') as f:
 f.write(json.dumps({'version':2,'width':100,'height':32,'title':'Mudarro pnpm workspace: child script in its own cwd'})+'\n')
 for event in events:f.write(json.dumps(event)+'\n')
(b/'workspace-child-menu.pty.txt').write_bytes(transcript)
(b/'workspace-child-menu-validation.json').write_text(json.dumps({'exit_code':0,'termios_restored':True,'genuine_pty':True,'child_cwd_verified':str(root/'packages/a'),'actual_pnpm_script_output_verified':True,'frames':len(events)},indent=2)+'\n')
print('PTY PASS',len(events))
