# MUD-055 — current example mudarro

Validation captured locally on 246048a before publication. Publication of MUD-055 was explicitly authorized; this revision includes the refreshed example and uses [skip ci], without claiming a new CI pass. Current presentation slogans removed while offline behavior is unchanged. Config.Name remains dynamic; the recording fixture uses project/service mudarro. Old Aurora samples/casts/PNGs/GIFs/metrics remain untouched historical evidence.

86 top-level tests passed with race and optional recording, zero failure/skip. Coverage OR-merged by source block: 1,866/2,477 (75.33%). Vet/Linux build and Darwin binary/terminal-test crosscompilation passed; no native macOS runtime claim. Dark/en and light/pt-BR real PTY captures each passed 24 chrome/header checks, verified project mudarro, no removed presentation slogans, and termios restoration. Resize/NO_COLOR/SIGTERM checks passed.

Fresh GIFs: dark 983×739/55 frames, light 983×739/53 frames, both 12.31s looping. All frames decoded; composed PNG frames at 4.2s visually inspected. PNGs are recording frames, not desktop screenshots. No GUI/browser/global terminal setting changed.

Nine galleries contain 40 image embeds, referencing current mudarro-shell assets and distinctly historical UI-configuration preview images. Files match case-sensitive relative paths. Library uses the same identities, version 4; previous versions retained.

Reproduction: build the current binary, then run scripts/record-ui-shell.py BINARY OUTPUT and scripts/validate-ui-shell-resize.py BINARY RESIZE_OUTPUT. Both helpers own disposable fixtures, check dynamic project headers and exercise actual CLI output.
