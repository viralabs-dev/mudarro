import os, pty, subprocess, time, select, json, struct, fcntl, termios, hashlib, pathlib, shutil
base=pathlib.Path('/tmp/mudarro-mud033-runtime')
master,slave=pty.openpty(); fcntl.ioctl(slave,termios.TIOCSWINSZ,struct.pack('HHHH',24,100,0,0))
env=os.environ.copy(); env['PATH']=str(base)+':/home/danielsouza/.bun/bin:'+env['PATH']; env['TERM']='xterm-256color'
p=subprocess.Popen([str(base/'mudarro'),'menu','--root',str(base/'fixture')],stdin=slave,stdout=slave,stderr=slave,env=env)
os.close(slave); start=time.monotonic(); sent=False
with open(base/'logs/menu.cast','w') as out:
 out.write(json.dumps({'version':2,'width':100,'height':24,'timestamp':int(time.time()),'env':{'TERM':env['TERM']},'title':'MUD033 genuine Bun menu PTY; injected q exits'})+'\n')
 while time.monotonic()-start<5:
  elapsed=time.monotonic()-start
  if elapsed>.3 and not sent: os.write(master,b'q'); sent=True
  ready,_,_=select.select([master],[],[],.1)
  if ready:
   try: data=os.read(master,65536)
   except OSError: break
   if not data: break
   out.write(json.dumps([round(elapsed,6),'o',data.decode('utf-8','replace')],ensure_ascii=False)+'\n')
  if p.poll() is not None and not ready: break
 if p.poll() is None: p.terminate()
 rc=p.wait(timeout=2)
os.close(master)
assert rc==0,rc
shutil.copytree('/home/danielsouza/dev/viralabs/mudarro/internal',base/'source-snapshot/internal',dirs_exist_ok=True)
shutil.copytree('/home/danielsouza/dev/viralabs/mudarro/cmd',base/'source-snapshot/cmd',dirs_exist_ok=True)
for name in ['go.mod','go.sum']: shutil.copy('/home/danielsouza/dev/viralabs/mudarro/'+name,base/'source-snapshot'/name)
exe=base/'mudarro'; (base/'logs/binary.sha256').write_text(hashlib.sha256(exe.read_bytes()).hexdigest()+'  '+str(exe)+'\n')
owned=[]
for path in pathlib.Path('/proc').glob('[0-9]*/cmdline'):
 try: args=path.read_bytes().split(b'\0')
 except (OSError,PermissionError): continue
 if args and (args[0]==str(exe).encode() or (b'bun' in args[0] and b'app.js' in args)):
  try: cwd=os.readlink(path.parent/'cwd')
  except OSError: cwd=''
  if args[0]==str(exe).encode() or cwd==str(base/'fixture'): owned.append({'pid':path.parent.name,'argv':[a.decode(errors='replace') for a in args if a],'cwd':cwd})
(base/'logs/process-cleanup.json').write_text(json.dumps({'remaining_owned_processes':owned},indent=2)+'\n')
assert not owned,owned
print('PTY cast exit0; source snapshot saved; no owned supervisor/Bun app remains')
