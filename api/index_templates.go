package api

import (
	"html/template"
)

var indexTemplate = template.Must(template.New("index").Parse(`
<html>
	<head>
		<meta charset="utf-8">
		<title>Megam Vertice {{.version}} </title>
		<style>body {font-family: Helvetica, Arial;}</style>
	</head>
	<body>
		<h1>Your are in Megam Vertice! {{.version}} </h1>
		<p>vertice is an omni scheduler for Megam Vertice that aims to make it easier to deploy vms, microservices, unikernel in production.</p>
		<h2>Build and deploy your application</h2>
		<p>Now you're ready to deploy an application to this vertice engine, please refer to the megam documentation for more details: <a href="http://docs.megam.io" title="Deploying an application in Vertice">docs.megam.io</a>.</p>
	</body>
</html>
`))
