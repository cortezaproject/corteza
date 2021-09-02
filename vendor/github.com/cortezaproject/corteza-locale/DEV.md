# Guide for frontend and backend developers

This guide is for corteza frontend and backend developers.
It shows how to connect translations from corteza-locale repository to frontend web applications and backend server.

## Prerequisites

Clone corteza-locale repository to a separate folder.

## Fronted developers

### Using corteza-server Docker container to support webapp development

#### Configuration with Docker Compose 

See the `docker-compose.yaml` in the root of the repository.

1. Run it by executing (inside cloned `corteza-locale` repository):
```shell
docker-compose up -d 
```

2. Fix `config.js` in the web application to point to the server
```js
window.CortezaAPI = `//localhost:1818/api`
```

#### Verify the loaded languages

```shell
docker-compose logs | grep locale | head
```

This will show you first couple (head) filtered (grep) log lines.
If some of them contain "language loaded" that reflect the setup in your `corteza-locale/src` you have successfully loaded translations into Corteza server.  

#### Verify by loading translations
```shell
curl 'http://localhost:1818/api/system/locale/en/corteza-webapp-admin' -H "Accept: application/json" -H 'Accept-Language: en'
```




## Backend developers

Corteza server loads, parses and serves all languages files for all frontend web applications and server. 

### Changes in configuration

The following chapter assumes your Corteza server development env is already set-up

Add `LOCALE_PATH` and `LOCALE_LOG` to your .env file
```dotenv
LOCALE_PATH=../corteza-locale/src
LOCALE_LOG=true 
```

Path can be absolute or relative and should contain subdirectories with languages.
You can remove or comment-out LOCALE_LOG if you find the setting to verbose.  


#### Verify the loaded languages

When you (re)start your corteza server yu should see the following log among the first logged lines:

```
20:24:42.881	INFO	locale	locale/locale.go:81	reloading	{"src": ["../corteza-locale/src"]}
20:24:42.899	INFO	locale	locale/load.go:66	language loaded	{"tag": "en", "config": "../corteza-locale/src/en/config.yaml"}
```

You will see additional log lines for every language loaded.

