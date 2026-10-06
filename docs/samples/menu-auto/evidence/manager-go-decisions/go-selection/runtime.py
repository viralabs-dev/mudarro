import pathlib,os,subprocess,json,hashlib
base=pathlib.Path('/tmp/mudarro-mud038-selection'); root=base/'runtime-final2';root.mkdir(exist_ok=True)
for name in ['member','other']:
 p=root/name;p.mkdir(exist_ok=True);(p/'go.mod').write_text('module example.test/'+name+'\ngo 1.23\n');(p/'main.go').write_text('package main\nimport("fmt";"os";"time")\nfunc main(){fmt.Println("REAL_GO_'+name.upper()+'",os.Getenv("GOWORK"));if len(os.Args)>1 {for {time.Sleep(time.Second)}}}\n')
(root/'go.work').write_text('go 1.23\nuse ./member\n')
env=dict(os.environ);env.update(PATH=str(base)+':/home/danielsouza/sdk/go1.27.1/bin:'+env['PATH'],GOCACHE=str(base/'go-cache'),GOTOOLCHAIN='local',GOPROXY='off',GOFLAGS='-buildvcs=false');env.pop('GOWORK',None)
binary=str(base/'mudarro');results=[]
def run(label,args,ok=True,cwd=root):
 p=subprocess.run(args,cwd=cwd,env=env,capture_output=True,text=True,timeout=60);(base/(label+'.log')).write_text(p.stdout+p.stderr);results.append({'case':label,'exit_code':p.returncode,'expected_success':ok});assert (p.returncode==0)==ok,(label,p.stdout,p.stderr);return p
run('member-inherit',['go','run','.'],cwd=root/'member')
run('nonmember-inherit',['go','run','.'],False,root/'other')
run('init',[binary,'init','--root',str(root),'--format','json'])
cfg=root/'mudarro.json';c=json.loads(cfg.read_text());other=next(s for s in c['services'] if s['dir']=='other');member=next(s for s in c['services'] if s['dir']=='member');oid=other['id'];mid=member['id']
run('cli-nonmember-inherit',[binary,'run',oid+':start','--root',str(root)],False)
other['go_workspace']='off';other['commands']['start']['args']=['go','run','.','stay'];other['commands']['probe']={'args':['go','run','.'],'group':'qualidade'};cfg.write_text(json.dumps(c,indent=2)+'\n');before=cfg.read_bytes()
run('cli-off',[binary,'run',oid+':probe','--root',str(root)])
(root/'go.work').write_text('go 1.23\nuse ./missing\n')
run('missing-inherit',[binary,'run',mid+':start','--root',str(root)],False)
run('missing-off',[binary,'run',oid+':probe','--root',str(root)])
run('generate',[binary,'generate','--root',str(root)])
files=lambda:{str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in root.rglob('*') if p.is_file() and ('.mudarro' in p.parts or p.name=='menu.sh')}
hashes=files();run('generate-repeat',[binary,'generate','--root',str(root)]);assert hashes==files();assert cfg.read_bytes()==before
run('wrapper-off',['bash',str(root/'.mudarro/scripts'/oid/'probe.sh')])
try:
 run('supervisor-up',[binary,'run',oid+':up','--root',str(root)])
 run('supervisor-status',[binary,'run',oid+':status','--root',str(root)])
 run('supervisor-logs',[binary,'run',oid+':logs','--root',str(root)])
finally:run('supervisor-down',[binary,'run',oid+':down','--root',str(root)])
assert cfg.read_bytes()==before
run('existing-init',[binary,'init','--root',str(root)],False);assert cfg.read_bytes()==before
(base/'runtime-summary.json').write_text(json.dumps({'passed':True,'checks':results,'manual_configuration_preserved':True,'generated_hashes_idempotent':True,'scoped_environment':{k:env[k] for k in ['GOCACHE','GOTOOLCHAIN','GOPROXY','GOFLAGS']},'global_GOWORK_not_modified':os.environ.get('GOWORK'),'membership_parser_validated':False},indent=2)+'\n')
print('real Go offline CLI/wrapper/supervisor PASS',len(results))
