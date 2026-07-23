# rebuild (run from the repo root)
 go build -o /opt/zimmer/zimmer .
 cp -r ./app /opt/zimmer/app
 cp .env /opt/zimmer/.env

 # re‑apply the bind‑capability if you rebuilt
 sudo setcap 'cap_net_bind_service=+ep' /opt/zimmer/zimmer

 # reload systemd and restart
 sudo systemctl daemon-reload
 sudo systemctl restart zimmer.service

