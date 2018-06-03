We check if arrays are empty to suppress warnings/notices.

{foreach array("GET", "POST") as $method}
HTTP {method}
{/foreach}

{foreach (array("GET", "POST") as $method)}
HTTP {method}
{/foreach}

{foreach $methods as $method}
HTTP {method}
{/foreach}
