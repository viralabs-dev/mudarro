from pathlib import Path
import os,json,subprocess,hashlib
b=Path('/tmp/mudarro-runtimes.K7F4J1');binary=Path('/tmp/mudarro-parser-checkpoint/mudarro');out=Path('/tmp/mudarro-parser-checkpoint/runtime-recheck');out.mkdir(exist_ok=True);results=[]
for tool,sub,action in [('pipenv','pipenv/final','check'),('uv','uv-final','test'),('poetry','poetry','test'),('yarn-classic','yarn/classic','test'),('yarn-modern','yarn/modern','test')]:
 env=os.environ.copy();root=b/'fixtures'/sub;priv={}
 if tool=='pipenv':priv=json.loads((b/'pipenv/runtime-env.json').read_text());prefix=b/'pipenv/tool-venv/bin'
 elif tool=='uv':
  prefix=b/'uv/tool-venv/bin';priv={'UV_CACHE_DIR':str(b/'cache/uv/runtime'),'UV_PYTHON':'/home/linuxbrew/.linuxbrew/bin/python3','UV_PYTHON_DOWNLOADS':'never','UV_OFFLINE':'true','UV_NO_CONFIG':'true'}
 elif tool=='poetry':
  prefix=b/'poetry/tool-venv/bin';priv={'POETRY_CACHE_DIR':str(b/'cache/poetry/runtime'),'POETRY_CONFIG_DIR':str(b/'poetry/config'),'POETRY_DATA_DIR':str(b/'poetry/data'),'POETRY_VIRTUALENVS_IN_PROJECT':'true','POETRY_KEYRING_ENABLED':'false','POETRY_NO_INTERACTION':'1'}
 else:
  name=tool.split('-')[1];priv=json.loads((b/f'logs/yarn/{name}-environment.json').read_text());priv.pop('PATH',None);prefix=b/tool/'bin'
 env.update(priv);env['PATH']=str(binary.parent)+':'+str(prefix)+':'+('/home/danielsouza/.nvm/versions/node/v24.15.0/bin:' if tool=='yarn-modern' else '')+env['PATH'];sid='app-javascript' if tool.startswith('yarn-') else 'app-python';log=[]
 for args in [[str(binary),'scan','--root',str(root),'--json'],[str(binary),'run',sid+':'+action,'--root',str(root)]]:
  p=subprocess.run(args,cwd=root,env=env,capture_output=True,text=True,timeout=30);log.append({'argv':args,'exit':p.returncode,'stdout':p.stdout,'stderr':p.stderr});assert p.returncode==0,(tool,args,p.stderr)
 (out/(tool+'.json')).write_text(json.dumps(log,indent=2)+'\n');results.append({'tool':tool,'scan_exit':0,'quality_exit':0})
(out/'summary.json').write_text(json.dumps({'binary_sha256':hashlib.sha256(binary.read_bytes()).hexdigest(),'checks':results,'scope':'post SemVer/Go fix actual scan and quality on all five already-installed versions; original full runtime evidence retained'},indent=2)+'\n');print('All five current-binary scans and actual quality actions PASS')
