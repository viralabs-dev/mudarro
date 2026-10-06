import pathlib,json,os,subprocess,pty,fcntl,termios,struct,time,select,hashlib
base=pathlib.Path('/tmp/mudarro-runtimes.K7F4J1');root=base/'fixtures/pipenv/final';logs=base/'logs/pipenv';bin=pathlib.Path('/tmp/mudarro-timeout-checkpoint/mudarro');tool=base/'pipenv/tool-venv/bin/pipenv'
env={k:v for k,v in os.environ.items() if not k.startswith(('PIP','VIRTUALENV','PYTHON','UV_','POETRY_')) and k!='VIRTUAL_ENV'};env.update(json.loads((base/'pipenv/runtime-env.json').read_text()));env['PATH']=str(tool.parent)+':'+str(bin.parent)+':/home/linuxbrew/.linuxbrew/bin:'+env['PATH']
p=subprocess.run(['bash',str(root/'menu.sh')],input='1\n3\n1\n0\n0\n0\n',cwd=root,env=env,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=10)
(logs/'menu.txt').write_text(p.stdout);assert p.returncode==0 and 'PIPENV REAL CHECK PASS' in p.stdout,(p.returncode,p.stdout)
master,slave=pty.openpty();fcntl.ioctl(slave,termios.TIOCSWINSZ,struct.pack('HHHH',24,100,0,0));env['TERM']='dumb'
p=subprocess.Popen(['bash',str(root/'menu.sh')],cwd=root,env=env,stdin=slave,stdout=slave,stderr=slave);os.close(slave);start=time.monotonic();raw=[];events=[(.15,b'1\n'),(.3,b'3\n'),(.45,b'1\n'),(1.4,b'\n'),(1.7,b'0\n'),(2,b'0\n'),(2.3,b'0\n')]
with (logs/'check.cast').open('w') as out:
 out.write(json.dumps({'version':2,'width':100,'height':24,'timestamp':int(time.time()),'title':'MUD034 genuine generated menu/check; actual PTY TERM=dumb fallback; injected navigation','env':{'TERM':'dumb'}})+'\n')
 while time.monotonic()-start<10:
  if events and time.monotonic()-start>=events[0][0]:
   when,data=events.pop(0);os.write(master,data);out.write(json.dumps([round(time.monotonic()-start,6),'i',data.decode()])+'\n')
  ready,_,_=select.select([master],[],[],.1)
  if ready:
   try:data=os.read(master,65536)
   except OSError:break
   if not data:break
   raw.append(data);out.write(json.dumps([round(time.monotonic()-start,6),'o',data.decode('utf-8','replace')])+'\n')
  elif p.poll() is not None:break
 if p.poll() is None:p.terminate()
 assert p.wait(timeout=2)==0
os.close(master);assert b'PIPENV REAL CHECK PASS' in b''.join(raw)
remaining=[]
for path in pathlib.Path('/proc').glob('[0-9]*/cmdline'):
 try:args=path.read_bytes().split(b'\0');cwd=os.readlink(path.parent/'cwd')
 except OSError:continue
 if args and args[0]==str(bin).encode() and str(root).encode() in args:remaining.append({'pid':path.parent.name,'argv':[a.decode(errors='replace') for a in args if a]})
 elif cwd==str(root) and b'app.py' in args:remaining.append({'pid':path.parent.name,'cwd':cwd})
assert not remaining,remaining
summary={'manager':'pipenv','version':'2026.8.0','python':'3.14.7','officialDirectSHA256':'8da889c636a8cab59a435c1366f30d252ecee465372341eb8cd4da43b055ad9d','officialVerifiedWheels':9,'install':'offline --require-hashes; wheel-only PyPIofficial','root':str(root),'mudarroBinarySHA256':hashlib.sha256(bin.read_bytes()).hexdigest(),'runtimeResults':str(logs/'runtime-results.json'),'menu':'genuine generated menu executed check','cast':str(logs/'check.cast'),'passed':['real lock','install --deploy --ignore-pipfile','sync','check','direct fail exit7','wrapper failure maps to CLI1','scan manager pipenv','explicit init script/start selection','second init rejected config preserved','generate2 hashes','doctor','install/check wrappers','custom sqlite seed2x idempotent','up2 idempotent','status/logs/restart/down/stopped','unknown action rejected','dotenv not loaded'],'remainingOwnedProcesses':remaining,'limitations':['zero third-party application dependencies','custom sqlite seed only; no actual Django/Alembic/Prisma/Goose runtime','startup explicitly selected not inferred','full Pipfile schema/hash semantics not validated by scanner','PTY capture not GUI/physical mouse acceptance'],'initialFixtureMistake':'seed script SQL quoting SyntaxError preserved separately; fixed only owned fixture, not production defect'}
(logs/'summary.json').write_text(json.dumps(summary,indent=2)+'\n');print(json.dumps(summary,indent=2))
