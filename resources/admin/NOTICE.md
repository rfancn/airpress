# Admin console — license notice

Everything in this directory is the compiled **AirPress admin console** (front-end).
It is **not** covered by the MIT license that applies to the rest of this repository.

## Origin

The console is derived from [halo-dev/console](https://github.com/halo-dev/console) (GPL-3.0),
through the go-sonic fork [go-sonic/console](https://github.com/go-sonic/console) (GPL-3.0).

The build artifacts come from [go-sonic/sonic](https://github.com/go-sonic/sonic) and,
apart from the rebranding listed below, are byte for byte identical to upstream.

## Local changes

All changes are string substitutions only — no logic was modified.

- `index.html` — page title, `<meta name="generator">` and noscript text renamed to AirPress.
- `js/app.fddf1567.js` — `alt:"Sonic Logo"` → `alt:"AirPress Logo"`;
  footer `Powered by Sonic` → `Powered by AirPress`;
  page-title suffix `De="Sonic"` → `De="AirPress"`.
- `js/328.65547e3f.js` — `alt:"Sonic Logo"` → `alt:"AirPress Logo"`;
  install wizard alert `欢迎使用 Sonic` → `欢迎使用 AirPress`;
  import dialog `是否为 Sonic 后台导出的文件` → `是否为 AirPress 后台导出的文件`.
- `js/40.4ad3ba1c.js`, `js/690.df587247.js` — `alt:"Sonic Logo"` → `alt:"AirPress Logo"`.
- `js/892.b26d4542.js` — About page: `感谢Halo，Sonic 的前端项目 Fork自Halo.` →
  `感谢Halo，AirPress 的前端项目 Fork自Halo.` (the Halo credit itself is kept);
  `检测到 Sonic 新版本` → `检测到 AirPress 新版本`;
  project link, contributors and release check point at `rfancn/airpress`
  instead of `go-sonic/sonic`.
- `js/123.bbbf91d4.js` — theme-ecology link → `github.com/rfancn/airpress#theme-ecology`.

Every patched `js/*.js` has its `js/*.js.gz` companion regenerated from the patched file.

The `https://github.com/go-sonic` link on the About page is kept on purpose, as credit to
the originating organization. The logo images are still the upstream artwork (see the
credits in the repository README).

## License

This directory is distributed under the **GNU General Public License v3.0**.
The full license text is in [LICENSE](LICENSE).
There is no warranty, to the extent permitted by law.

## Corresponding source

As GPL-3.0 requires, the corresponding source code of this front-end is available:

- Console source: <https://github.com/go-sonic/console> (GPL-3.0), itself a fork of
  <https://github.com/halo-dev/console> (GPL-3.0)
- The AirPress modifications are the string substitutions listed under
  [Local changes](#local-changes) above; applying them to the upstream source reproduces
  this build, and `index.html` is included here in its modified form.
