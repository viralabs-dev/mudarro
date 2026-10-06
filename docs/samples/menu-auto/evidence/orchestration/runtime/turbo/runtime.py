import pathlib,os,sys,subprocess,json,shutil,hashlib,time
base=pathlib.Path('/tmp/mudarro-mud037-turbo-runtime'); root=base/'fixture'; logs=base/'logs'; logs.mkdir(exist_ok=True)
installed=pathlib.Path(sys.argv[1]); turbo=installed/'turbo/node_modules/.bin/turbo'
nodebin=pathlib.Path('/home/danielsouza/.nvm/versions/node/v24.15.0/bin'); binary=pathlib.Path('/tmp/mudarro-mud037-root/mudarro')
source=pathlib.Path('/home/danielsouza/dev/viralabs/mudarro/docs/samples/menu-auto/evidence/orchestration/samples/turbo')
if not root.exists(): shutil.copytree(source,root)
for d in ['config','cache','tmp','xdg-config','xdg-cache']: (base/d).mkdir(exist_ok=True)
for f in ['user.rc','global.rc']: (base/'config'/f).write_text('')
env={'PATH':str(binary.parent)+':'+str(turbo.parent)+':'+str(nodebin)+':/usr/bin:/bin','LANG':'C.UTF-8','LC_ALL':'C.UTF-8','TERM':'dumb','TMPDIR':str(base/'tmp'),'XDG_CONFIG_HOME':str(base/'xdg-config'),'XDG_CACHE_HOME':str(base/'xdg-cache'),'TURBO_TELEMETRY_DISABLED':'1','TURBO_CACHE':'local:rw','TURBO_CACHE_DIR':str(base/'cache'),'NPM_CONFIG_USERCONFIG':str(base/'config/user.rc'),'NPM_CONFIG_GLOBALCONFIG':str(base/'config/global.rc'),'NPM_CONFIG_CACHE':str(base/'cache/npm'),'NPM_CONFIG_AUDIT':'false','NPM_CONFIG_FUND':'false','NPM_CONFIG_IGNORE_SCRIPTS':'true'}
results=[]
def run(name,args,expect=0):
 t=time.monotonic();p=subprocess.run([str(a) for a in args],cwd=root,env=env,capture_output=True,text=True,timeout=60)
 (logs/(name+'.stdout')).write_text(p.stdout);(logs/(name+'.stderr')).write_text(p.stderr)
 results.append({'name':name,'argv':[str(a) for a in args],'exit':p.returncode,'elapsed':round(time.monotonic()-t,3),'stdout':name+'.stdout','stderr':name+'.stderr'})
 (base/'results.json').write_text(json.dumps(results,indent=2))
 assert p.returncode==expect,(name,p.returncode,p.stdout,p.stderr)
 return p
run('version',[turbo,'--version']);assert (logs/'version.stdout').read_text().strip()=='2.11.7'
run('lock',[nodebin/'npm','--userconfig',base/'config/user.rc','--globalconfig',base/'config/global.rc','--cache',base/'cache/npm','install','--package-lock-only','--ignore-scripts','--no-audit','--no-fund','--offline'])
run('workspace-install',[nodebin/'npm','--userconfig',base/'config/user.rc','--globalconfig',base/'config/global.rc','--cache',base/'cache/npm','ci','--ignore-scripts','--no-audit','--no-fund','--offline'])
common=['--cache=local:rw','--cache-dir',base/'cache/turbo','--summarize']
run('build-cold',[turbo,'run','build',*common])
for name,value in [('lib',21),('app',42)]:assert json.loads((root/'packages'/name/'dist/result.json').read_text())['value']==value
run('test-real',[turbo,'run','test',*common])
run('build-repeat-cache',[turbo,'run','build',*common])
(root/'packages/lib/input.txt').write_text('owned input invalidation\n')
run('build-input-change',[turbo,'run','build',*common])
run('failure',[turbo,'run','fail',*common],7)
run('missing-task',[turbo,'run','absent-task',*common],1)
scan=json.loads(run('scan',[binary,'scan','--root',root,'--json']).stdout);sus=[s for s in scan['suggestions'] if s.get('purpose')=='orchestration'];assert len(sus)==3,sus
sel=','.join(s['service']+':'+s['name'] for s in sus);run('init',[binary,'init','--root',root,'--select',sel]);config=(root/'mudarro.yaml').read_bytes();run('init-preservation',[binary,'init','--root',root],1);assert (root/'mudarro.yaml').read_bytes()==config
run('generate-1',[binary,'generate','--root',root]);hashes=lambda:{str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in root.rglob('*') if p.is_file() and 'node_modules' not in p.parts and '.turbo' not in p.parts};one=hashes();run('generate-2',[binary,'generate','--root',root]);assert hashes()==one
for s in sus:run('preview-'+s['name'],[binary,'preview','--root',root,s['service']+':'+s['name']])
for task,expected in [('build',0),('test',0),('fail',1)]:
 su=next(s for s in sus if s['name']=='turbo-'+task)
 run('mudarro-'+task,[binary,'run','--root',root,su['service']+':'+su['name']],expected)
 run('wrapper-'+task,[root/'.mudarro/scripts'/su['service']/(su['name']+'.sh')],expected)
# Preserve source/lock manifest list, runtime binaries themselves stay outside samples.
(base/'source-hashes.json').write_text(json.dumps(hashes(),indent=2))
print('Turbo runtime and Mudarro gates complete; inspect run summaries for cache/order before acceptance')
