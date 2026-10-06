from pathlib import Path
import os,json,subprocess,hashlib
b=Path('/tmp/mudarro-parser-checkpoint/mudarro');root=Path('/tmp/mudarro-mud038-parser/runtime/valid');env=os.environ.copy();env.update(PATH=str(b.parent)+':/home/danielsouza/sdk/go1.27.1/bin:'+env['PATH'],GOCACHE='/tmp/mudarro-go-cache',GOMODCACHE='/tmp/mudarro-mud038-parser/module-cache',GOPROXY='off',GOTOOLCHAIN='local',GOFLAGS='-buildvcs=false');env.pop('GOWORK',None);checks=[]
for args in [['scan','--json'],['run','app-space-go:probe'],['run','app-space-go:test']]:
 p=subprocess.run([str(b),*args,'--root',str(root)],cwd=root,env=env,text=True,capture_output=True,timeout=60);checks.append({'argv':args,'exit':p.returncode,'stdout':p.stdout,'stderr':p.stderr});assert p.returncode==0,(args,p.stderr)
Path('/tmp/mudarro-parser-checkpoint/go-recheck.json').write_text(json.dumps({'binary_sha256':hashlib.sha256(b.read_bytes()).hexdigest(),'checks':checks},indent=2)+'\n');print('Final binary real Go scan, cross-module import, test PASS')
