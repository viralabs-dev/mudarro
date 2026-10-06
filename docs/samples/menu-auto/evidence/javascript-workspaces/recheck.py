from pathlib import Path
import json,os,subprocess,hashlib
root=Path('/tmp/mudarro-mud036-runtime');binary=Path('/tmp/mudarro-js-checkpoint/mudarro');logs=[]
for name in ['npm','pnpm','classic','modern','bun']:
 fixture=root/name;env=os.environ.copy();env.update(json.loads((root/(name+'-environment.json')).read_text()));env['PATH']=str(binary.parent)+':'+env['PATH'];env.update(TMPDIR=str(fixture/'private-tmp'),BUN_INSTALL_CACHE_DIR=str(fixture/'private-bun-cache'),pnpm_config_store_dir=str(fixture/'private-pnpm-store'))
 for args in [['scan','--json'],['run','packages-a-javascript:test']]:
  p=subprocess.run([str(binary),*args,'--root',str(fixture)],cwd=fixture,env=env,capture_output=True,text=True,timeout=60);logs.append({'tool':name,'argv':args,'exit':p.returncode,'stdout':p.stdout,'stderr':p.stderr});assert p.returncode==0,(name,p.stderr)
Path('/tmp/mudarro-js-checkpoint/runtime-recheck.json').write_text(json.dumps({'binary_sha256':hashlib.sha256(binary.read_bytes()).hexdigest(),'checks':logs},indent=2)+'\n');print('Final036binary allfive actual scan+childtest PASS')
