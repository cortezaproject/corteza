---
title: Expression functions
description: Every function of the expression language, with its signature and an example.
outline: [2, 3]
---

<!-- This file is auto-generated from server/pkg/expr/expr_functions.yaml. -->

# Expression functions

Every [expression](./) can call these functions: field value expressions,
value sanitization and validation, contextual roles and the While loop in TAQs.
A name is a function only where it is called: `split(a, b)` is the function,
while a bare `split` reads a variable of that name.

Every example on this page is evaluated by the server's test suite, so the
result shown is the result you get.

## String

### `trim` {#trim}

```ts
trim(s: String): String
```

Removes leading and trailing whitespace.

```js
trim("  hello  ")
// → "hello"
```

### `trimLeft` {#trimLeft}

```ts
trimLeft(s: String, cutset: String): String
```

Removes every leading character that appears in `cutset`. The second argument is a set of characters, not a prefix.

```js
trimLeft("xxhixx", "x")
// → "hixx"
```

### `trimRight` {#trimRight}

```ts
trimRight(s: String, cutset: String): String
```

Removes every trailing character that appears in `cutset`. The second argument is a set of characters, not a suffix.

```js
trimRight("xxhixx", "x")
// → "xxhi"
```

### `toLower` {#toLower}

```ts
toLower(s: String): String
```

Converts every letter to lower case.

```js
toLower("HeLLo")
// → "hello"
```

### `toUpper` {#toUpper}

```ts
toUpper(s: String): String
```

Converts every letter to upper case.

```js
toUpper("hello")
// → "HELLO"
```

### `shortest` {#shortest}

```ts
shortest(first: String, ...rest: String): String
```

Returns the shortest of the given strings, measured in bytes. On a tie the earlier argument wins.

```js
shortest("apple", "fig", "banana")
// → "fig"
```

### `longest` {#longest}

```ts
longest(first: String, ...rest: String): String
```

Returns the longest of the given strings, measured in bytes. On a tie the earlier argument wins.

```js
longest("apple", "fig", "banana")
// → "banana"
```

### `format` {#format}

```ts
format(format: String, ...values: Any): String
```

Formats the values with Go `fmt.Sprintf` verbs such as `%s`, `%v` and `%.2f`.

::: warning
Number literals are floats, so `%d` prints `%!d(float64=3)`. Use `%v`, or `%.0f` for a whole number.
:::

```js
format("%s has %v items", "cart", 3)
// → "cart has 3 items"
```

### `title` {#title}

```ts
title(s: String): String
```

Capitalises the first word and leaves the rest of the string as it is.

::: warning
A single-word input gets a trailing space: `title("hello")` returns `"Hello "`.
:::

```js
title("hello big world")
// → "Hello big world"
```

### `untitle` {#untitle}

```ts
untitle(s: String): String
```

Lower-cases the first letter of the first word.

::: warning
A single-word input gets a trailing space, and an empty string makes the expression fail.
:::

```js
untitle("Hello World")
// → "hello World"
```

### `repeat` {#repeat}

```ts
repeat(s: String, count: Integer): String
```

Returns `s` repeated `count` times. A fractional count is truncated to a whole number.

```js
repeat("ab", 3)
// → "ababab"
```

### `replace` {#replace}

```ts
replace(s: String, old: String, new: String, n: Integer): String
```

Replaces the first `n` occurrences of `old` with `new`; an `n` of `-1` replaces all of them. All four arguments are required.

```js
replace("a-b-c", "-", "+", -1)
// → "a+b+c"
```

### `isUrl` {#isUrl}

```ts
isUrl(s: String): Boolean
```

Reports whether the string is a URL. The scheme is optional, so a bare host such as `example.com` passes.

```js
isUrl("https://example.com")
// → true
```

### `isEmail` {#isEmail}

```ts
isEmail(s: String): Boolean
```

Reports whether the string is an email address.

```js
isEmail("jane@example.com")
// → true
```

### `split` {#split}

```ts
split(s: String, sep: String): Array
```

Splits the string at every occurrence of `sep` and returns the parts.

```js
split("a,b,c", ",")
// → ["a","b","c"]
```

### `join` {#join}

```ts
join(list: Array, sep: String): String
```

Joins the items of the list into one string, with `sep` between them. Items are converted to strings; an empty value gives an empty string.

```js
join(["a", "b", "c"], "-")
// → "a-b-c"
```

### `hasSubstring` {#hasSubstring}

