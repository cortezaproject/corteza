{foreach $apis as $api}
# {api.title}

{api.description}

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
{foreach $api.apis as $name => $call}
| `{call.method}` | `{api.path}{call.path}` | {call.title} |
{/foreach}


{foreach $api.apis as $name => $call}
## {call.title}

{call.description}

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `{api.path}{call.path}` | {if $api.protocol}{api.protocol}{else}HTTP/S{/if} | {call.method} | {eval echo implode(", ", $api.authentication)} |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
{foreach $call.parameters as $method => $params}
{foreach $params as $param}
| {param.name} | {param.type} | {method|toupper} | {param.title} | {if empty($param.default)}N/A{else}{param.default}{/if}
{if $param.values}<br><br>Values:<br><br><ul>{foreach $param.values as $value}<li>`{value}`</li>{/foreach}{/if} | {if $param.required}YES{else}NO{/if} |
{/foreach}
{/foreach}


{/foreach}

---

{/foreach}