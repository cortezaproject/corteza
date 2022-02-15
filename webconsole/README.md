# Corteza Server Web Console

When enabled (`HTTP_WEB_CONSOLE_ENABLED=true`), it allows insight and management of corteza internals. 

## Web Console development

When developing web-console backend (API) base URL must be set to the actual server. That can be achieved by setting a URL as value local store item with `console-api-base-url` as key:

```javascript
localStorage.setItem('console-api-base-url', '//localhost:3000/!console')
```  
