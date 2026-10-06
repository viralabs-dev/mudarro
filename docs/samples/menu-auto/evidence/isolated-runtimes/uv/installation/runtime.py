from pathlib import Path
import subprocess,json,os,hashlib,time,pty,select,signal
base=Path('/tmp/mudarro-runtimes.K7F4J1');root=base/'fixtures/uv-final';root.mkdir(parents=True,exist_ok=True);logs=base/'logs/uv';bin=Path('/tmp/mudarro-timeout-checkpoint/mudarro');uv=base/'uv/tool-venv/bin/uv';python='/home/linuxbrew/.linuxbrew/bin/python3'
env=os.environ.copy();env.update(PATH=str(uv.parent)+':'+str(bin.parent)+':'+env['PATH'],UV_CACHE_DIR=str(base/'cache/uv/runtime'),UV_PYTHON=python,UV_PYTHON_DOWNLOADS='never',UV_OFFLINE='true',UV_NO_CONFIG='true',XDG_CONFIG_HOME=str(base/'uv/config'),XDG_CACHE_HOME=str(base/'cache/uv/xdg'))
(root/'pyproject.toml').write_text('[project]\nname="mudarro-uv-runtime"\nversion="0.1.0"\nrequires-python=">=3.10"\ndependencies=[]\n[tool.uv]\npackage=false\n')
(root/'app.py').write_text('import os,time,signal\nprint("UV startup",os.getpid(),flush=True)\ndef stop(sig,frame):\n print("UV shutdown",sig,flush=True);raise SystemExit(0)\nsignal.signal(signal.SIGTERM,stop)\nwhile True:time.sleep(.1)\n')
(root/'check.py').write_text('import sys\nprint("UV real check PASS",sys.version.split()[0])\n')
(root/'seed.py').write_text('import sqlite3\nc=sqlite3.connect("seed.sqlite3")\nc.execute("CREATE TABLE IF NOT EXISTS sample(id INTEGER PRIMARY KEY, name TEXT NOT NULL)")\nc.execute("INSERT OR IGNORE INTO sample VALUES(1,\'mudarro\')")\nc.commit();assert c.execute("SELECT COUNT(*) FROM sample").fetchone()[0]==1\nprint("UV SQLite custom seed idempotent PASS")\n')
(root/'fail.py').write_text('print("expected controlled failure")\nraise SystemExit(7)\n')
steps=[];out=(logs/'e2e.txt').open('w')
def run(args,expected=0,cwd=root):
 out.write('\nCOMMAND '+repr([str(a) for a in args])+'\n');out.flush();p=subprocess.run([str(a) for a in args],cwd=cwd,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True,timeout=30);out.write(p.stdout);out.write('EXIT '+str(p.returncode)+'\n');out.flush();steps.append({'argv':[str(a) for a in args],'exit':p.returncode,'expected':expected});assert p.returncode==expected,(args,p.returncode,p.stdout);return p.stdout
run([uv,'--version']);run([uv,'lock','--offline']);run([bin,'scan','--root',root,'--json']);scan=json.loads(run([bin,'scan','--root',root,'--json']));(logs/'scan.json').write_text(json.dumps(scan,indent=2)+'\n');assert scan['config']['services'][0]['manager']=='uv';run([bin,'init','--root',root]);before=(root/'mudarro.yaml').read_bytes();run([bin,'init','--root',root],1);assert before==(root/'mudarro.yaml').read_bytes();run([bin,'generate','--root',root,'--dry-run']);assert not (root/'menu.sh').exists()
# Explicit start/custom scripts satisfy Python's documented non-inference contract.
config=scan['config'];config['name']='uv-runtime-explicit';s=config['services'][0];s.pop('pending',None)
for action,file,group in [('start','app.py','aplicacao'),('test','check.py','qualidade'),('lint','check.py','qualidade'),('build','check.py','aplicacao'),('seed','seed.py','banco'),('fail','fail.py','qualidade')]:s['commands'][action]={'args':['uv','run','--offline','python',file],'group':group}
(root/'mudarro.yaml').unlink();(root/'mudarro.json').write_text(json.dumps(config,indent=2)+'\n');run([bin,'generate','--root',root])
def hashes():return {str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in root.rglob('*') if p.is_file() and (p.name=='menu.sh' or '.mudarro/scripts/' in str(p) or p.name=='generated.json')}
h=hashes();run([bin,'generate','--root',root]);assert h==hashes();(logs/'generated-hashes.json').write_text(json.dumps(h,indent=2)+'\n');run([bin,'doctor','--root',root]);run(['bash',root/'.mudarro/scripts/app-python/install.sh']);run(['bash',root/'.mudarro/scripts/app-python/install.sh'])
for action in ['test','lint','build','seed','seed']:run(['bash',root/f'.mudarro/scripts/app-python/{action}.sh'])
run(['bash',root/'.mudarro/scripts/app-python/fail.sh'],1);run([bin,'run','app-python:unknown','--root',root],1)
try:
 run([bin,'run','app-python:up','--root',root]);run([bin,'run','app-python:up','--root',root]);time.sleep(.3);assert 'running' in run([bin,'run','app-python:status','--root',root]);assert 'UV startup' in run([bin,'run','app-python:logs','--root',root]);run([bin,'run','app-python:restart','--root',root]);time.sleep(.3);run([bin,'run','app-python:logs','--root',root])
finally:run([bin,'run','app-python:down','--root',root])
assert 'stopped' in run([bin,'run','app-python:status','--root',root])
# Genuine menu PTY, only fixture application output.
pid,fd=pty.fork()
if pid==0:
 os.chdir(root);os.execve('/bin/bash',['bash',str(root/'menu.sh')],env)
frames=[];start=time.monotonic();sent=False
try:
 while time.monotonic()-start<8:
  ready,_,_=select.select([fd],[],[],.1)
  if ready:
   try:data=os.read(fd,65536)
   except OSError:break
   if not data:break
   frames.append([round(time.monotonic()-start,6),'o',data.decode(errors='replace')])
  if not sent and time.monotonic()-start>.7:os.write(fd,b'0\n');sent=True
finally:
 os.close(fd);_,status=os.waitpid(pid,0)
assert os.waitstatus_to_exitcode(status)==0
with (logs/'menu.cast').open('w') as f:
 f.write(json.dumps({'version':2,'width':80,'height':24,'command':'bash fixture/menu.sh','env':{'TERM':env.get('TERM','xterm')}})+'\n')
 for frame in frames:f.write(json.dumps(frame)+'\n')
(root/'menu.sh').write_text((root/'menu.sh').read_text()+'\n# manual edit preservation probe\n');bh=hashes();run([bin,'generate','--root',root],1);assert bh==hashes()
summary={'manager':'uv','version':'0.12.23','root':str(root),'steps':steps,'runtime':'PASS','lock':'generated by real uv offline','generated_idempotence':True,'custom_sqlite_seed':'two runs, one row','menu_cast':'genuine headless PTY; no desktop screenshot','start':'explicit configuration, not inferred Python entrypoint','cleanup':'down/status stopped','limits':['zero dependencies; no ORM/framework integration','native macOS not executed','failed command tested synchronously; arbitrary descendant cleanup not claimed']};(logs/'summary.json').write_text(json.dumps(summary,indent=2)+'\n');out.close();print(json.dumps(summary,indent=2))