```ts
hasSubstring(s: String, substring: String, caseSensitive: Boolean): Boolean
```

Reports whether `substring` occurs in `s`. The third argument is required: `false` ignores case, `true` matches it exactly.

```js
hasSubstring("Hello World", "world", false)
// → true
```

### `substring` {#substring}

```ts
substring(s: String, start: Integer, end: Integer): String
```

Returns the part of `s` from `start` to `end`, both positions included and counted in bytes from 0. An `end` of `-1` runs to the end of the string; a `start` past the end gives an empty string.

```js
substring("Hello World", 0, 4)
// → "Hello"
```

### `hasPrefix` {#hasPrefix}

```ts
hasPrefix(s: String, prefix: String): Boolean
```

Reports whether `s` starts with `prefix`. Case-sensitive.

```js
hasPrefix("Hello", "He")
// → true
```

### `hasSuffix` {#hasSuffix}

```ts
hasSuffix(s: String, suffix: String): Boolean
```

Reports whether `s` ends with `suffix`. Case-sensitive.

```js
hasSuffix("Hello", "lo")
// → true
```

### `shorten` {#shorten}

```ts
shorten(s: String, unit: String, count: Integer): String
```

Cuts `s` to its first `count` words, or to its first `count` bytes when `unit` is `"char"`, then drops one trailing punctuation mark and appends `" …"`. A string already short enough is returned unchanged.

```js
shorten("The quick brown fox", "word", 2)
// → "The quick …"
```

### `camelize` {#camelize}

```ts
camelize(s: String): String
```

Joins space-separated words into camelCase.

::: warning
The result ends with a trailing space; wrap it in `trim()` where that matters.
:::

```js
trim(camelize("hello big world"))
// → "helloBigWorld"
```

### `snakify` {#snakify}

```ts
snakify(s: String): String
```

Joins space-separated words into lower-case snake_case.

```js
snakify("Hello big World")
// → "hello_big_world"
```

### `match` {#match}

```ts
match(s: String, pattern: String): Boolean
```

Reports whether the regular expression matches anywhere in `s`. The pattern uses Go (RE2) syntax; anchor it with `^` and `$` to match the whole string.

```js
match("abc123", "[0-9]+")
// → true
```

### `base64encode` {#base64encode}

```ts
base64encode(value: String): String
```

Encodes a string or binary value as standard base64. Any other value gives an empty string.

```js
base64encode("hello")
// → "aGVsbG8="
```

## Numeric

### `min` {#min}

```ts
min(...values: Float): Float
```

Returns the smallest of the numbers. Arguments that are not numbers, numeric strings included, are skipped; with no number at all the result is `0`.

```js
min(3, 1, 2)
// → 1
```

### `max` {#max}

```ts
max(...values: Float): Float
```

Returns the largest of the numbers. Arguments that are not numbers, numeric strings included, are skipped; with no number at all the result is `0`.

```js
max(3, 1, 2)
// → 3
```

### `round` {#round}

```ts
round(value: Float, decimals: Float): Float
```

Rounds to the given number of decimal places, halves away from zero. Both arguments are required.

```js
round(3.14159, 2)
// → 3.14
```

### `floor` {#floor}

```ts
floor(value: Float): Float
```

Rounds down to the nearest whole number.

```js
floor(3.7)
// → 3
```

### `ceil` {#ceil}

```ts
ceil(value: Float): Float
```

Rounds up to the nearest whole number.

```js
ceil(3.2)
// → 4
```

### `abs` {#abs}

```ts
abs(value: Float): Float
```

Returns the absolute value.

```js
abs(-4.5)
// → 4.5
```

### `log` {#log}

```ts
log(value: Float): Float
```

Returns the base-10 logarithm.

```js
log(1000)
// → 3
```

### `pow` {#pow}

```ts
pow(base: Float, exponent: Float): Float
```

Raises `base` to the power of `exponent`.

```js
pow(2, 10)
// → 1024
```

### `sqrt` {#sqrt}

```ts
sqrt(value: Float): Float
```

Returns the square root.

```js
sqrt(16)
// → 4
```

### `sum` {#sum}

```ts
sum(...values: Any): Float
```

Adds the arguments. Numeric strings count; anything that cannot be read as a number is skipped.

::: warning
A list passed as one argument is skipped too, so `sum([1, 2, 3])` is `0`. Pass the numbers as separate arguments.
:::

```js
sum(1, 2, "3.5")
// → 6.5
```

