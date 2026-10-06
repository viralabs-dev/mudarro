from pathlib import Path
import subprocess,json,hashlib,time,os,uuid
base=Path('/tmp/mudarro-approved-runtimes.6AsmDd/dotnet');prefix=base/'prefix/dotnet-10.0.401';sdk=prefix/'dotnet';binary=Path('/tmp/mudarro-next-languages-root/mudarro');results=[]
runroot=base/'runs'/('verified-'+uuid.uuid4().hex[:8]);(runroot/'fixtures').mkdir(parents=True)
env={'PATH':str(prefix)+':'+str(binary.parent)+':/usr/bin:/bin','LANG':'C.UTF-8','LC_ALL':'C.UTF-8','TERM':'dumb','DOTNET_ROOT':str(prefix),'DOTNET_CLI_HOME':str(base/'cache/dotnet-home'),'NUGET_PACKAGES':str(base/'cache/nuget-packages'),'NUGET_HTTP_CACHE_PATH':str(base/'cache/nuget-http'),'NUGET_PLUGINS_CACHE_PATH':str(base/'cache/nuget-plugins'),'DOTNET_CLI_TELEMETRY_OPTOUT':'1','DOTNET_SKIP_FIRST_TIME_EXPERIENCE':'1','DOTNET_ADD_GLOBAL_TOOLS_TO_PATH':'0','DOTNET_GENERATE_ASPNET_CERTIFICATE':'false','DOTNET_NOLOGO':'1','DOTNET_MULTILEVEL_LOOKUP':'0','DOTNET_CLI_WORKLOAD_UPDATE_NOTIFY_DISABLE':'1','DOTNET_CLI_USE_MSBUILD_SERVER':'0','DOTNET_CLI_DO_NOT_USE_MSBUILD_SERVER':'1','MSBUILDDISABLENODEREUSE':'1','TMPDIR':str(base/'tmp')}
# Preserve inherited HOME verbatim; all SDK/NuGet writable homes remain explicit private paths.
if 'HOME' in os.environ: env['HOME']=os.environ['HOME']
for name in ['home','cache/dotnet-home','cache/nuget-packages','cache/nuget-http','cache/nuget-plugins','tmp','fixtures','logs']:(base/name).mkdir(parents=True,exist_ok=True)
(base/'evidence/runtime-environment.json').write_text(json.dumps(env,indent=2)+'\n')
def run(profile,name,argv,expected=0,cwd=None):
 cwd=Path(cwd) if cwd else base;started=time.monotonic();p=subprocess.run([str(a)for a in argv],cwd=cwd,env=env,capture_output=True,text=True,timeout=90);record={'profile':profile,'name':name,'argv':[str(a)for a in argv],'cwd':str(cwd),'exit':p.returncode,'expectedExit':expected,'elapsedSeconds':round(time.monotonic()-started,3),'stdout':p.stdout,'stderr':p.stderr};results.append(record);(base/'evidence/command-results.json').write_text(json.dumps(results,indent=2)+'\n');assert p.returncode==expected,(profile,name,p.returncode,p.stdout,p.stderr);return p
