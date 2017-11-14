<html>
<head>
 <title> isula composer service </title>
 <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
</head>
<pre>
	{{range $key, $value := .files}}
	<a href="{{$value}}">{{$value}}</a>
	{{end}}
</pre>
</html>
