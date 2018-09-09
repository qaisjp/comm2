# community

Build using `go build`.

Copy `config.yaml.example` to `config.yaml`.

Start with `config=config.yaml ./community`

## OAuth

- Send the user to `https://forum.mtasa.com/oauth/authorize/?client_id={CLIENT_ID}&response_type=code&redirect_uri=http://localhost:8080`