### `average` {#average}

```ts
average(...values: Any): Float
```

Returns the mean of the arguments that can be read as numbers, skipping the rest. With no number at all the result is NaN.

```js
average(1, 2, 3)
// → 2
```

### `random` {#random}

```ts
random(a: Float, b?: Float): Float
```

With one argument returns a random number from 0 up to `a`; with two, from `a` up to `b`. The result is fractional; wrap it in `int()` or `floor()` for a whole number. Negative bounds fail.

```js
random(10)
// → a number from 0 up to 10
```

### `int` {#int}

```ts
int(value: Any): Integer
```

Converts to a whole number, truncating any fraction. Strings are read as base-10, so a leading zero is not octal; a string that is not a number gives `0`.

```js
int("042")
// → 42
```

### `float` {#float}

```ts
float(value: Any): Float
```

Converts to a number. A string that is not a number makes the expression fail.

```js
float("3.5")
// → 3.5
```

## Array

### `push` {#push}

```ts
push(list: Array, ...values: Any): Array
```

Returns a copy of the list with the values appended. The original list is not changed; an empty list value gives just the pushed values.

```js
push([1, 2], 3)
// → [1,2,3]
```

### `pop` {#pop}

```ts
pop(list: Array): Any
```

Returns the last item without removing it. An empty list gives an empty value.

```js
pop([1, 2, 3])
// → 3
```

### `shift` {#shift}

```ts
shift(list: Array): Any
```

Returns the first item without removing it. An empty list gives an empty value.

```js
shift([1, 2, 3])
// → 1
```

### `count` {#count}

```ts
count(list: Array, ...values: Any): Integer
```

With only a list, returns its length. With values, returns how many of those values occur in the list — each value counts once, however often it appears. On a string, it instead counts every occurrence of each value.

::: warning
`count([1, 2, 1], 1)` is `1`, not `2`, while `count("banana", "a")` is `3`.
:::

```js
count(["a", "b", "c"], "a", "z")
// → 1
```

### `has` {#has}

```ts
has(list: Array, ...values: Any): Boolean
```

Reports whether any of the values occurs in the list or string. On a key-value object it checks keys instead.

```js
// with {"rec": {"a": 1, "b": 2}}
has(rec, "b")
// → true
```

### `hasAll` {#hasAll}

```ts
hasAll(list: Array, ...values: Any): Boolean
```

Reports whether every one of the values occurs in the list. Works on lists only.

```js
hasAll([1, 2, 3], 1, 3)
// → true
```

### `find` {#find}

```ts
find(list: Array, value: Any): Integer
```

Returns the position of the first item equal to `value`, counting from 0, or `-1` when there is none. Types must match: `"1"` is not found in `[1, 2]`.

```js
find(["a", "b", "c"], "b")
// → 1
```

### `sort` {#sort}

```ts
sort(list: Array, descending: Boolean): Array
```

Returns a sorted copy of the list, ascending when `descending` is `false`. Both arguments are required.

::: warning
Only lists of one comparable type sort — the result of `split()`, or a list field of a record. A list literal such as `[3, 1, 2]` fails with "cannot compare".
:::

```js
sort(split("b,c,a", ","), false)
// → ["a","b","c"]
```

### `splice` {#splice}

```ts
splice(list: Array, start: Float, end: Float): Array
```

Returns the items from position `start` up to, but not including, `end`, counting from 0. An `end` of `-1` runs to the end of the list. A `start` past the end returns the whole list.

```js
splice([1, 2, 3, 4], 1, 3)
// → [2,3]
```

## Date and time

Functions that take a date accept a DateTime value or a string in a common date format such as RFC 3339 (`2024-03-01T10:00:00Z`) or a plain date (`2024-03-01`).

### `now` {#now}

```ts
now(): DateTime
```

Returns the current date and time.

```js
now()
// → the current date and time
```

### `parseISOTime` {#parseISOTime}

```ts
parseISOTime(s: String): DateTime
```

Parses an RFC 3339 timestamp. Unlike the other date functions it accepts nothing else — a plain date such as `2024-03-01` fails.

```js
parseISOTime("2024-03-01T10:00:00Z")
// → "2024-03-01T10:00:00Z"
```

### `parseDuration` {#parseDuration}

```ts
parseDuration(s: String): Duration
```

Parses a duration such as `90m`, `1h30m` or `-2h`. Valid units are `ns`, `us`, `ms`, `s`, `m` and `h`; there is no unit for days.

