from pathlib import Path
import json,os,subprocess
base=Path('/tmp/mudarro-mud036-runtime');binary='/tmp/mudarro-mud038-parser/mudarro'
profiles={'npm':'npm@11.19.0','pnpm':'pnpm@11.3.0','classic':'yarn@1.22.22','modern':'yarn@4.18.1','bun':'bun@1.4.0'}
for name,pm in profiles.items():
 root=base/name;root.mkdir(exist_ok=True)
 manifest={'name':'mudarro-workspace-'+name,'private':True,'packageManager':pm,'workspaces':['packages/*'],'scripts':{'test':'node root.cjs'}}
 (root/'package.json').write_text(json.dumps(manifest,indent=2)+'\n');(root/'root.cjs').write_text("console.log('ROOT_ONLY');\n")
 if name=='pnpm':(root/'pnpm-workspace.yaml').write_text('packages:\n  - packages/*\n')
 for child in ['a','b']:
  p=root/'packages'/child;p.mkdir(parents=True,exist_ok=True)
  (p/'package.json').write_text(json.dumps({'name':'mudarro-child-'+name+'-'+child,'private':True,'scripts':{'start':'node app.cjs','test':'node probe.cjs','hello':'node probe.cjs'}},indent=2)+'\n')
  (p/'app.cjs').write_text("console.log('WORKSPACE_CHILD_"+child+"',process.cwd());setInterval(()=>{},1000);\n")
  (p/'probe.cjs').write_text("require('node:fs').writeFileSync('OWN_CWD_MARKER',process.cwd());console.log('WORKSPACE_CHILD_"+child+"',process.cwd());\n")
 env=dict(os.environ)
 if name in ['classic','modern']:env.update(json.loads(Path('/tmp/mudarro-runtimes.K7F4J1/logs/yarn/'+name+'-environment.json').read_text()))
 env.update(NPM_CONFIG_USERCONFIG='/tmp/mudarro-runtimes.K7F4J1/config/yarn-npm-user.rc',NPM_CONFIG_GLOBALCONFIG='/tmp/mudarro-runtimes.K7F4J1/config/yarn-npm-global.rc')
 env={k:v for k,v in env.items() if k in ['PATH','NPM_CONFIG_USERCONFIG','NPM_CONFIG_GLOBALCONFIG'] or k.startswith('YARN_') or k.startswith('COREPACK_')}
 (base/(name+'-environment.json')).write_text(json.dumps(env,indent=2)+'\n')
 p=subprocess.run([binary,'scan','--json','--root',str(root)],capture_output=True,text=True,timeout=20);(base/(name+'-baseline-scan.json')).write_text(p.stdout);(base/(name+'-baseline-scan.stderr')).write_text(p.stderr);assert p.returncode==0
 data=json.loads(p.stdout);print(name,[(s['dir'],s.get('manager'),s.get('pending')) for s in data['config']['services']])