run('toolchain','sdk-version',[sdk,'--version']);run('toolchain','sdk-info',[sdk,'--info']);assert results[0]['stdout'].strip()=='10.0.401'
global_json=json.dumps({'sdk':{'version':'10.0.401','rollForward':'disable','allowPrerelease':False}},indent=2)+'\n'
nuget='<configuration><packageSources><clear/></packageSources><fallbackPackageFolders><clear/></fallbackPackageFolders><disabledPackageSources><clear/></disabledPackageSources></configuration>\n'
project='<Project Sdk="Microsoft.NET.Sdk"><PropertyGroup><OutputType>Exe</OutputType><TargetFramework>net10.0</TargetFramework><ImplicitUsings>disable</ImplicitUsings><UseSharedCompilation>false</UseSharedCompilation><RestoreDisableParallel>true</RestoreDisableParallel></PropertyGroup></Project>\n'
source='''using System;
internal static class Program {
 static int Main(string[] args) {
  int actual=21*2;
  if(args.Length==1 && args[0]=="--controlled-failure") { Console.WriteLine("DOTNET_CONTROLLED_FAILURE_EXPECTED_EXIT_7");return 7; }
  if(args.Length==1 && args[0]=="--self-test") {if(actual!=42){Console.Error.WriteLine("SELF_TEST_FAIL");return 9;}Console.WriteLine("DOTNET_REAL_STDLIB_SELF_TEST_PASS: 21 * 2 = 42");return 0;}
  Console.WriteLine("DOTNET_REAL_CONSOLE_OUTPUT: "+actual);return 0;
 }
}
'''
fixture=runroot/'fixtures/console';fixture.mkdir(exist_ok=True)
for name,body in [('global.json',global_json),('NuGet.Config',nuget),('Sample.csproj',project),('Program.cs',source)]: (fixture/name).write_text(body)
sourcehash=lambda:{name:hashlib.sha256((fixture/name).read_bytes()).hexdigest()for name in ['global.json','NuGet.Config','Sample.csproj','Program.cs']};before=sourcehash()
scan=json.loads(run('console','scan',[binary,'scan','--root',fixture,'--json']).stdout);assert scan==json.loads(run('console','scan-repeat',[binary,'scan','--root',fixture,'--json']).stdout);services=scan['config']['services'];assert len(services)==1 and services[0]['manager']=='dotnet' and not services[0].get('commands') and not services[0].get('framework');sid=services[0]['id'];assert [s['name']for s in scan['suggestions']]==['build','test'];assert not any(s['name']=='start'for s in scan['suggestions'])
run('console','init-select-build',[binary,'init','--root',fixture,'--format','json','--select',sid+':build']);config=json.loads((fixture/'mudarro.json').read_text());service=config['services'][0];assert service['commands']['build']['args']==['dotnet','build','./Sample.csproj'];assert 'test'not in service['commands'] and 'start'not in service['commands']
# User-declared argv extends the opted-in build with private/offline runtime policy.
service['commands']['build']['args']=['dotnet','build','./Sample.csproj','--no-restore','--disable-build-servers','-p:UseSharedCompilation=false','-nodeReuse:false']
for name,args in [('console',['dotnet','bin/Debug/net10.0/Sample.dll']),('self-test',['dotnet','bin/Debug/net10.0/Sample.dll','--self-test']),('controlled-failure',['dotnet','bin/Debug/net10.0/Sample.dll','--controlled-failure'])]:service['commands'][name]={'args':args,'group':'qualidade'if name!='console'else'aplicacao'}
(fixture/'mudarro.json').write_text(json.dumps(config,indent=2)+'\n');configbefore=(fixture/'mudarro.json').read_bytes()
run('console','restore-clear-feeds',[sdk,'restore','./Sample.csproj','--configfile',fixture/'NuGet.Config','--disable-parallel','--verbosity','minimal','-p:NuGetAudit=false'],cwd=fixture)
assets=json.loads((fixture/'obj/project.assets.json').read_text());assert not assets.get('libraries');assert not assets['project']['restore'].get('sources');(base/'evidence/restore-policy.json').write_text(json.dumps({'applicationPackageReference':False,'assetLibraries':assets.get('libraries',{}),'resolvedSources':assets['project']['restore'].get('sources',{}),'targetFramework':'net10.0','frameworkReferences':assets['project']['frameworks']['net10.0'].get('frameworkReferences',{}),'sdkPinned':'10.0.401','rollForward':'disable','nativeTestFrameworkExecuted':False},indent=2)+'\n')
run('console','generate-1',[binary,'generate','--root',fixture]);generatedhash=lambda:{str(p.relative_to(fixture)):hashlib.sha256(p.read_bytes()).hexdigest()for p in fixture.rglob('*.sh')};generatedbefore=generatedhash();run('console','generate-2',[binary,'generate','--root',fixture]);assert generatedbefore==generatedhash() and configbefore==(fixture/'mudarro.json').read_bytes()
run('console','init-existing-preserve',[binary,'init','--root',fixture],1);assert configbefore==(fixture/'mudarro.json').read_bytes()
for task,expected in [('build',0),('console',0),('self-test',0),('controlled-failure',1)]:
 run('console','preview-'+task,[binary,'preview','--root',fixture,sid+':'+task]);p=run('console','mudarro-run-'+task,[binary,'run','--root',fixture,sid+':'+task],expected);run('console','wrapper-'+task,[fixture/'.mudarro/scripts'/sid/(task+'.sh')],expected)
 if task=='self-test':assert 'DOTNET_REAL_STDLIB_SELF_TEST_PASS' in p.stdout
 if task=='console':assert 'DOTNET_REAL_CONSOLE_OUTPUT: 42'in p.stdout