```js
parseDuration("1h30m")
// → "1h30m0s"
```

### `earliest` {#earliest}

```ts
earliest(first: DateTime, ...rest: DateTime): DateTime
```

Returns the earliest of the dates.

```js
earliest("2024-03-01T00:00:00Z", "2024-01-15T00:00:00Z")
// → "2024-01-15T00:00:00Z"
```

### `latest` {#latest}

```ts
latest(first: DateTime, ...rest: DateTime): DateTime
```

Returns the latest of the dates.

```js
latest("2024-03-01T00:00:00Z", "2024-01-15T00:00:00Z")
// → "2024-03-01T00:00:00Z"
```

### `modTime` {#modTime}

```ts
modTime(date: DateTime, duration: Duration): DateTime
```

Adds a duration to the date. The duration is a Duration value or a string in `parseDuration()` syntax; a negative one moves back.

```js
modTime("2024-03-01T10:00:00Z", "90m")
// → "2024-03-01T11:30:00Z"
```

### `modDate` {#modDate}

```ts
modDate(date: DateTime, days: Integer): DateTime
```

Adds a number of days to the date; a negative number moves back.

```js
modDate("2024-03-01T10:00:00Z", -1)
// → "2024-02-29T10:00:00Z"
```

### `modWeek` {#modWeek}

```ts
modWeek(date: DateTime, weeks: Integer): DateTime
```

Adds a number of weeks to the date; a negative number moves back.

```js
modWeek("2024-03-01T10:00:00Z", 2)
// → "2024-03-15T10:00:00Z"
```

### `modMonth` {#modMonth}

```ts
modMonth(date: DateTime, months: Integer): DateTime
```

Adds a number of months to the date; a negative number moves back. A day that does not exist in the target month rolls over into the next one.

```js
modMonth("2024-01-31T10:00:00Z", 1)
// → "2024-03-02T10:00:00Z"
```

### `modYear` {#modYear}

```ts
modYear(date: DateTime, years: Integer): DateTime
```

Adds a number of years to the date; a negative number moves back. 29 February rolls over to 1 March in a year that has no leap day.

```js
modYear("2024-02-29T10:00:00Z", 1)
// → "2025-03-01T10:00:00Z"
```

### `strftime` {#strftime}

```ts
strftime(date: DateTime, format: String): String
```

Formats the date with POSIX `strftime` directives such as `%Y`, `%m`, `%d`, `%H` and `%M`.

::: warning
`%b` prints milliseconds and `%L` the Unix timestamp in seconds, not their usual POSIX meanings (month name, and none).
:::

```js
strftime("2024-03-01T10:00:00Z", "%d.%m.%Y %H:%M")
// → "01.03.2024 10:00"
```

### `isLeapYear` {#isLeapYear}

```ts
isLeapYear(date: DateTime): Boolean
```

Reports whether the year of the date is a leap year.

```js
isLeapYear("2024-03-01")
// → true
```

### `isWeekDay` {#isWeekDay}

```ts
isWeekDay(date: DateTime): Boolean
```

Reports whether the date falls on Monday to Friday, in the date's own time zone.

```js
isWeekDay("2024-03-02")
// → false
```

### `sub` {#sub}

```ts
sub(later: DateTime, earlier: DateTime): Integer
```

Returns the time from `earlier` to `later` in milliseconds.

::: warning
The first date must be strictly later than the second; equal dates, or the other order, make the expression fail.
:::

```js
sub("2024-03-01T10:00:00Z", "2024-03-01T09:00:00Z")
// → 3600000
```

## Generic

### `coalesce` {#coalesce}

```ts
coalesce(...values: Any): Any
```

Returns the first argument that is not empty-valued (nil). An empty string or `0` is a value and is returned.

```js
// with {"nickname": null}
coalesce(nickname, "anonymous")
// → "anonymous"
```

### `isNil` {#isNil}

```ts
isNil(value: Any): Boolean
```

Reports whether the value is nil — unset, not merely empty.

```js
// with {"nickname": null}
isNil(nickname)
// → true
```

### `isEmpty` {#isEmpty}

```ts
isEmpty(value: Any): Boolean
```

Reports whether the value is nil or its type's zero value: an empty string, `0`, `false`, or a list or object with no items.

```js
isEmpty("")
// → true
```

### `length` {#length}

```ts
length(value: Any): Integer
```

Returns the number of items in a list or object, or of bytes in a string. Any other value, an empty one included, gives `0`.

