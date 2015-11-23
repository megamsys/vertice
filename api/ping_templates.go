package api

import (
	"html/template"
)

var pingTemplate = template.Must(template.New("ping").Parse(`
  <html>
      <head>
          <meta charset="utf-8">
          <title>Oja ping</title>
					<style>
            * { font-family: sans-serif; }
            h1 { font-size: 20px; }
            h2 { font-size: 16px; }
            p { font-size: 15px; }
						textarea
						{
  					border:1px solid #999999;
  					width:100%;
  					margin:5px 0;
  					padding:3px;
						}
        </style>
	 </head>
      <body>
          <h1>Oja ping</h1>
					<textarea rows="100" cols="2">
					{{ .data}}
					</textarea>
			</body>
  </html>
`))
