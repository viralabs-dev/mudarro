# MUD045 preparation contract

Existing tools only: GNU Make4.3, GCC/G++13.3.0, all /usr/bin. CMake and Meson are absent and outside this phase. No installation, network, production, repository, Vault or Git changes.

Two independent zero-dependency fixtures: C program prints21; C++ program prints42. Makefile explicit build/test/fail targets. Build compiles with strict warnings; test runs binary and checks exact output; repeated build keeps output hash/mtime; intentional failure recipe exits7 but GNU Make reports2; missing target reports2. Source/Makefile hashes stay unchanged throughout. No startup inference/service supervisor used.

Current Mudarro scan's existing Make suggestions can inventory explicit Make targets and use custom service fallback. This is evidence of Make integration only, not an implemented C/C++ adapter. Manual custom configuration declares build/test/fail with `make <target>`; real generated wrappers can demonstrate those explicit commands. Existing manual YAML must survive repeated scan/init rejection/generate twice.

Proposed future offline detection: bounded regular contained source/Makefile reads; .c C evidence, .cpp/.cc/.cxx C++ evidence, explicit Make targets as opt-in suggestions only; no compiler, Make, CMake or Meson execution in scan; no generated install/start/framework/DB guesses. Make-only declaration remains generic custom evidence. Mixed C/C++ sources and source trees without Makefile require explicit contract choice before production.

Reusable source allowlist: package-independent `main.c` or `main.cpp`, Makefile, manual mudarro.yaml, instructions and source hash manifest. Exclude compiled build/, generated outputs, .mudarro state/wrappers, caches, installed tools and binaries. Runtime logs separately archived.
