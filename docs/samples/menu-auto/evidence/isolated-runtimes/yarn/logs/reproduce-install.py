import pathlib,os,subprocess,json,urllib.request,hashlib,base64,tarfile
root=pathlib.Path('/tmp/mudarro-runtimes.K7F4J1');cfg=root/'config';cfg.mkdir(exist_ok=True)
for name in ['yarn-npm-user.rc','yarn-npm-global.rc']:(cfg/name).write_text('');(cfg/name).chmod(0o600)
for rel in ['yarn-classic','yarn-modern','cache/npm-classic','cache/corepack','logs/yarn','fixtures/yarn']:(root/rel).mkdir(parents=True,exist_ok=True)
env={k:v for k,v in os.environ.items() if not k.lower().startswith('npm_config_') and k not in ['NPM_TOKEN','NODE_AUTH_TOKEN']}
node='/home/danielsouza/.nvm/versions/node/v24.15.0/bin/node';cp='/home/danielsouza/.nvm/versions/node/v24.15.0/lib/node_modules/corepack/dist/corepack.js'
versions={}
for name,pkg,version in [('classic','yarn','1.22.22'),('modern','@yarnpkg/cli-dist','4.18.1')]:
 url='https://registry.npmjs.org/'+urllib.parse.quote(pkg,safe='')+'/'+version
 metadata=json.load(urllib.request.urlopen(url));(root/'logs/yarn'/f'{name}-metadata.json').write_text(json.dumps(metadata,indent=2)+'\n')
 tarball=urllib.request.urlopen(metadata['dist']['tarball']).read();actual='sha512-'+base64.b64encode(hashlib.sha512(tarball).digest()).decode();assert actual==metadata['dist']['integrity'],name
 path=root/('yarn-'+name)/(f'{name}-{version}.tgz');path.write_bytes(tarball)
 versions[name]={'version':version,'metadata_url':url,'tarball_url':metadata['dist']['tarball'],'integrity_sha512':actual,'tarball_sha256':hashlib.sha256(tarball).hexdigest()}
 if name=='modern':
  with tarfile.open(path) as archive:body=archive.extractfile('package/bin/yarn.js').read()
  (root/'yarn-modern/verified-official-yarn.js').write_bytes(body);versions[name]['official_runtime_sha256']=hashlib.sha256(body).hexdigest()
args=['npm','--userconfig',str(cfg/'yarn-npm-user.rc'),'--globalconfig',str(cfg/'yarn-npm-global.rc'),'--cache',str(root/'cache/npm-classic'),'--registry','https://registry.npmjs.org','install','--prefix',str(root/'yarn-classic'),'--ignore-scripts','--no-audit','--no-fund','yarn@1.22.22']
with (root/'logs/yarn/classic-install.log').open('w') as log:subprocess.run(args,cwd=root/'yarn-classic',env=env,stdout=log,stderr=subprocess.STDOUT,check=True)
lock=json.loads((root/'yarn-classic/package-lock.json').read_text());assert lock['packages']['node_modules/yarn']['integrity']==versions['classic']['integrity_sha512']
env.update(COREPACK_HOME=str(root/'cache/corepack'),COREPACK_ENABLE_PROJECT_SPEC='0')
with (root/'logs/yarn/modern-install.log').open('w') as log:subprocess.run([node,cp,'yarn@4.18.1','--version'],cwd=root/'yarn-modern',env=env,stdout=log,stderr=subprocess.STDOUT,check=True)
cache_files=list((root/'cache/corepack').glob('v1/yarn/4.18.1/**/yarn.js'));assert cache_files, 'Modern cached runtime missing'
assert any(hashlib.sha256(p.read_bytes()).hexdigest()==versions['modern']['official_runtime_sha256'] for p in cache_files), 'Corepack runtime differs from integrity verified npm artifact'
shim=root/'yarn-modern/bin';shim.mkdir(exist_ok=True);(shim/'yarn').write_text('#!/bin/sh\nexport COREPACK_HOME='+str(root/'cache/corepack')+'\nexport COREPACK_ENABLE_PROJECT_SPEC=0\nexport COREPACK_ENABLE_NETWORK=0\nexec '+node+' '+cp+' yarn@4.18.1 "$@"\n');(shim/'yarn').chmod(0o755)
for label,exe in [('classic',str(root/'yarn-classic/node_modules/.bin/yarn')),('modern',str(shim/'yarn'))]:
 version=subprocess.check_output([exe,'--version'],env=env,text=True,cwd=root/('yarn-'+label)).strip();assert version==versions[label]['version'];versions[label]['runtime_path']=exe;versions[label]['reported_version']=version
(root/'logs/yarn/install-summary.json').write_text(json.dumps(versions,indent=2)+'\n');print('Classic1.22.22 and Modern4.18.1 installed privately; official integrity verified; Corepack runtime bytes matched verified tarball')
