from pathlib import Path
import os,subprocess,json,hashlib,shutil
base=Path('/tmp/mudarro-mud038-parser/runtime');binary='/tmp/mudarro-mud038-parser/mudarro';env=dict(os.environ);env.update(PATH='/tmp/mudarro-mud038-parser:/home/danielsouza/sdk/go1.27.1/bin:'+env['PATH'],GOCACHE=str(base/'go-cache'),GOMODCACHE='/tmp/mudarro-mud038-parser/module-cache',GOTOOLCHAIN='local',GOPROXY='off',GOFLAGS='-buildvcs=false');env.pop('GOWORK',None);checks=[]
def run(label,args,root,ok=True):
 p=subprocess.run(args,cwd=root,env=env,capture_output=True,text=True,timeout=60);(base/(label+'.log')).write_text(p.stdout+p.stderr);checks.append({'case':label,'exit_code':p.returncode,'expected_success':ok});assert (p.returncode==0)==ok,(label,p.stdout,p.stderr);return p
# Additional bounded/metadata scan-only fixtures.
extras={'nonmodule':'go 1.23\nuse ./plain\n','duplicate':'go 1.23\nuse (./app ./app/.)\n','fifo':None,'oversized':'x'*(4*1024*1024+1)}
for name,work in extras.items():
 r=base/name;r.mkdir(exist_ok=True);(r/'app').mkdir(exist_ok=True);(r/'app/go.mod').write_text('module example.test/independent\ngo 1.23\n');(r/'plain').mkdir(exist_ok=True)
 if work is None:os.mkfifo(r/'go.work')
 else:(r/'go.work').write_text(work)
for name in ['valid','missing','invalid','outside','symlink',*extras]:
 r=base/name;p=run(name+'-after-scan',[binary,'scan','--json','--root',str(r)],r);data=json.loads(p.stdout)
 if name=='valid':
  w=data['workspaces'][0];assert w['members']==['app space','lib'];assert w['modules']=={'app space':'example.test/app','lib':'example.test/lib'};assert w['go_version']=='1.23';assert len(data['config']['services'])==4
 else:assert data.get('warnings'),name;assert data['config']['services'],name
root=base/'valid';run('after-init',[binary,'init','--root',str(root),'--format','json','--select','app-space-go:go-entry-1'],root)
cfg=root/'mudarro.json';raw=cfg.read_bytes();c=json.loads(raw);app=next(s for s in c['services'] if s['id']=='app-space-go');app['commands']['start']['args'].append('stay');cfg.write_text(json.dumps(c,indent=2)+'\n');manual=cfg.read_bytes()
run('after-build',[binary,'run','app-space-go:build','--root',str(root)],root);run('after-test',[binary,'run','app-space-go:test','--root',str(root)],root)
app['commands']['probe']={'args':['go','run','.'],'group':'qualidade'};cfg.write_text(json.dumps(c,indent=2)+'\n');manual=cfg.read_bytes()
run('after-probe-inherit',[binary,'run','app-space-go:probe','--root',str(root)],root)
app['go_workspace']='off';cfg.write_text(json.dumps(c,indent=2)+'\n');run('after-probe-off-expected-fail',[binary,'run','app-space-go:probe','--root',str(root)],root,False);app.pop('go_workspace');cfg.write_bytes(manual)
run('after-generate',[binary,'generate','--root',str(root)],root)
def hashes():return {str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in root.rglob('*') if p.is_file() and ('.mudarro' in p.parts or p.name=='menu.sh')}
h=hashes();run('after-generate-repeat',[binary,'generate','--root',str(root)],root);assert h==hashes();assert cfg.read_bytes()==manual
run('after-wrapper',['bash',str(root/'.mudarro/scripts/app-space-go/probe.sh')],root)
try:
 run('after-up',[binary,'run','app-space-go:up','--root',str(root)],root);run('after-status',[binary,'run','app-space-go:status','--root',str(root)],root);run('after-logs',[binary,'run','app-space-go:logs','--root',str(root)],root)
finally:run('after-down',[binary,'run','app-space-go:down','--root',str(root)],root)
assert cfg.read_bytes()==manual
(base/'after-summary.json').write_text(json.dumps({'passed':True,'checks':checks,'owned_service_stopped':True,'manual_configuration_preserved':True,'generated_hashes_idempotent':True,'binary':binary,'binary_sha256':hashlib.sha256(Path(binary).read_bytes()).hexdigest()},indent=2)+'\n');print('PASS',len(checks))
