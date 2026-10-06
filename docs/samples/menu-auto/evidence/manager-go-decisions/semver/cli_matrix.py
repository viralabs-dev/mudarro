import pathlib,json,subprocess,sys,hashlib
base=pathlib.Path('/tmp/mudarro-mud035-semver');binary=pathlib.Path(sys.argv[1]);phase=sys.argv[2];results=[]
for i,value in enumerate(['npm@1.2.3','yarn@4.18.1+sha224.abcdef','npm@1.2.3-alpha.1+build.001','npm@latest','npm@01.2.3','npm@1.2.3-alpha.01','npm@1.2.3+build_1']):
 root=base/(phase+'-fixture-'+str(i));root.mkdir(exist_ok=True)
 body=json.dumps({'packageManager':value,'scripts':{'test':'touch NEVER_EXECUTE'}});(root/'package.json').write_text(body)
 args=[str(binary),'scan','--root',str(root),'--json'];p=subprocess.run(args,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=3)
 report=json.loads(p.stdout);starts=[s.get('manager') for s in report['config']['services']];warnings=report.get('warnings') or []
 for s in report['config']['services']:warnings+=s.get('pending',[])
 result={'packageManager':value,'argv':args,'exit':p.returncode,'managers':starts,'diagnostics':warnings};results.append(result)
 assert (root/'package.json').read_text()==body and not(root/'NEVER_EXECUTE').exists()
 if phase=='after':
  if i<3:assert starts[0] in ['npm','yarn'],result
  else:assert not any(starts) and any('packageManager' in w for w in warnings),result
(base/('cli-'+phase+'.json')).write_text(json.dumps({'binarySHA256':hashlib.sha256(binary.read_bytes()).hexdigest(),'results':results},indent=2)+'\n');print(phase,len(results),'CLI scans complete; preservationPASS')