run('console','direct-controlled-failure',[sdk,'bin/Debug/net10.0/Sample.dll','--controlled-failure'],7,cwd=fixture)
run('console','repeat-selected-build',[binary,'run','--root',fixture,sid+':build']);run('console','repeat-selected-self-test',[binary,'run','--root',fixture,sid+':self-test'])
# Compile failure is genuine compiler feedback; restore the source afterwards.
(fixture/'Program.cs').write_text(source+'\nTHIS_IS_A_CONTROLLED_COMPILE_ERROR\n');run('console','controlled-compile-failure',[binary,'run','--root',fixture,sid+':build'],1);(fixture/'Program.cs').write_text(source);run('console','recovery-build',[binary,'run','--root',fixture,sid+':build']);run('console','recovery-self-test',[binary,'run','--root',fixture,sid+':self-test']);assert sourcehash()==before and configbefore==(fixture/'mudarro.json').read_bytes()
for profile in ['multiple','invalid','excluded']:
 root=runroot/'fixtures'/profile;root.mkdir(exist_ok=True)
 if profile=='multiple':
  for name in ['-flag.csproj','space +.csproj']:(root/name).write_text('<Project/>')
 elif profile=='invalid':(root/'broken.csproj').write_text('<!DOCTYPE Project SYSTEM "https://example.invalid/must-not-fetch"><Project/>')
 else:
  (root/'Sample.csproj').write_text('<Project/>');(root/'ignored').mkdir(exist_ok=True);(root/'ignored/foreign.csproj').write_text('x'*(4*1024*1024+1))
 args=[binary,'scan','--root',root,'--json'];
 if profile=='excluded':args+=['--exclude','ignored']
 report=json.loads(run(profile,'scan',args).stdout)
 if profile=='multiple':
  assert len(report['suggestions'])==4 and all(s['command']['args'][2].startswith('./')for s in report['suggestions']);run(profile,'explicit-selection',[binary,'init','--root',root,'--format','json','--select',report['suggestions'][0]['service']+':'+report['suggestions'][0]['name']]);cfg=json.loads((root/'mudarro.json').read_text());assert len(cfg['services'][0]['commands'])==1
 elif profile=='invalid':assert all(s.get('language')!='csharp' and not s.get('commands') for s in report['config'].get('services',[])) and report.get('warnings')
 else:assert len(report['config']['services'])==1 and not report.get('warnings')
packages=list((base/'cache').rglob('*.nupkg'));assert not packages
summary={'result':'PASS','commands':len(results),'sdk':'10.0.401','runtime':'net10.0 console compile and explicit stdlib self-test','mudarroBinary':str(binary),'mudarroBinarySHA256':hashlib.sha256(binary.read_bytes()).hexdigest(),'profiles':['console','multiple-static-choice','invalid-static','excluded-static'],'feedsCleared':True,'externalApplicationPackages':0,'nupkgFilesInPrivateCaches':0,'nativeTestFrameworkExecuted':False,'dotnetTestExecuted':False,'workloadsCertificatesGlobalToolsInstalled':False,'persistentServicesStarted':False,'generateIdempotent':True,'manualConfigAndSourcesPreserved':True,'sourceHashes':before,'generatedHashes':generatedbefore,'controlledFailureExit':{'directConsole':7,'mudarroRun':1,'wrapper':1},'controlledCompileFailure':True,'sourceRestoredAndRecoveryPassed':True,'inheritedHomePreserved':True,'fixture':str(fixture),'collectorInitialExpectationCorrected':'Invalid XML yields warning plus generic custom fallback, no C# commands'};(base/'evidence/summary.json').write_text(json.dumps(summary,indent=2)+'\n');print(json.dumps(summary,indent=2))
