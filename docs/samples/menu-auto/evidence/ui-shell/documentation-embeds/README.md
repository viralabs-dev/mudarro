# Visible documentation media — MUD-054

Nine pages, 36 standalone Markdown image embeds, zero broken relative links. Each gallery uses #recorded-screens (English) or #telas-gravadas (pt-BR). Root README shows dark-English; pt-BR README shows light-Portuguese. Guides/samples show both. Existing recordings and Library version 3 are unchanged.

Images use standard Markdown image syntax with relative PNG/GIF paths, filenames match on a case-sensitive filesystem, and each image has a visible prose caption. PNGs were decoded and visually inspected. GIFs decode all frames, have nonzero durations/multiple frames and loop indefinitely. This validates the local media and embed structure; it does not claim a browser or actual GitHub page preview was opened.

Read-only remote HEAD query confirmed main at 9869dc6302d4e131a4033e5a080f93fcb5c7449b. That tree contains earlier UI-configuration media but neither the new ui-shell assets nor these README embeds. Publication of MUD-032 and MUD-054 was subsequently authorized on 2026-10-06. This revision includes the embeds and assets with `[skip ci]`; the earlier remote query is retained as the pre-publication observation. No settings change or CI dispatch/rerun is authorized.

Optional validation through GitHub POST /markdown was not executed: automatic approval review rejected sending full local page text to the external API without explicit disclosure authorization. No retry/workaround sent that content. Local link/media validation completed. No browser, desktop capture, installation or code changes in this task.
