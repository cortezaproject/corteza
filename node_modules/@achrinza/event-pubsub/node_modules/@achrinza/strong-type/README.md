# @achrinza/strong-type

Type checking module for anywhere javascript can run ES6 modules. This includes node, electron and the browsers. Fully isomorphic.

npm @achrinza/strong-type info :  [See npm trends and stats for @achrinza/strong-type](http://npm-stat.com/charts.html?package=@achrinza/strong-type&author=&from=&to=)  
![@achrinza/strong-type npm version](https://img.shields.io/npm/v/@achrinza/strong-type.svg) ![supported node version for @achrinza/strong-type](https://img.shields.io/node/v/@achrinza/strong-type.svg) ![total npm downloads for @achrinza/strong-type](https://img.shields.io/npm/dt/@achrinza/strong-type.svg) ![monthly npm downloads for @achrinza/strong-type](https://img.shields.io/npm/dm/@achrinza/strong-type.svg) ![npm licence for @achrinza/strong-type](https://img.shields.io/npm/l/@achrinza/strong-type.svg)

` npm install --save @achrinza/strong-type `  

[![RIAEvangelist](https://avatars3.githubusercontent.com/u/369041?v=3&s=100)](https://github.com/RIAEvangelist)

GitHub info :  
[![@achrinza/strong-type GitHub Release](https://img.shields.io/github/release/achrinza/strong-type.svg) ![GitHub license @achrinza/strong-type license](https://img.shields.io/github/license/achrinza/strong-type.svg) ![open issues for @achrinza/strong-type on GitHub](https://img.shields.io/github/issues/achrinza/strong-type.svg)](http://riaevangelist.github.io/@achrinza/strong-type/)


Build Info :  
Travis CI (linux,windows & Mac) : [![Build Status](https://travis-ci.org/RIAEvangelist/strong-test.svg?branch=main)](https://travis-ci.org/RIAEvangelist/strong-test)

## What does @achrinza/strong-type do?
`@achrinza/strong-type` allows easy type enforcement for all JS types objects and classes. It is also extensible and provides simple to use type checks for your own custom classes and types should you want to use them.

## Example | strict vs. non-strict

Using `strict` and `non-strict` modes. By default `@achrinza/strong-type` runs `strict` and will `throw` a verbose err you can handle and use if the check fails, or return true if it passes.

To use `non-strict` mode, simply pass `false` to the constructor. In this mode, `@achrinza/strong-type` will return `false` instead of throwing when the check fails, and will still return `true` when it passes.


#### strict
```javascript 
import Is from '@achrinza/strong-type';

//strict
const is = new Is;
const is = new Is();
const is = new Is(true);

//throws
is.string(1);
```

#### non-strict
```javascript
import Is from '@achrinza/strong-type';

//non-strict
const is = new Is(false);

//returns false
is.string(1);

```

## Example | Basic type checking

`@achrinza/strong-type` is intended to be very simple to use.

```javascript 
import Is from '@achrinza/strong-type';

//strict
const is = new Is;

function strongTypeRequired(aNumber,aString,anAsyncFunction){
    is.number(aNumber);
    is.string(aString);
    is.asyncFunction(anAsyncFunction);
}


//this will throw because we do not pass an async Function, but rather a normal function.
strongTypeRequired(1,'a',function(){})


```
#### browser
![Basic Type Checking Example Web](https://raw.githubusercontent.com/achrinza/strong-type/main/docs/img/basicExampleWeb.PNG)

#### node
![Basic Type Checking Example Node](https://raw.githubusercontent.com/achrinza/strong-type/main/docs/img/basicExampleNode.PNG)

## Example | Generator type checking

Generators are notoriously confusing to type check for many devs. This is why we chose to use them as an example.

```javascript 
import Is from '@achrinza/strong-type';

//strict
const is = new Is;

//empty generator for this example's sake
function* myGenFunc(){};
const myGen=myGenFunc();

//we'll show async as well
async function* myAsyncGenFunc(){}
const myAsyncGen=myAsyncGenFunc();

//empty function
function myFunc(){};

//will pass and allow contunue
is.generator(myGen);

/*
will fail because a generatorFunction is a 
GeneratorFunction, not a Generator
*/
try{
    is.generator(myGenFunc);
}catch(err){
    console.log(err);
}

//will pass and allow contunue
is.generatorFunction(myGenFunc);

//will fail because a function is not a generatorFunction
try{
    is.generatorFunction(myFunc);
}catch(err){
    console.log(err);
}

/*
will fail because this is @achrinza/strong-type, a 
generatorFunction is explicitly a GeneratorFunction,
and not a Function
*/
try{
    is.function(myGenFunc);
}catch(err){
    console.log(err);
}

//will fail because a function is not a generatorFunction
try{
    is.generatorFunction(myFunc);
}catch(err){
    console.log(err);
}

//will pass and allow contunue
is.asyncGeneratorFunction(myAsyncGenFunc);

//will pass and allow contunue
is.asyncGenerator(myAsyncGen);

/*
will fail becase @achrinza/strong-type 
asyncGenerators and generators are explicitly different
this is the same for generatorFunctions and functions
*/
try{
    is.asyncGenerator(myGen);
}catch(err){
    console.log(err);
}

try{
    is.generator(myAsyncGen);
}catch(err){
    console.log(err);
}


```

#### browser
![Generator Type Checking Example Web](https://raw.githubusercontent.com/achrinza/strong-type/main/docs/img/generatorExampleWeb.PNG)

#### node
![Generator Type Checking Example Node](https://raw.githubusercontent.com/achrinza/strong-type/main/docs/img/generatorExampleNode.PNG)

## Date example

```javascript 
import Is from '@achrinza/strong-type';

const is = new Is;

//returns true
is.date(new Date()); 

//throws in strict or returns false in non-strict
is.date(1975);

```

## isNaN() vs is.NaN()

Javascripts types are weak by nature, so the built in `isNaN()` function returns true for anything that not a number, but `is.NaN()` only returns true if it is explicitly passed `NaN`.

```js 
import Is from '@achrinza/strong-type';

const is = new Is;

//built in JS isNaN
//returns false
isNaN(1);

//all return true
isNaN(NaN);
isNaN(undefined);
isNaN('a'); 

//@achrinza/strong-type is.NaN all return false in non-strict mode,
//or throw in default strict mode
is.NaN(1);
is.NaN(undefined);
is.NaN('a');

//in @achrinza/strong-type only this returns true
is.NaN(NaN);

```

#### browser
![Date Type Checking Example Web](https://raw.githubusercontent.com/achrinza/strong-type/main/docs/img/dateExampleWeb.PNG)

#### node
![Date Type Checking Example Node](https://raw.githubusercontent.com/achrinza/strong-type/main/docs/img/dateExampleNode.PNG)



## Running the browser and node type support examples
run `npm i` in the root dir of this module to make sure you get the devDependencies installed.

#### node example :  
`npm run nodeExample` The whole screen should be green as all of the types are supported in node.

#### browser examples and tests
`npm start`  
this will spin up a `node-http-server` in this modules root on port 8000. The browser examples are in the  `./example/web/` folder, and the tests are in `./test/`. You can see them by going to this local address : [http://localhost:8000/](http://localhost:8000/example/web/index.html)

Chrome, Opera, and Edge support all the types so all rows will be green.

You will see some red rows in Firefox as it does not yet support all types. The unsupported types will throw type errors when checked/validated.

#### Digital Ocean Static App

We use the free Digital Ocean Static Apps to host a version of the local server. It is exactly the same as if you ran `npm start` on your machine. You can also use this like a CDN as it automatically rebuilds from main/master each time the branch is updated. [@achrinza/strong-type CDN home](https://@achrinza/strong-type-hiprj.ondigitalocean.app/) : https://@achrinza/strong-type-hiprj.ondigitalocean.app/

## Type check methods 

All of these methods take just one arg, the `value` to check. 

|Most Common Type Methods|args|
|-|-|
|`is.globalThis`|`value`|
|`is.array`|`value`|
|`is.bigint`|`value`|
|`is.boolean`|`value`|
|`is.date`|`value`|
|`is.finite`|`value`|
|`is.generator`|`value`|
|`is.asyncGenerator`|`value`|
|`is.infinity`|`value`|
|`is.map`|`value`|
|`is.NaN`|`value`|
|`is.null`|`value`|
|`is.number`|`value`|
|`is.object`|`value`|
|`is.promise`|`value`|
|`is.regExp`|`value`|
|`is.set`|`value`|
|`is.string`|`value`|
|`is.symbol`|`value`|
|`is.undefined`|`value`|
|`is.weakMap`|`value`|
|`is.weakSet`|`value`|

|Function Type Methods|args|
|-|-|
|`is.function`|`value`|
|`is.asyncFunction`|`value`|
|`is.generatorFunction`|`value`|
|`is.asyncGeneratorFunction`|`value`|

|Error Type Methods|args|
|-|-|
|`is.error`|`value`|
|`is.evalError`|`value`|
|`is.rangeError`|`value`|
|`is.referenceError`|`value`|
|`is.syntaxError`|`value`|
|`is.typeError`|`value`|
|`is.URIError`|`value`|

|Buffer/Typed Array Type Methods|args|
|-|-|
|`is.arrayBuffer`|`value`|
|`is.dataView`|`value`|
|`is.sharedArrayBuffer`|`value`|
|`is.bigInt64Array`|`value`|
|`is.bigUint64Array`|`value`|
|`is.float32Array`|`value`|
|`is.float64Array`|`value`|
|`is.int8Array`|`value`|
|`is.int16Array`|`value`|
|`is.int32Array`|`value`|
|`is.uint8Array`|`value`|
|`is.uint8ClampedArray`|`value`|
|`is.uint16Array`|`value`|
|`is.uint32Array`|`value`|

|Intl Type Methods|args|
|-|-|
|`is.intlDateTimeFormat`|`value`|
|`is.intlCollator`|`value`|
|`is.intlDisplayNames`|`value`|
|`is.intlListFormat`|`value`|
|`is.intlLocale`|`value`|
|`is.intlNumberFormat`|`value`|
|`is.intlPluralRules`|`value`|
|`is.intlRelativeTimeFormat`|`value`|

|Garbage Collection Type Methods|args|
|-|-|
|`is.finalizationRegistry`|`value`|
|`is.weakRef`|`value`|

## Core methods 

You can use these to directly check your own types / classes Or extend the Is class to add your own methods in which you use these for checking more types, especially custom types and classes.

|Method|args|description|
|-|-|-|
|`is.throw`|`valueType`, `expectedType`| this will use the valueType and expectedValueType to create and throw a new `TypeError` |
|`is.typeCheck`|`value`, `type`| this will check the javascript spec types returned from `typeof`. So the `type` arg would be a string of `'string'`, `'boolean'`, `'number'`, `'object'`, `'undefined'`, `'bigint'` etc. |
|`is.instanceCheck`|`value`=`new Fake`, `constructor`=`FakeCore`| The core defaults the args to a `Fake` instance and the `FakeCore` class. This allows unsupported js spec types to fail as expected with a `TypeError` instead of a `Reference` or other Error (see the `./example/web/` example in firefox which is missing some support for `Intl` classes). This method compares `value` with the `constructor` to insure the value is an `instanceof` the constructor. |
|`is.symbolStringCheck`|`value`, `type`| This can be used to check the `Symbol.toStringTag` it works on all types, but in the core we only use it to check `generator`, `GeneratorFunction`, `async function`, and `async GeneratorFunction` as these have no other way to check their type. A generator ***for example*** has a type of `[object generator]` this way. So you pass in an expected `generator` as `value` and the string `'generator'` as the type, and we handle the rest including lowercasing everything to insure cross browser and platform checking |
|`is.compare`|`value`, `targetValue`, `typeName`| this will do an explicit compare on the `value` and `targetValue`. In the core, we only use this for JS primitives/constants that have no other way to check such as `Infinity` and `globalThis`. The type name is the string representation of the class type, or a very explicit error string as the only place this arg is ever used is when the `compare` results in a `throws`. |

## Extending the Is class for your own Types

If you are using type checking on your own types in production, its probably wise for yout to just go ahead and extend the module rather than calling the more cumbersome `Core Methods` many times.

#### custom Pizza type
```javascript
//custom class type constructor
class Pizza{
    constructor(topping){
        this.eat=true;
    }
}

export {default:Pizza, Pizza}
```

#### extension
```javascript
import Is from '@achrinza/strong-type';
import Pizza from 'my-delicious-pizza';

class IsMy extends Is{
    //custom pizza type
    pizza(value){
        return this.instanceCheck(value,Pizza);
    }
}

export={default:IsMy, IsMy};
```

#### test
```javascript
import IsMy from 'my-delicious-typechecks';
import Pizza from 'my-delicious-pizza';

const is=new IsMy;

//will throw because 42 is not a Pizza Type
//and 
is.pizza(42)

```

#### browser
![Pizza Type Checking Example Web](https://raw.githubusercontent.com/achrinza/strong-type/main/docs/img/pizzaExampleWeb.PNG)

#### node
![Pizza Type Checking Example Node](https://raw.githubusercontent.com/achrinza/strong-type/main/docs/img/pizzaExampleNode.PNG)

