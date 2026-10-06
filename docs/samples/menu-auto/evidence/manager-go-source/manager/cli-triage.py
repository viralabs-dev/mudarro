import hashlib,json,pathlib,subprocess,tempfile
binary='/tmp/mudarro-mud035/mudarro'
cases=[('unknown','unknown'),('relative-path','../manager'),('absolute-path','/tmp/manager'),('url','https://example.invalid/tool'),('trailing-space','npm '),('leading-space',' npm'),('empty-version','npm@'),('empty-name','@1.2.3'),('tab','npm\t'),('newline','npm\n'),('number',123),('boolean',True),('array',[]),('object',{})]
results=[]
for name,value in cases:
 with tempfile.TemporaryDirectory(prefix='mudarro-pm-cli-') as tmp:
  root=pathlib.Path(tmp)
  for folder,body in [('bad',{'packageManager':value,'scripts':{'test':'touch SHOULD_NEVER_EXECUTE'}}),('good',{'packageManager':'npm','scripts':{'test':'node --version'}})]:
   (root/folder).mkdir();(root/folder/'package.json').write_text(json.dumps(body))
  (root/'menu.sh').write_text('user launcher preserved\n');(root/'.mudarro').mkdir();(root/'.mudarro/generated.json').write_text('historical fixture bytes preserved\n')
  def hashes():return {str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in root.rglob('*') if p.is_file()}
  before=hashes();run=subprocess.run([binary,'scan','--root',str(root),'--json'],capture_output=True,text=True,check=True);report=json.loads(run.stdout)
  services=report['config']['services'];bad=[s for s in services if s['dir']=='bad'];good=[s for s in services if s['dir']=='good']
  diagnostics=report.get('warnings',[]) or []
  invalid_args=any(any(c.get('args') for c in (s.get('commands') or {}).values()) for s in bad)
  invalid_args=invalid_args or any(s.get('service')=='bad-javascript' and (s.get('command') or {}).get('args') for s in report.get('suggestions') or [])
  assert good and good[0]['manager']=='npm' and not invalid_args and hashes()==before and not (root/'SHOULD_NEVER_EXECUTE').exists(),name
  assert any('packagemanager' in d.lower() for d in diagnostics), (name,diagnostics)
  result={'case':name,'exit_code':run.returncode,'diagnostic':True,'invalid_executable_plan':False,'valid_sibling':True,'files_preserved':True,'offline_scan_no_script_execution':True};results.append(result);print(json.dumps(result))
pathlib.Path('/tmp/mudarro-mud035/cli-triage-summary.json').write_text(json.dumps({'passed':len(results),'cases':results,'binary_sha256':hashlib.sha256(pathlib.Path(binary).read_bytes()).hexdigest(),'runtime':'real built CLI scan; no package-manager command invoked'},indent=2)+'\n')
