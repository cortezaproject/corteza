Monotek Mini Template
=====================

Zadnje posodobljeno - Tue 18 Dec 2012 11:24:30 AM CET

Avtor: Tit Petrič ( tit.petric@monotek.net ) / Monotek d.o.o.<br/>


Predstavitev
------------

Stabilna Produkcijska različica podpira sestavljanje šablone in je združljiva
s PHP4 in PHP5. Za ohranjanje majhne velikosti kode uporablja zelo napredne
funkcije, ob tem pa poskrbi za izjemno uporabnost.

Celotna koda je velika 13,5 KB in vsebuje minimalno količino komentarjev in
zamik kode z uporabo tabulatorjev za berljivost.

Koda je zaščitena z [Creative Commons Priznanje Avtorstva - Deljenje pod
enakimi pogoji](http://creativecommons.org/licenses/by-sa/3.0/) licenco.


Namestitev in potrebna programska oprema
----------------------------

Za delovanje Monotek Mini Template potrebujete PHP različice 4.3.0 ali novejšo.
Sistem lahko deluje tudi s starejšimi PHP različicami, če priskrbite svojo
[file_get_contents()}(http://php.net/file_get_contents) funkcijo.

Glede na mapo, kjer uporabljate template objekt, potrebujete mape:

> templates/<br/>
> templates/cache/ <span class="red">*</span>

Ker pa se ta template sistem sestavlja, boste potrebovali tudi cache mapo.

Jezikovna vezava
------------------

Ker je način uporabe template sistema podoben Smarty template sistemu in PHP
programskem jeziku, vam bo poznavanje PHP osnov v pomoč pri uporabi tega
template sistema.

1. <a href="#variables">Sestava jezika</a>
2. <a href="#loops">Zanke</a>
3. <a href="#conditions">Pogoji</a>
4. <a href="#blocks">Block in Inline definicije</a>
5. <a href="#advanced">Napredno, vstavljanje PHP kode</a>
6. <a href="#php">Uporaba template sistema v PHP</a>

----------------------------------------------------------

<h3 id="variables">1. Jezikovna sestava <a href="#top">Δ Skok na vrh</a></h3>

Jezikovna sestava vam ponuja razlago delovanja sistema, da boste lahko pričeli
z ustvarjanjem in uporabo svojih template datotek. Prva stvar, katero je
potrebno povedat je, da so vse jezikovne sestave definirane med zavitima
oklepajema (`{` in `}`). Če template sistem ne prepozna jezikovne sestave, je
ne prevaja. To omogoča uporabo javascript in json zapisov brez težav.

#### 1.1. Spremenljivke

Uporaba spremenljivk je enostavna. Spremenljivke so zaprte med zavite oklepaje:
`{spremenljivka}`. Ker so spremenljivke prevedene v zadnji stopnji, je uporaba
znaka dolar pred njimi neobvezna. Zato ni nobene razlike med zapisom
`{$spremenljivka}` in zapisom `{spremenljivka}`, oba sta pravilna.

~~~~~~~~~~~~~
{spremenljivka} je enako kot {$spremenljivka}
~~~~~~~~~~~~~

#### 1.2. Množice

Obstaja kratek način uporabe množic znotraj template. Če želite izpisati
spremenljivko, ki se nahaja znotraj množice, za naslavljanje uporabite operator
`.` (pika). Množice lahko naslavljate v poljubno globino kot `$mnozica.tocka.0`
ali pa tudi z uporabo spremenljivk, ki se začnejo z znakom dolar
`{$mnozica.$novice.naslov}`. Prav tako je mogoča uporaba PHP jezikovne sestave
kot je razvidno iz spodnjega primera.

~~~~~~~~~~~~
{$array.items} je enako kot {$array['items']}
{$array.items.0} je enako kot {$array['items'][0]}
{$array.$items.0} je enako kot {$array[$items][0]}
~~~~~~~~~~~~

#### 1.3. Prikrojevalci

Prikrojevalci so PHP funkcije, ki iz poljubne vrednosti (ponavadi besede)
vrnejo besedo, ki se prikaže. Nekatere PHP funkcije se lahko uporabljajo
kot prikrojevalci (strtoupper, ucfirst, str_rot13, strrev, count, ...).

~~~~~~~~~~~~~~~~~~~~~~~~
/** Primer uporabe prikrojevalcev */

function add_it_up($array)
{
	$size = 0;
	foreach ($array as $value) {
		$size += $value['size'];
	}
	return $size;
}
~~~~~~~~~~~~~~~~~~~~~~~~

Za uporabo zgornje funkcije kot prikrojevalca v template datoteko vpišite
`{$variable|add_it_up}`. Prikrojevalec bo šel čez vse elemente spremenljivke
in seštel vrednosti pod ključem `size`, ter izpisal seštevek.

Podajanje dodatnih parametrov prikrojevalcem ni omogočeno. Če ne morete okoli
tega, potem si oglejte kategorijo <a href="#advanced">naprednih zmožnosti</a>
kjer boste zvedeli za način, kako vključiti PHP kodo v template datoteko.

Template sistem vključuje poseben prikrojevalec z imenom `escape`, kater se
zamenja z klicem `htmlspecialchars($left, ENT_QUOTES);`. Ta prikrojevalec vam
služi za izpis podatkov znotraj HTML, kateri mogoče vsebujejo narekovaje ali
pa znake `<` (manjše) in `>` (večje). V primeru je demonstrirana pravilna
uporaba uporaba tega prikrojevalca.

~~~~~~~~~~~~~~~
<input type="title" type="text" value="{title|escape}"/>
<textarea name="content">{content|escape}</textarea>
<a href="{news.link}" title="{news.title|escape}">Read more ...</a>
<h3>{site.title|escape}</h3>
{content} {* pričakujemo HTML, tukaj ne rabimo escape *}
~~~~~~~~~~~~~~~

Dodatna template prikrojevalca sta še `toupper` za `strtoupper` in `tolower`
za `strtolower`. Za naštete prikrojevalce ni narejenih posebnih funkcij.

~~~~~~~~~~~~~~
{variable|escape}
{variable|toupper}
{variable|add_id_up}
~~~~~~~~~~~~~~

Zgornja koda se prevede v naslednje:

~~~~~~~~~~~~~~
<?php
	$_v = &$this->vars;
	echo htmlspecialchars($_v['variable'], ENT_QUOTES);
	echo strtoupper($_v['variable']);
	echo add_id_up($_v['variable']);
?>
~~~~~~~~~~~~~~

V praksi se prikrojevalci pogosto uporabljajo za izpis podatkov v json notaciji
za uporabo znotraj javascript knjižnic (uporaba funkcije `json_encode`)

#### 1.4. Objekti

Objekte lahko kličete tudi kot prikrojevalce, za globalne ali lokalne objekte.
Ob prevajanju se določi, če obstaja globalni objekt z navedenim imenom in
se glede na to prilagodi prevedeno kodo.

~~~~~~~~~~~
{$variable|$memcache->get}
~~~~~~~~~~~

Če globalni objekt ne obstaja se predvideva, da obstaja lokalni objekt.

~~~~~~~~~~~
<?php
	$_v = &$this->vars;
	echo $_v['memcache']->get($_v['variable']);
?>
~~~~~~~~~~~

Če globalni objekt v času prevajanja obstaja se template prevede tako:

~~~~~~~~~~
<?php
	$_v = &$this->vars;
	global $memcache;
	echo $memcache->get($_v['variable']);
?>
~~~~~~~~~~

Na enak način lahko uporabljate tudi spremenljivke iz objektov.

~~~~~~~~~~~~
{$memcache->variable}
~~~~~~~~~~~~

Kar se prevede se v naslednje:

~~~~~~~~~~
<?php
	$_v = &$this->vars;
	echo $_v['memcache']->variable;
?>
~~~~~~~~~~

ali

~~~~~~~~~~
<?php
	$_v = &$this->vars;
	global $memcache;
	echo $memcache->variable;
?>
~~~~~~~~~~

#### 1.5. Konstante

Vrednost, ki ima na začetku znak `_` se smatra za konstanto. Template sistem
bo izpisal vrednost konstante, če je ta določena, ali pa ime konstante če ni.
Zaradi teh pravil bo niz `{_MOJA_KONSTANTA}` uporabljen kot konstanta, medtem
ko bo `{$_MOJA_NESPREMENLJIVKA}` uporabljen kot spremenljivka, ker se začne z
znakom za dolar `$`.

~~~~~~~~~~~
{_MY_CONSTANT}
{_this_is_also_a_constant}
{$_my_variable}
~~~~~~~~~~~

Zgornja koda se prevede v naslednje:

~~~~~~~~~~~
<?php
	$_v = &$this->vars;
	echo _MY_CONSTANT;
	echo _this_is_also_a_constant;
	echo $_v['_my_variable'];
?>
~~~~~~~~~~~

#### 1.6. Vključitev ločenih template datotek

Niz `{include imedatoteke.tpl}` se zamenja z vsebino določene template
datoteke ob začetku prevajanja. V primeru, da želite spremeniti vsebino
vključene template datoteke, morate obnoviti prevedeno kodo glavne template
datoteke. To naredite ročno (z izbrisom vsebine cache direktorija ali obnovo
časa posodobitve glavne template datoteke). V template datotekah lahko
vključite poljubno število drugih template datotek.

~~~~~~~~~~~
<html>
<head>
	<title>{$title}</title>
</head>
<body>
{include site_header.tpl}
{$contents}
{include site_footer.tpl}
</body>
</html>
~~~~~~~~~~~

Če želite imeti dejansko nalaganje skupnih datotek, kar reši nekatere težave
okoli dinamičnega nalaganja in hranjenja cache datotek, potem uporabite `{load}`:

Imejte v mislih, da s to metodo ni možna uporaba `inline` in `blok` definicij
izven naložene template datoteke. To lahko dosežete z uporabo `include`.

~~~~~~~~~~~
<html>
<head>
	<title>{$title}</title>
</head>
<body>
{load $dynamic_header}
{load $dynamic_contents_template}
{load $dynamic_footer}
</body>
</html>
~~~~~~~~~~~

#### 1.7. Komentarji

Komentarji so zaprti med `{*` in `*}`. Komentarji se odstranijo med prevajalnim časom in se ne prikažejo v prevedeni kodi ali izpisu. Uporabni so za izčrpno dokumentacijo, s čimer ne povzročajo povečanega časa obdelave.

~~~~~~~~~~~
{* This is a comment that won't be shown anywhere,
   except in the source template, only to developers. *}

Hello world!
~~~~~~~~~~~

Kar se prevede v:

~~~~~~~~~~~
`Hello world!`
~~~~~~~~~~~

<h3 id="loops">2. Zanke <a href="#top">Δ Skok na vrh</a></h3>

#### for, foreach, foreach / else, while

Cilj sintakse je pomagati razvijalcu, zato je podobna PHP sintaksi.
Medtem ko to ni idealno za spletne programerje, ki izdelujejo template
datoteke, je idealno za veliko večino PHP razvijalcev, katerim se ni
potrebno naučiti novega programskega jezika za template sistem.
Pravtako, osnove PHP jezika kot so predstavljene, nebi smele predstavljati
večjih problemov za oblikovalce, ki že zdaj izdelujejo template datoteke
npr. za Smarty. Običajno programer naredi prvo template datoteko, oblikovalci
jo pa za tem popravljajo.

~~~~~~~~~~~~
<html><body>
{for $i=0; $i<count($array); $i++}
	Hop number {i}<br/>
{/for}
{foreach $array as $value}
	I like values like I like: {value}<br/>
{/foreach}
{foreach $array as $key=>$value}
	I like my value {value} to have a key {key}.<br/>
{/foreach}
{foreach $array as $key=>$value}
	I like my value {value} to have a key {key}.<br/>
{else}
	I'm sorry, I have nothing in the array.<br/>
{/foreach}
{while ($k++ < 10)}
	Well, k is {k}<br/>
{/while}
</body></html>
~~~~~~~~~~~~

Kot vidite, je template sistem zelo prijazen PHP razvijalcem.
Primerjava z originalno PHP kodo (dva načina):

~~~~~~~~~~~~
echo '<html><body>';
for ($i=0; $i<count($array); $i++) {
	echo 'Hop number '.$i.'<br/>';
}
foreach ($array as $value) {
	echo 'I like values like I like: '.$value.'<br/>';
}
foreach ($array as $key=>$value) {
	echo 'I like my value '.$value.' to have a key '.$key.'.<br/>';
}
if (!empty($array)) {
	foreach ($array as $key=>$value) {
		echo 'I like my value '.$value.' to have a key '.$key.'.<br/>';
	}
} else {
	echo "I'm sorry, I have nothing in the array.<br/>";
}
while ($k++ < 10) {
	echo "Well, k is ".$k."<br/>";
}
echo '</body></html>';
~~~~~~~~~~~~

Medtem ko ta primer prikazuje podobnost v skladnji, prikazuje tudi, kako vam
pisanje template datotek lahko vzame manj časa in je končni izdelek bolj
berljiv. Če ste morda opazili, izjave podpirajo celotno PHP sintakso.
Medtem ko to ni očitno ko pride do `foreach` zanke, lahko vidite, da je
pri ostalih primerih prevedena koda skoraj identična template kodi.

~~~~~~~~~~~~~~~~~~~
<p>{foreach $items.news as $newsitem}<b>-</b>{/foreach}</p>
<p>{foreach $items['news'] as $newsitem}<b>+</b>{/foreach}</p>
<p>{for $i=0; $i<count($items.news); $i++}<b>!</b>{/for}</p>
<p>{for $i=0; $i<count($items.news)+count($items['news']); $i+=2}<b>?</b>{/for}</p>
~~~~~~~~~~~~~~~~~~~

<h3 id="conditions">3. Pogoji <a href="#top">Δ Skok na vrh</a></h3>

Sintaksa za `if` in `elseif` izjave je enaka kot v PHP. Z njimi lahko kličete
funkcije, metode, uporabljate konstante in izvajate aritmetične operacije.
Spremenljivke so prevedene v lokalne spremenljivke, katere ste določili pri
prikazu template datoteke. Istočasno se lahko uporablja PHP sintakso za
naslavljanje množic kot pa tudi olajšano (shorthand) sintakso z uporabo pike.
V prvem primeru je demonstrirana različna sintaksa za množice.

#### 3.1. if izjave

~~~~~~~~~~~~~
{if $is.admin && $user['name']=="black"}
        {* Hello black! Only you can edit things,
           but only as long as you stay an admin. *}
  ...
{/if}
~~~~~~~~~~~~~

#### 3.2. if/elseif/else izjave

~~~~~~~~~~~~~
{if $is.admin && $user['name']=="black"}
        {* Hello black! Only you can edit things,
           but only as long as you stay an admin. *}
  ...
{elseif $is_moderator}
        {* Hello moderator! You can do some things. *}
{else}
	{* You are a nobody and you earn nothing! *}
{/if}
~~~~~~~~~~~~~

#### 3.3. foreach/else izjave

~~~~~~~~~~~~~
{foreach $newsitems.$section.items as $item}
	<div class="newsitem">
		<h3 class="title">{item.title}</h3>
		<div class="content">{item.content}</div>
	</div>
{else}
	<div class="notice">
	No newsitems exist in the chosen section.
	</div>
{/if}
~~~~~~~~~~~~~

#### 3.4. nocache

MiniTPL template datoteke se prevajajo v PHP datoteke po potrebi
za najboljšo možno hitrost izvajanja. Z uporabo `{include}` direktiv
se osveževanje php datotek rahlo zakomplicira. Za potrebe razvoja smo
dodali `nocache` direktivo. Z uporabo direktive zagotovite, da se vaš
template cache zbriše po vsaki uporabi.

~~~~~~~~~~~~
{*nocache*}
~~~~~~~~~~~~

To funkcionalnost se mora omogočiti z setiranjem `$tpl->_nocache` na `true`.

<h3 id="blocks">4. Bloki in Inline definicije <a href="#top">Δ Skok na vrh</a></h3>

Odvisno od uporabe boste morda želeli večkrat uporabiti delce iste kode v
eni ali večih template datotekah. To lahko dosežete z uporabo `block` in
`inline` definicij. Razlika med `block` definicijo in `inline` definicijo je
v tem, da se `block` definicija prevede v PHP funkcijo in jo lahko uporabljamo
rekurzivno, naprimer za ustvarjanje drevesne strukture. Uporaba `inline`
definicije vam omogoča le večkratno uporabo enake kode.

#### 4.1 block definicija in uporaba

~~~~~~~~~~~
{block recurse}
{if $i++ < 10}
        call {i}
        {block:recurse}
{/if}
{/block}

{block:recurse}
~~~~~~~~~~~

Prevedena template datoteka bo izgledala nekako takole:

~~~~~~~~~~~
<?php

/* This code has been cleaned up some, for
   documentation purposes. It illustrates
   the compile aspect of blocks, but is not
   a carbon copy of the compiled template. */

	$_v = &$this->vars;
	function recurse_1213891842_673($_v) {
		if ($_v['i']++ < 10) {
	        	echo "call ".$_v['i'];
        		recurse_1213891842_673(&$_v);
		}
	}
	recurse_1213891842_673(&$_v);
?>
~~~~~~~~~~~

Kot lahko vidite, se `block` definicija prevede v funkcijo, ki uporablja
iste spremenljivke, katere ste določili template datoteki.

Ustvarjanje ali izris drevesne strukture ta način ni zelo praktičen, je pa
možen. Vrjetno nikoli ne boste potrebovali tega. Po naših izkušnjah je
rekurzija uporabljena le redko, če pa že, pa v PHP kodi in ne v template
datotekah.

Uporabite lahko `block` konstrukt namesto `inline`, če vas skrbi za velikost
prevedene template datoteke. Sam bi to predlagal, ko pridete do tega, da je
ena `inline` definicija dolga več kilobyteov in se jo uporablja večkrat v
isti template datoteki.

#### 4.2 inline definicija in uporaba

Ključna beseda `inline` prihaja iz programskega jezika `C++`, kjer prevajalnik
zamenja klice inline funkcij z vsebino funkcije. S tem v `C++` pridobimo na
hitrosti, ker je klicanje funkcije bolj praktično za berljivost, kot pa
kopiranje kode po najvišjem nivoju. Iz enakega razloga se uporablja tudi
v tem template sistemu.

~~~~~~~~~~
{inline newsitem}
<div class="news">
	<h3 class="title">{news.title}</h3>
	<div class="content">{news.content}</div>
</div>
{/inline}

<div class="top_news">
	{foreach $newsitems.top.items as $news}{inline:newsitem}{/foreach}
</div>

<div class="other_news">
{foreach $newsitems.$section.items as $news}{inline:newsitem}{/foreach}
</div>
~~~~~~~~~~

Uporaba `inline` definicij zmanjša velikost template datotek pred prevajanjem
in pripomore k lepši obliki template datotek. Določeni sestavni deli so lahko
torej večkrat uporabljeni v template datoteki.

<h3 id="advanced">5. Napredno, vstavljanje PHP kode <a href="#top">Δ Skok na vrh</a></h3>

Medtem ko zgornji jezikovni sestavki poskrbijo le za izpis, včasih potrebujemo
tudi posodobiti spremenljivke znotraj template datoteke. Za to poskrbijo
sestavki `eval`, `eval_literal` in `php`. Sestavek `php` je namenjen
zahtevnejšim operacijam in ima kot `eval` dostop direktno do spremenljivk
katere ste podali template datoteki.

#### 5.1 eval

Eval sestavek vam omogoča hitre operacije z lokalnimi template spremenljivkami.
Če želite naprimer narediti tabelo, ki ima v vsaki vrstici drugačen CSS stil,
bi to naredili nekako takole:

~~~~~~~~~~~~~
{eval $style="even";}
<table>
{foreach $rows as $row}
{eval $style = ($style=="even") ? "odd" : "even"}
<tr class="{style}">
<td>{row.message}</td>
</tr>
{/foreach}
</table>
~~~~~~~~~~~~~

Uporabo spremenljivke `$_v` odsvetujemo, saj lahko povzročajo napake
pri prevajanju, ali pa z uporabo te spremenljivke prepišete podatke, kateri
so bili namenjeni prikazu.

Uporaba `{` in `}` znotraj `eval` in `eval_literal` sestavkov ni mogoča.

#### 5.2 eval_literal

Ko potrebujete za izvajanje kode globalne spremenljivke ali objekte,
lahko uporabite `eval_literal` sestavek. Koda znotraj sestavka se ne prevede,
ostane torej taka kot je.

~~~~~~~~~~~
{eval_literal
	global $cms_module;
	$_v['menu_data'] =
		$cms_module->get_menu("branch", array("item","menu")); }
~~~~~~~~~~~

Uporaba `{` in `}` znotraj `eval` in `eval_literal` sestavkov ni mogoča.

#### 5.3 php

~~~~~~~~~~~~~~~
{php}
function mygettime()
{
        return array("time"=>time(),"date"=>date("r"),"microtime"=>microtime());
}
$mygettime = mygettime();
{/php}

{mygettime|var_dump}
~~~~~~~~~~~~~~~

Zgornja koda se prevede v:

~~~~~~~~~
<?php
	$_v = &$this->vars;

	function mygettime()
	{
	        return array("time"=>time(),"date"=>date("r"),"microtime"=>microtime());
	}

	$_v['mygettime'] = mygettime();

	echo var_dump($_v['mygettime']);
?>
~~~~~~~~~

To je najmanj uporabljen in najmanj testiran sestavek. Če želite zgraditi
objekte ali funkcije ali uporabiti veliko PHP kode znotraj tamplate datoteke,
definitivno delate nekaj narobe, čeprav template sistem to podpira.

Uporabo spremenljivke `$_v` odsvetujemo, saj lahko povzročajo napake
pri prevajanju, ali pa z uporabo te spremenljivke prepišete podatke, kateri
so bili namenjeni prikazu.

<h3 id="php">6. Uporaba template sistema v PHP <a href="#top">Δ Skok na vrh</a></h3>

V tej kategoriji bomo pregledali zmogljivosti template sistema v PHP.
Ker že znate izdelovati template datoteke vam bomo v tem oddelku prikazali,
kako narediti PHP kodo s katero uporabimo template datoteke.

#### 6.1 Osnovna uporaba

Za osnovno uporabo morate vključiti datoteko `class.template.php`. Uporabljene
funkcije bodo po potrebi same vključile datoteko `class.template_compiler.php`.
Za razvijalce so na voljo naslednje metode:

`load`, `assign`, `render`, `get`, `set_paths`, `compile`

Ponavadi boste uporabljali le prve štiri metode, če seveda ne želite
spremeniti map, kjer sistem išče template datoteke ali pa bi želeli sami
prevajati template datoteke.

##### 6.1. load ( string $filename )

S to metodo določite katera template datoteka se naj naloži. Metoda bo
preiskala nastavljene poti (Privzeto: `templates/`) in naložila želeno
template datoteko. Klic te metode izprazni podatke namenjene template
datoteki, katerega naknadno določite z uporabo `assign` metode.

##### 6.2. assign ( mixed $key, [ mixed $value = '' ] )

Za programerja je ta metoda najbolj ključna. Rezultat metode je različen glede
na količino in vrstni red podanih parametrov.

Če je prvi parameter množica, drug parameter pa ni podan, potem se template
datoteki določi ena vrednost z ključem vnosa v množico, za vsak vnos.

Če je prvi parameter množica in drug parameter niz znakov, bo narejen vnos
za vsako vrednost v množici. Za ključ vnosa bo uporabljen niz v drugem
parametru, kateremu sledi `_` in potem ključ vnosa iz množice.

Če je prvi parameter niz znakov, potem bo drug parameter določen kot vrednost
dosegljiva pod ključem prvega parametra. Tip drugega parametra ni pomemben.

~~~~~~~~
$data = array(); // some example data
$data['title'] = "Leno promises smooth transition to O'Brien";
$data['content'] = "For months, Fallon has been widely considered
                    the top choice to succeed O’Brien when he steps
                    down next year. On Thursday, published reports ...";

/* This will define an entry {timestamp} */

$tpl->assign("timestamp", time());

/* This will define entries {title} and {content} */

$tpl->assign($data);

/* This example will define {news_title} and {news_content} */

$tpl->assign($data, "news");

/* This will define the item {news}, which contains an array.
   You can output the fields with {news.title} and {news.content} */

$tpl->assign("news", $data);
~~~~~~~~

##### 6.3. render, get

Te metode nimajo parametrov. Prevajanje template datotek je izvedeno v ozadju,
če je potrebno. Metoda `render` izpiše podatke v navadni obliki, medtem ko
jih medota `get` vrne v obliki katero lahko uporabimo v PHP.

##### 6.4. set_paths ( [ $paths = false ] )

Metoda `set_paths` vzame množico z možnimi lokacijami za template datoteke.
Lokacija za prevajanje template datotek je mapa `cache`, katera mora obstajati
v vsaki od podanih map. Podane mape se morajo končati z vrezom (slash).

Template sistem bo pregledal `$paths` množico, dokler ne najde template
datoteke, katero želimo naložiti preko metode `load`. Če datoteke ne najde,
izpiše napako, izvajanje se konča.

Metoda je uporabna v primeru, če želite imeti več map s template datotekami.
Naprimer, če želite imeti CMS strukturo.

~~~~~~~~~~~~~~~
$paths = array();

/* This is the most important location, everything
   can be overriden from inside the theme. */

$paths[] = "theme/templates/";

/* This is the second most important location,
   it usually defines the look of the cms modules */

$paths[] = "modules/".$module_name."/templates/";

/* This is the least important template location,
   it usually provides system wide templates, like
   a paginator template, an xml / rss template, or
   other very general templates. */

$paths[] = "include/templates/";

$tpl->set_paths($paths);
~~~~~~~~~~~~~~~

Konstrukt `{include}` v template datotekah upošteva nastavitve podane preko
metode `set_paths`. Pred prevajanjem sistem išče datoteke za vključitev na
istih lokacijah, dokler jih ne najde.