```js
length([1, 2, 3])
// → 3
```

## JSON

### `toJSON` {#toJSON}

```ts
toJSON(value: Any): String
```

Encodes the value as a JSON string.

::: warning
A key-value object of typed values encodes each value with its type, as `{"@value": 1, "@type": "Float"}`.
:::

```js
toJSON(["a", "b"])
// → "[\"a\",\"b\"]"
```

## Key-value objects

These work on key-value objects: Vars (typed values), KV (string values) and KVV (lists of strings). Each returns a new object and leaves its input unchanged.

### `set` {#set}

```ts
set(object: Vars, key: String, value: Any): Vars
```

Returns a copy of the object with `key` set to `value`.

```js
// with {"rec": {"a": 1, "b": 2}}
set(rec, "c", 3)
// → {"a":1,"b":2,"c":3}
```

### `merge` {#merge}

```ts
merge(object: Vars, ...others: Vars): Vars
```

Returns a copy of the object with the keys of the others added. On a shared key the later object wins. Every argument must be an object value such as a variable — an inline `{…}` literal fails.

```js
// with {"rec": {"a": 1, "b": 2}, "extra": {"b": 20, "c": 3}}
merge(rec, extra)
// → {"a":1,"b":20,"c":3}
```

### `filter` {#filter}

```ts
filter(object: Vars, ...keys: String): Vars
```

Returns a copy of the object with only the given keys.

```js
// with {"rec": {"a": 1, "b": 2}}
filter(rec, "a")
// → {"a":1}
```

### `omit` {#omit}

```ts
omit(object: Vars, ...keys: String): Vars
```

Returns a copy of the object without the given keys.

```js
// with {"rec": {"a": 1, "b": 2}}
omit(rec, "a")
// → {"b":2}
```

## Access control

These functions read the user group hierarchy (Admin Area → System → User Groups, where a group "reports to" parent groups). They are meant for contextual role expressions. The examples assume a group `managers` with user 1, and two groups under it — `sales` (user 2), linked with the relationship name `read`, and `support` (user 3), linked with no name.

### `isDescendantOf` {#isDescendantOf}

```ts
isDescendantOf(userID: ID, owners: Any, ...relationships: String): Boolean
```

True when a group of one of the `owners` sits somewhere below the user's group. `owners` is a user ID or a list of them, typically `resource.ownedBy`. With `relationships`, only links with one of those names — or with no name — are followed. The user's own group does not count, and it is false when either user is in no group.

```js
// with {"userID": 1, "resource": {"ownedBy": 2}}
isDescendantOf(userID, resource.ownedBy)
// → true
```

### `isDescendantOfR` {#isDescendantOfR}

```ts
isDescendantOfR(userID: ID, owners: Any): Boolean
```

`isDescendantOf` following only links named `read`, or with no name.

```js
// with {"userID": 1, "resource": {"ownedBy": 2}}
isDescendantOfR(userID, resource.ownedBy)
// → true
```

### `isDescendantOfC` {#isDescendantOfC}

```ts
isDescendantOfC(userID: ID, owners: Any): Boolean
```

`isDescendantOf` following only links named `create`, or with no name.

```js
// with {"userID": 1, "resource": {"ownedBy": 3}}
isDescendantOfC(userID, resource.ownedBy)
// → true
```

### `isDescendantOfU` {#isDescendantOfU}

```ts
isDescendantOfU(userID: ID, owners: Any): Boolean
```

`isDescendantOf` following only links named `update`, or with no name. The `read` link to `sales` does not count, so this is false.

```js
// with {"userID": 1, "resource": {"ownedBy": 2}}
isDescendantOfU(userID, resource.ownedBy)
// → false
```

### `isDescendantOfD` {#isDescendantOfD}

```ts
isDescendantOfD(userID: ID, owners: Any): Boolean
```

`isDescendantOf` following only links named `delete`, or with no name. It only looks down: a member of `support` is not below itself or above `managers`.

```js
// with {"userID": 3, "resource": {"ownedBy": 1}}
isDescendantOfD(userID, resource.ownedBy)
// → false
```

### `isDescendantOfW` {#isDescendantOfW}

```ts
isDescendantOfW(userID: ID, owners: Any): Boolean
```

`isDescendantOf` following only links named `create`, `update` or `delete`, or with no name.

```js
// with {"userID": 1, "resource": {"ownedBy": 3}}
isDescendantOfW(userID, resource.ownedBy)
// → true
```
