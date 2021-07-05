package tpl

const defaultHTML = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta http-equiv="x-dns-prefetch-control" content="on" />
<link rel="dns-prefetch" href="//cdn.jsdelivr.net" />
<meta name="viewport" content="width=device-width,minimum-scale=1,initial-scale=1,maximum-scale=5,viewport-fit=cover">
<title>Notes</title>
<meta name="robots" content="noindex, nofollow">
<link rel="apple-touch-icon" sizes="180x180" href="https://cdn.jsdelivr.net/gh/x2ox/memo@f535d41e57af7da22a96d3fbab8bd8f3bddab3c1/.data/static/apple-touch-icon.png">
<link rel="icon" type="image/png" sizes="32x32" href="https://cdn.jsdelivr.net/gh/x2ox/memo@f535d41e57af7da22a96d3fbab8bd8f3bddab3c1/.data/static/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="16x16" href="https://cdn.jsdelivr.net/gh/x2ox/memo@f535d41e57af7da22a96d3fbab8bd8f3bddab3c1/.data/static/favicon-16x16.png">
<link rel="manifest" href="https://cdn.jsdelivr.net/gh/x2ox/memo@f535d41e57af7da22a96d3fbab8bd8f3bddab3c1/.data/static/site.webmanifest">
<link rel="mask-icon" href="https://cdn.jsdelivr.net/gh/x2ox/memo@f535d41e57af7da22a96d3fbab8bd8f3bddab3c1/.data/static/safari-pinned-tab.svg" color="#ffffff">
<meta name="msapplication-TileColor" content="#ffffff">
<meta name="theme-color" content="#ffffff">
<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/x2ox/memo@d6e74cc79f80ce7da9977df834f593346806f964/.data/static/markdown.css"/>
<style>
#preview-box {
	background-color: rgba(1,1,1,0.1);
	padding: 1px 10px;
	border-radius: 10px;
}
img {
    width: 100%;
    height: 100%;
}
</style>

</head>
<body>

{{ . }}

</body>
</html>`
