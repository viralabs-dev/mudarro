import pathlib,subprocess,os,json,sys
base=pathlib.Path('/tmp/mudarro-mud038-next/diagnostics');root=base/'cli-fixture';tools=base/'cli-tools';root.mkdir(exist_ok=True);tools.mkdir(exist_ok=True)
config={'version':1,'name':'Go diagnostic exclusive','services':[{'id':'app','dir':'.','language':'go','go_workspace':'off','infrastructure':{'kind':'custom'},'commands':{'test':{'args':['go','test','./...'],'group':'qualidade'}}}]}
(root/'mudarro.json').write_text(json.dumps(config));(tools/'go').write_text('#!/bin/sh\nprintf bad > "'+str(root/'MUST_NOT_EXECUTE')+'"\n');os.chmod(tools/'go',0o755)
results=[]
for args in [['preview','app:test'],['run','app:test','--dry-run'],['doctor']]:
 full=[sys.argv[1]]+args+['--root',str(root)];p=subprocess.run(full,env=dict(os.environ,PATH=str(tools)),text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=3);results.append({'argv':full,'exit':p.returncode,'output':p.stdout})
assert not(root/'MUST_NOT_EXECUTE').exists()
if sys.argv[2]=='after':assert all('GOWORK=off' in r['output'] for r in results[:2]) and results[2]['exit']==1 and 'AUSENTE app: env' in results[2]['output']
(base/('cli-'+sys.argv[2]+'.json')).write_text(json.dumps(results,indent=2)+'\n');print(sys.argv[2],len(results),'CLI diagnostics completed; no action execution')
