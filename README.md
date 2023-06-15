<img src="images/lcm_logo_mint.png"  width="50%" height="20%">

---
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
The first time you build/run the image, you'll need to run the container with `GENERATE_CERTS=true`
With this, certs will be generated locally and placed in your `/certs` directory.
These certs will also be used to serve the service over SSL.

Subsequent builds will not need to generate certs, so you should _not_ run with this
environment variable set unless you wish to wipe your old certs/regenerate if and
when certs expire.

### Auth
You'll also need to set `AUTH_PASSWORD=<web-app-password>` in order to gain access
to the web service in the browser. You can either export this `export AUTH_PASSWORD=<web-app-password>`
in the shell you'll be running `docker-compose up` from, or supply it in the compose
command:
```
$ export AUTH_PASSWORD=<web-app-password> # or injected into your env via .env, or some other way
$ ...
$ ...
$ docker-compose up
```
or
```
$ AUTH_PASSWORD=<web-app-password> docker-compose up
```

### First Run / Rebuild Certs
With your `AUTH_PASSWORD` set, run the following to build/run the image and (re)generate certs:
```
$ GENERATE_CERTS=true docker-compose up
```

### Subsequent Runs / No Cert Rebuild
After certs are generated, you can simply run the following when you want to start the service:
```
$ docker-compose up
```

### Access the UI/Enable SSL
You can now access/download your certs to use locally by visiting `https://localhost:8443/`

Download your public CA and add it to your local trust to get rid of the "not secure"
browser warning, as the service will start over a TLS connection, consuming the
certs you just generated!
