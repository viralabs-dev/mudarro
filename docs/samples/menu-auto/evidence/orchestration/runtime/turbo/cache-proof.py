import pathlib,json
script=pathlib.Path('/tmp/mudarro-mud037-turbo-runtime/runtime.py').read_text();exec(script.split("run('version'")[0])
results=json.loads((base/'results.json').read_text())
conf=json.loads((root/'turbo.json').read_text());conf['globalDependencies']=['scripts/**']
for name in ['build','test']:conf['tasks'][name]['inputs']=['package.json','input.txt']
(root/'turbo.json').write_text(json.dumps(conf,indent=2)+'\n')
(base/'env.json').write_text(json.dumps(env,indent=2))
common=['--cache=local:rw','--cache-dir',base/'cache/turbo','--summarize']
def proof(name):
 before=set((root/'.turbo/runs').glob('*.json'));run(name,[turbo,'run','build',*common]);after=set((root/'.turbo/runs').glob('*.json'));new=after-before;assert len(new)==1
 f=new.pop();d=json.loads(f.read_text());(logs/(name+'.summary.json')).write_text(json.dumps(d,indent=2));return d
cold=proof('build-explicit-inputs-cold');hit=proof('build-explicit-inputs-repeat')
assert all(t['cache']['status']=='MISS' for t in cold['tasks'])
assert all(t['cache']['status']=='HIT' and t['cache']['local'] and not t['cache']['remote'] for t in hit['tasks'])
(root/'packages/lib/input.txt').write_text('owned input invalidation second revision\n')
changed=proof('build-explicit-inputs-change');assert all(t['cache']['status']=='MISS' for t in changed['tasks'])
assert {t['taskId']:t['hash'] for t in hit['tasks']}!={t['taskId']:t['hash'] for t in changed['tasks']}
lib=next(t for t in cold['tasks'] if t['taskId'].endswith('lib#build'));app=next(t for t in cold['tasks'] if t['taskId'].endswith('app#build'));assert lib['execution']['endTime']<=app['execution']['startTime']
run('test-explicit-inputs-real',[turbo,'run','test',*common]);run('wrapper-test-final',[root/'.mudarro/scripts'/next(s['service'] for s in json.loads((logs/'scan.stdout').read_text())['suggestions'] if s['name']=='turbo-test')/'turbo-test.sh'])
summary={'tool':'turbo','version':'2.11.7','node':'24.15.0','npm':'11.12.1','runtimeRoot':str(installed),'fixture':str(root),'env':'env.json','commands':len(results),'buildValues':{'lib':21,'app':42},'dependencyOrder':'lib.endTime <= app.startTime','cache':{'cold':'MISS both','repeat':'local HIT both; remote false','inputChange':'MISS both; hashes changed'},'negatives':{'childTaskExit':7,'turboExit':7,'mudarroAndWrapperExit':1,'missingTaskExit':1},'mudarro':['scan','init --select','init preserves existing','generate twice identical','preview all','run build/test/fail','generated wrappers build/test/fail'],'initialLimitations':['--cache and --force mutually exclusive in 2.11.7; corrected to initially empty local cache','No Git default inputs included output/log files; corrected fixture to explicit task inputs + globalDependencies scripts/**'],'cacheProofLogs':['build-explicit-inputs-cold.summary.json','build-explicit-inputs-repeat.summary.json','build-explicit-inputs-change.summary.json'],'noPersistentAppStarted':True}
(base/'summary.json').write_text(json.dumps(summary,indent=2));(base/'source-hashes.json').write_text(json.dumps({str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in root.rglob('*') if p.is_file() and 'node_modules' not in p.parts and '.turbo' not in p.parts},indent=2))
print('Actual Turbo cache/order/invalidation verified; recorder env ready')
