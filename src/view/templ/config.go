package templ

const Config = `
	app:
	  name: "{{.ProjectName}}"
	  host: "localhost:9003"
	  service: "content"
	  env: "local"
	  port: "9003"
`
