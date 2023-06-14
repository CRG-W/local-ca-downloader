# Local CA Downloader
Create a local CA and generate self-signed certs to use for your local home lab
services, based on configuration.

## Configuration
When building and running the image for the first time, or to regenerate certs locally,
you'll need to configure a handful of settings:

Export a passphrase to be used during CA and cert generation. Note that this can
also be added to your `.env` like the remaining config below, but we recommend
you generate a random passphrase in your environment so that it's ephemeral and
not stored in plain text.
```
$ export CERT_PASSPHRASE=<passphrase-for-cert-gen>
```

Configure cert-specific settings by copying `.env.example` to `.env` locally. Here are the
configurable options:

| Environment Variable | Description                                                                               |
| -------------------- | ----------------------------------------------------------------------------------------- |
| `CA_TTL_DAYS`        | the duration in days you wish your CA cert to be valid                                    |
| `CERT_TTL_DAYS`      | the duration in days you wish your cert to be valid                                       |
| `CA_SUBJECT`         | the CA cert subject, formatted as: `/C=US/ST=State/L=City/O=Organization/CN=CA Name`      |
| `CERT_CN`            | the CN for the cert, from your CA subject above, formatted as: `/CN=CA Name`              |
| `CERT_ALT_NAMES`     | the DNS/IPs you wish the cert to be valid for, formatted as: `DNS:localhost,IP:127.0.0.1` |

## Build and Run
The first time you build/run the image, certs will be generated locally and placed
in your `/certs` directory. Subsequent builds will not overwrite these certs. If you
wish to generate new certs, delete the cert files in your `/certs` directory and
rebuild.

Run the following to generate certs (when none found) and start the web app, which
will be served on `http://localhost:8081`
```
$ docker-compose up
```

You can now download your certs to use locally!
