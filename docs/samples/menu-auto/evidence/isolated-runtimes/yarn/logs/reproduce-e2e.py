import os,pathlib,json,subprocess,hashlib,time
base=pathlib.Path('/tmp/mudarro-runtimes.K7F4J1');logs=base/'logs/yarn';bin='/tmp/mudarro-timeout-checkpoint/mudarro';results={}
for name in ['classic','modern']:
 root=base/'fixtures/yarn'/name;env=os.environ.copy();env.update(json.loads((logs/(name+'-environment.json')).read_text()))
 checks=[]
 def run(label,args,expected=0,input=None):
  p=subprocess.run(args,cwd=root,env=env,input=input,capture_output=True,text=True,timeout=20);(logs/(name+'-'+label+'.log')).write_text(p.stdout+p.stderr);assert p.returncode==expected,(name,label,p.returncode,p.stdout,p.stderr);checks.append({'case':label,'exit_code':p.returncode});return p
 def hashes():return {str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in [root/'package.json',root/'yarn.lock',root/'mudarro.yaml',root/'menu.sh',root/'.mudarro/generated.json',*list((root/'.mudarro/scripts').rglob('*.sh'))]}
 before=hashes()
 run('generated-install',['bash',str(root/'.mudarro/scripts/app-javascript/install.sh')])
 assert hashes()==before,(name,'installchangedmanifestlockorwrapper')
 run('wrapper-test',['bash',str(root/'.mudarro/scripts/app-javascript/test.sh')]);run('wrapper-hello',['bash',str(root/'.mudarro/scripts/app-javascript/hello.sh')])
 (root/'probe-marker.json').unlink()
 menu_result=run('menu-test',['bash',str(root/'menu.sh')],input='1\n4\n1\n0\n0\n0\n')
 assert 'Yarn quality real OK' in menu_result.stdout,(name,'menu did not execute quality action')
 assert json.loads((root/'probe-marker.json').read_text())['cwd']==str(root)
 assert (root/'hello-marker.txt').read_text()==str(root)
 run('init-repeat',['/'+str(pathlib.Path(bin)).lstrip('/'),'init','--root',str(root)],expected=1)
 run('generate-idempotence',[bin,'generate','--root',str(root)]);assert hashes()==before
 run('unknown-action',[bin,'run','app-javascript:unknown','--root',str(root)],expected=1)
 try:
  run('up',[bin,'run','app-javascript:up','--root',str(root)]);run('up-repeat',[bin,'run','app-javascript:up','--root',str(root)])
  status=run('status',[bin,'run','app-javascript:status','--root',str(root)]);assert 'running' in status.stdout
  run('restart',[bin,'run','app-javascript:restart','--root',str(root)]);run('logs',[bin,'run','app-javascript:logs','--root',str(root)])
  run('down',[bin,'run','app-javascript:down','--root',str(root)]);status=run('status-stopped',[bin,'run','app-javascript:status','--root',str(root)]);assert 'stopped' in status.stdout
 finally:
  subprocess.run([bin,'run','app-javascript:down','--root',str(root)],cwd=root,env=env,capture_output=True,timeout=10)
 # Restore only our deliberate user edit, never regenerate over it.
 menu=(root/'menu.sh').read_bytes();(root/'menu.sh').write_bytes(menu+b'# owned preservation probe\n');edited=hashes();run('manual-edit-rejection',[bin,'generate','--root',str(root)],expected=1);assert hashes()==edited;(root/'menu.sh').write_bytes(menu)
 checks.append({'case':'manifest_lock_generated_hashes_preserved','passed':True})
 if name=='modern':
  frozen=run('frozen-compatibility',['yarn','install','--frozen-lockfile']);checks.append({'case':'frozen-lockfile_compatibility','passed':True,'note':'accepted but deprecated; inspect log; prefer immutable for Modern'})
 results[name]={'passed':True,'checks':checks,'fixture':str(root),'lockfile_sha256':hashlib.sha256((root/'yarn.lock').read_bytes()).hexdigest(),'configuration_sha256':hashlib.sha256((root/'mudarro.yaml').read_bytes()).hexdigest(),'processes_stopped':True,'runtime':'real Yarn '+('1.22.22' if name=='classic' else '4.18.1'),'network_disabled_during_fixture_runtime':name=='modern','classic_offline_initial_install':name=='classic'}
 print(name+': menu/wrappers/install/supervisor/negative/idempotence PASS, stopped')
(logs/'e2e-summary.json').write_text(json.dumps(results,indent=2)+'\n')
