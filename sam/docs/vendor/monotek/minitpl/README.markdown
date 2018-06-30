# MiniTPL

The goal of the MiniTPL template engine is to provide a miniature
framework which allows you to rapidly create and consume
Smarty-like templates without adding the overhead of Smarty to
your choice of a PHP framework.

In benchmarks the speed of Mini TPL is very close to PHP itself.
All that is usually needed for Mini TPL is a 3KB PHP code overhead.
So it beats Smarty, and usual PHP vsprintf and str_replace functionality.

With a total size of about 13KB and the functionality contained, this is
one of the smallest full featured template engines for PHP to date.

MiniTPL is available on [packagist as monotek/minitpl](https://packagist.org/packages/monotek/minitpl).

To start using MiniTPL in your project with [composer](http://getcomposer.org/), create a composer.json file:
```
{
    "require": {
        "monotek/minitpl": ">=1.0"
    }
}
```

And run `composer install`. You can start using MiniTPL right away

```
<?php

include("vendor/autoload.php");

$tpl = new Monotek\MiniTPL\Template;

$tpl->load("test.tpl");
$tpl->render();
```