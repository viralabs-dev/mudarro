import os,pathlib,subprocess,json
base=pathlib.Path('/tmp/mudarro-runtimes.K7F4J1'); logs=base/'logs/yarn';bin='/tmp/mudarro-timeout-checkpoint/mudarro'
classicbin=base/'yarn-classic/bin';classicbin.mkdir(exist_ok=True)
classicrc=base/'config/yarn-classic.rc'
classicrc.write_text('prefix '+str(base/'yarn-classic/private-prefix')+'\nglobal-folder '+str(base/'yarn-classic/global')+'\ncache-folder '+str(base/'cache/npm-classic/runtime')+'\n')
(classicbin/'yarn').write_text('#!/bin/sh\nexec '+str(base/'yarn-classic/node_modules/.bin/yarn')+' --no-default-rc --use-yarnrc '+str(classicrc)+' "$@"\n');(classicbin/'yarn').chmod(0o755)
for name,version in [('classic','1.22.22'),('modern','4.18.1')]:
 root=base/'fixtures/yarn'/name;root.mkdir(exist_ok=True)
 (root/'package.json').write_text(json.dumps({'name':'mudarro-yarn-'+name,'private':True,'packageManager':'yarn@'+version,'scripts':{'start':'node app.cjs','test':'node probe.cjs','hello':'node hello.cjs','fail':'node fail.cjs'}},indent=2)+'\n')
 (root/'app.cjs').write_text("console.log('owned yarn app ready'); setInterval(()=>{}, 1000);\n")
 (root/'probe.cjs').write_text("const fs=require('node:fs');fs.writeFileSync('probe-marker.json',JSON.stringify({cwd:process.cwd(),args:process.argv.slice(2)}));console.log('Yarn quality real OK');\n")
 (root/'hello.cjs').write_text("require('node:fs').writeFileSync('hello-marker.txt',process.cwd());console.log('Yarn custom real OK');\n")
 (root/'fail.cjs').write_text("console.error('Yarn real failure');process.exit(7);\n")
 (root/'.mudarro-isolated-yarnrc.yml').write_text('enableGlobalCache: false\nnodeLinker: node-modules\nenableTelemetry: false\n')
 env={k:v for k,v in os.environ.items() if not k.startswith('YARN_') and not k.lower().startswith('npm_config_') and k not in ['NPM_TOKEN','NODE_AUTH_TOKEN']}
 env.update(PATH=str(base/('yarn-'+name)/'bin')+':'+str(pathlib.Path(bin).parent)+':'+env['PATH'],COREPACK_HOME=str(base/'cache/corepack'),COREPACK_ENABLE_PROJECT_SPEC='0',COREPACK_ENABLE_NETWORK='0',YARN_CACHE_FOLDER=str(base/('cache/npm-classic/runtime' if name=='classic' else 'yarn-modern/cache')),YARN_GLOBAL_FOLDER=str(base/('yarn-'+name)/'global'))
 env.update(NPM_CONFIG_USERCONFIG=str(base/'config/yarn-npm-user.rc'),NPM_CONFIG_GLOBALCONFIG=str(base/'config/yarn-npm-global.rc'))
 if name=='modern':env['PATH']='/home/danielsouza/.nvm/versions/node/v24.15.0/bin:'+env['PATH']
 if name=='modern':env.update(YARN_RC_FILENAME='.mudarro-isolated-yarnrc.yml',YARN_ENABLE_NETWORK='0',YARN_ENABLE_TELEMETRY='0',YARN_NPM_REGISTRY_SERVER='https://registry.npmjs.org')
 args=['yarn','install']+(['--offline','--ignore-scripts','--non-interactive'] if name=='classic' else ['--mode=skip-build'])
 run=subprocess.run(args,cwd=root,env=env,capture_output=True,text=True);(logs/(name+'-fixture-install.log')).write_text(run.stdout+run.stderr);assert run.returncode==0,(name,run.stdout,run.stderr)
 assert (root/'yarn.lock').exists(),name
 for label,args in [('scan',['scan','--json']),('init',['init','--select','app-javascript:hello,app-javascript:fail']),('generate',['generate']),('generate-repeat',['generate']),('test',['run','app-javascript:test']),('hello',['run','app-javascript:hello']),('negative-fail',['run','app-javascript:fail'])]:
  run=subprocess.run([bin,*args,'--root',str(root)],cwd=root,env=env,capture_output=True,text=True);(logs/(name+'-'+label+'.log')).write_text(run.stdout+run.stderr)
  assert (run.returncode!=0 if label=='negative-fail' else run.returncode==0),(name,label,run.stdout,run.stderr)
 (logs/(name+'-environment.json')).write_text(json.dumps({k:v for k,v in env.items() if k=='PATH' or k.startswith('YARN_') or k.startswith('COREPACK_')},indent=2)+'\n')
 print(name+': install real lock + scan/init/generate2/test/hello/failure PASS')
