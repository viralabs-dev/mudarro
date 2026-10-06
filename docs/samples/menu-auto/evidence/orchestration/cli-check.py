import pathlib,tempfile,shutil,subprocess,json,hashlib
binary='/tmp/mudarro-mud037-root/mudarro';base=pathlib.Path('/home/danielsouza/dev/viralabs/mudarro/docs/samples/menu-auto/evidence/orchestration/samples');results=[]
for profile in ['turbo','nx']:
 root=pathlib.Path(tempfile.mkdtemp(prefix='mudarro-037-cli-'))/profile;shutil.copytree(base/profile,root)
 def run(args,expected=0):
  p=subprocess.run([binary,*args,'--root',str(root)],capture_output=True,text=True);assert p.returncode==expected,(args,p.returncode,p.stderr);results.append({'profile':profile,'args':args,'exit':p.returncode,'stdout':p.stdout,'stderr':p.stderr});return p
 scan=json.loads(run(['scan','--json']).stdout);suggestions=[s for s in scan['suggestions'] if s.get('purpose')=='orchestration'];assert suggestions
 sel=','.join(s['service']+':'+s['name'] for s in suggestions);run(['init','--select',sel]);config=root/'mudarro.yaml';before=config.read_bytes();run(['init'],1);assert config.read_bytes()==before
 run(['generate']);files=lambda:{str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in root.rglob('*') if p.is_file()};one=files();run(['generate']);assert files()==one
 for s in suggestions:run(['preview',s['service']+':'+s['name']])
 results.append({'profile':profile,'manualConfigPreserved':True,'generateIdempotent':True,'orchestrationSuggestions':len(suggestions),'runtimeExecuted':False})
pathlib.Path('/tmp/mudarro-mud037-root/cli-evidence.json').write_text(json.dumps(results,indent=2));print('CLI scan/select/manual preservation/generate twice/previews: PASS both samples')
